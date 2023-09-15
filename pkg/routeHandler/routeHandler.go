package routeHandler

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"ncpdpDestination/pkg/pbmlib"

	"time"
)

func HandleRoute(nc *nats.Conn, pbm pbmlib.PBM, route, natsPublicSubject, natsPrivateSubject string, timeout time.Duration) (*nats.Subscription, *nats.Subscription, error) {

	sub, err := nc.Subscribe(natsPublicSubject, func(msg *nats.Msg) {
		data := msg.Data
		headers := map[string][]string(msg.Header)
		//Claim Object

		go postToPBM(pbm, data, headers, timeout, func(response *pbmlib.Response, respHeader map[string][]string, err *pbmlib.ErrorInfo) {
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
		return nil, nil, err
	}

	//private sub
	privSub, err := nc.Subscribe(natsPrivateSubject, func(msg *nats.Msg) {
		data := msg.Data
		headers := map[string][]string(msg.Header)

		go postToPBM(pbm, data, headers, timeout, func(resp []byte, respHeader map[string][]string, err pbmlib.ErrorInfo) {

			nc.PublishMsg(&nats.Msg{
				Data: resp,
			})
		})

	})
	if err != nil {
		return sub, nil, err
	}

	return sub, privSub, nil
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

	responseBuffer, responseHeaders, erroInfo := pbm.Post(clm, headers, timeout)
	if erroInfo.Code != "TRX00" {
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
