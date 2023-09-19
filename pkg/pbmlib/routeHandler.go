package pbmlib

//import (
//	"github.com/google/uuid"
//	"github.com/nats-io/nats.go"
//	"time"
//)
//
//func (ph *PBMHandler) HandleRoute(route, natsPublicSubject, natsPrivateSubject, natsQueue string, timeout time.Duration) (*nats.Subscription, *nats.Subscription, error) {
//
//	privateId := uuid.New().String()
//	natsPrivateSubject = natsPrivateSubject + "." + privateId + "." + route
//	natsPublicSubject = natsPublicSubject + "." + route
//
//	sub, err := subscribeToSubject(natsQueue, natsPublicSubject, natsPrivateSubject, timeout)
//	if err != nil {
//		return nil, nil, err
//	}
//	privateSub, err := subscribeToSubject(nc, pbm, natsQueue, natsPrivateSubject, natsPrivateSubject, timeout)
//	if err != nil {
//		return nil, nil, err
//	}
//
//	return sub, privateSub, nil
//}
//
//func subscribeToSubject(nc *nats.Conn, pbm PBM, natsQueue, subject, natsPrivateSubject string, timeout time.Duration) (*nats.Subscription, error) {
//	sub, err := nc.QueueSubscribe(subject, natsQueue, func(msg *nats.Msg) {
//		data := msg.Data
//		headers := map[string][]string(msg.Header)
//		//Claim Object
//
//		go postToPBM(pbm, data, headers, timeout, func(response *Response, respHeader map[string][]string, err *ErrorInfo) {
//			if respHeader == nil {
//				respHeader = make(map[string][]string)
//			}
//			respMsg := nats.Msg{
//				Data:    response.ToJSON(),
//				Header:  nats.Header(respHeader),
//				Subject: msg.Reply,
//			}
//			respMsg.Header.Add("privateSubject", natsPrivateSubject)
//			nc.PublishMsg(&respMsg)
//		})
//
//	})
//	if err != nil {
//		return nil, err
//	}
//	return sub, nil
//}
