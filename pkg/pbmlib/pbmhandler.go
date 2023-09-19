package pbmlib

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/transactrx/ncpdpDestination/pkg/natshelper"
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

type PBMHandler struct {
	natsURL                  string
	natsJWT                  string
	natsKey                  string
	natsPublicSubject        string
	natsPrivateSubjectPrefix string
	natsQueue                string
	routes                   []string
	nc                       *nats.Conn
	pbms                     []handledPBM
	publicSubscriptions      map[string]*nats.Subscription
	privateSubscriptions     map[string]*nats.Subscription
	overAllTimeOut           time.Duration
}

type handledPBM struct {
	id             string
	privateSubject string
	pbm            PBM
	activeCalls    atomic.Int32
}

func NewPBMHandler() (*PBMHandler, error) {
	ph, err := createHandlerFromConfig()
	if err != nil {
		return nil, err
	}

	nc, err := natshelper.CreateNatsClient(ph.natsURL, ph.natsKey, ph.natsURL)
	if err != nil {
		return nil, err
	}
	ph.nc = nc
	ph.privateSubscriptions = make(map[string]*nats.Subscription)
	ph.publicSubscriptions = make(map[string]*nats.Subscription)
	return ph, err
}

func (ph *PBMHandler) HandlePBMS(pbm []PBM, routes []string) error {

	ph.routes = routes
	ph.pbms = make([]handledPBM, len(pbm))
	for i, pbm := range pbm {
		id := uuid.New().String()
		ph.pbms[i] = handledPBM{
			id:             id,
			pbm:            pbm,
			privateSubject: ph.natsPrivateSubjectPrefix + "." + id,
		}
	}

	err := ph.handlePublicRoutes(routes)
	if err != nil {
		return err
	}
	err = ph.handlePrivateRoutes(routes)
	if err != nil {
		return err
	}

	return nil

}

func (ph *PBMHandler) handlePrivateRoutes(routes []string) error {

	for i := 0; i < len(ph.pbms); i++ {
		sub, err := ph.nc.QueueSubscribe(ph.pbms[i].privateSubject, ph.natsQueue, func(msg *nats.Msg) {

			//select leastBusyPbm with the least active calls;
			var privatePBM *handledPBM = nil
			for i := 0; i < len(ph.pbms); i++ {

				if ph.pbms[i].privateSubject == msg.Subject {
					privatePBM = &ph.pbms[i]
					break
				}
			}
			if privatePBM == nil {
				//error
				log.Printf("Unexepected error: private subject %s not found", msg.Subject)
			}

			go privatePBM.post(msg.Data, map[string][]string(msg.Header), ph.overAllTimeOut, true, func(response *Response, respHeader map[string][]string, err *ErrorInfo) {

				if respHeader == nil {
					respHeader = make(map[string][]string)
				}
				respMsg := nats.Msg{
					Data:    response.ToJSON(),
					Header:  nats.Header(respHeader),
					Subject: msg.Reply,
				}
				respMsg.Header.Add("privateSubject", privatePBM.privateSubject)
				ph.nc.PublishMsg(&respMsg)
			})
		})
		if err != nil {
			return err
		}
		ph.privateSubscriptions[ph.pbms[i].privateSubject] = sub

	}

	return nil
}

func (ph *PBMHandler) handlePublicRoutes(routes []string) error {
	for i := 0; i < len(ph.routes); i++ {
		subject := ph.natsPublicSubject + "." + routes[i]
		sub, err := ph.nc.QueueSubscribe(subject, ph.natsQueue, func(msg *nats.Msg) {

			//select leastBusyPbm with the least active calls;
			var leastBusyPbm *handledPBM = nil
			for i := 0; i < len(ph.pbms); i++ {

				if leastBusyPbm == nil {
					leastBusyPbm = &ph.pbms[i]
				} else if ph.pbms[i].activeCalls.Load() < leastBusyPbm.activeCalls.Load() {
					leastBusyPbm = &ph.pbms[i]
				}
			}

			go leastBusyPbm.post(msg.Data, map[string][]string(msg.Header), ph.overAllTimeOut, false, func(response *Response, respHeader map[string][]string, err *ErrorInfo) {

				if respHeader == nil {
					respHeader = make(map[string][]string)
				}
				respMsg := nats.Msg{
					Data:    response.ToJSON(),
					Header:  nats.Header(respHeader),
					Subject: msg.Reply,
				}
				respMsg.Header.Add("privateSubject", leastBusyPbm.privateSubject)
				ph.nc.PublishMsg(&respMsg)
			})
		})
		if err != nil {
			return err
		}
		ph.publicSubscriptions[subject] = sub
	}
	return nil
}

func (ph *PBMHandler) HandleShutdownAndDrainNats() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range sigChan {
			switch sig {
			case syscall.SIGINT:
				for route, sub := range ph.publicSubscriptions {
					log.Printf("Draining private subscription for route %s", route)
					_ = sub.Drain()
				}
				for route, sub := range ph.privateSubscriptions {
					log.Printf("Draining public subscription for route %s", route)
					_ = sub.Drain()
				}
				os.Exit(1)

			case syscall.SIGTERM:
				fmt.Println("Received SIGTERM!")
				// Handle SIGTERM specific tasks here
				os.Exit(1)
			}
		}
	}()

	select {}
}

func (hpbm *handledPBM) post(requestBuffer []byte, headers map[string][]string, timeout time.Duration, privateMessage bool, f func(response *Response, respHeader map[string][]string, err *ErrorInfo)) {

	clm := Claim{}
	hpbm.activeCalls.Add(1)
	defer hpbm.activeCalls.Add(-1)
	err := json.Unmarshal(requestBuffer, &clm)
	if err != nil {
		//build response with unable to parse request claim
		claim := Claim{}
		claim.TransactionData.NcpdpData = string(requestBuffer)
		resp := Response{}
		resp.BuildResponseError(claim, ErrorCode.TRX01, time.Now())
		f(&resp, nil, &ErrorCode.TRX01)
		return
	}
	clm.TimeRcvd = time.Now()
	responseBuffer, responseHeaders, erroInfo := hpbm.pbm.Post(clm, headers, timeout, privateMessage)
	if erroInfo.Code == ErrorCode.TRX00.Code {
		//build response
		resp := Response{}
		resp.BuildResponseSuccess(clm, clm.TimeRcvd, responseBuffer)
		f(&resp, responseHeaders, &ErrorCode.TRX00)
		return
	}

	resp := Response{}
	resp.BuildResponseError(clm, erroInfo, clm.TimeRcvd)

	f(&resp, responseHeaders, &erroInfo)
	//build response

}

func createHandlerFromConfig() (*PBMHandler, error) {
	var err error
	pbmHandler := PBMHandler{}

	pbmHandler.natsURL, err = getEnvironmentVariable("NATS_URL")
	if err != nil {
		return nil, err
	}

	pbmHandler.natsJWT, err = getEnvironmentVariable("NATS_JWT")
	if err != nil {
		return nil, err
	}

	pbmHandler.natsKey, err = getEnvironmentVariable("NATS_KEY")
	pbmHandler.natsJWT, err = getEnvironmentVariable("NATS_JWT")
	if err != nil {
		return nil, err
	}

	pbmHandler.natsPrivateSubjectPrefix, err = getEnvironmentVariable("NATS_PRIVATE_SUBJECT_PREFIX")
	pbmHandler.natsJWT, err = getEnvironmentVariable("NATS_JWT")
	if err != nil {
		return nil, err
	}

	pbmHandler.natsPublicSubject, err = getEnvironmentVariable("NATS_PUBLIC_SUBJECT")
	pbmHandler.natsJWT, err = getEnvironmentVariable("NATS_JWT")
	if err != nil {
		return nil, err
	}

	pbmHandler.natsQueue = getEnvironmentVariableOrDefault("NATS_QUEUE", "EXAMPLE_DEST")
	pbmHandler.natsJWT, err = getEnvironmentVariable("NATS_JWT")
	if err != nil {
		return nil, err
	}

	return &pbmHandler, nil
}

func getEnvironmentVariable(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("environment variable %s is missing", key)
	}
	return value, nil
}
func getEnvironmentVariableOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
