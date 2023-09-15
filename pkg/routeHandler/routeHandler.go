package routeHandler

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"ncpdpDestination/pkg/pbmlib"

	"time"
)

func HandleRoute(nc *nats.Conn, pbm pbmlib.PBM, route, natsPublicSubject, natsPrivateSubject, natsQueue string, timeout time.Duration) (*nats.Subscription, *nats.Subscription, error) {

	privateId := uuid.New().String()
	natsPrivateSubject = natsPrivateSubject + "." + privateId + "." + route
	natsPublicSubject = natsPublicSubject + "." + route

	sub, err := subscribeToSubject(nc, pbm, natsQueue, natsPublicSubject, natsPrivateSubject, timeout)
	if err != nil {
		return nil, nil, err
	}
	privateSub, err := subscribeToSubject(nc, pbm, natsQueue, natsPrivateSubject, natsPrivateSubject, timeout)
	if err != nil {
		return nil, nil, err
	}

	return sub, privateSub, nil
}

func subscribeToSubject(nc *nats.Conn, pbm pbmlib.PBM, natsQueue, subject, natsPrivateSubject string, timeout time.Duration) (*nats.Subscription, error) {
	sub, err := nc.QueueSubscribe(subject, natsQueue, func(msg *nats.Msg) {
		data := msg.Data
		headers := map[string][]string(msg.Header)
		//Claim Object

		go postToPBM(pbm, data, headers, timeout, func(response *pbmlib.Response, respHeader map[string][]string, err *pbmlib.ErrorInfo) {
			if respHeader == nil {
				respHeader = make(map[string][]string)
			}
			respMsg := nats.Msg{
				Data:    response.ToJSON(),
				Header:  nats.Header(respHeader),
				Subject: msg.Reply,
			}
			respMsg.Header.Add("privateSubject", natsPrivateSubject)
			nc.PublishMsg(&respMsg)
		})

	})
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func postToPBM(pbm pbmlib.PBM, requestBuffer []byte, headers map[string][]string, timeout time.Duration, f func(response *pbmlib.Response, respHeader map[string][]string, err *pbmlib.ErrorInfo)) {

	clm := pbmlib.Claim{}
	err := json.Unmarshal(requestBuffer, &clm)
	if err != nil {
		//build response with unable to parse request claim
		claim := pbmlib.Claim{}
		claim.TransactionData.NcpdpData = string(requestBuffer)
		resp := pbmlib.Response{}
		resp.BuildResponseError(claim, pbmlib.ErrorCode.TRX01, time.Now())
		f(&resp, nil, &pbmlib.ErrorCode.TRX01)
		return
	}
	clm.TimeRcvd = time.Now()
	responseBuffer, responseHeaders, erroInfo := pbm.Post(clm, headers, timeout)
	if erroInfo.Code == pbmlib.ErrorCode.TRX00.Code {
		//build response
		resp := pbmlib.Response{}
		resp.BuildResponseSuccess(clm, clm.TimeRcvd, responseBuffer)
		f(&resp, responseHeaders, &pbmlib.ErrorCode.TRX00)
		return
	}

	resp := pbmlib.Response{}
	resp.BuildResponseError(clm, erroInfo, clm.TimeRcvd)

	f(&resp, responseHeaders, &erroInfo)
	//build response

}
