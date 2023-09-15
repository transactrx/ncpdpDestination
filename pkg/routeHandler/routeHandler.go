package routeHandler

import (
	"github.com/nats-io/nats.go"
	"ncpdpDestination/pkg/pbmlib"

	"time"
)

func HandleRoute(nc *nats.Conn, pbm pbmlib.PBM, route, natsPublicSubject, natsPrivateSubject string, timeout time.Duration) (*nats.Subscription, *nats.Subscription, error) {

	sub, err := nc.Subscribe(natsPublicSubject, func(msg *nats.Msg) {
		data := msg.Data
		headers := map[string][]string(msg.Header)

		go postToPBM(pbm, data, headers, timeout, func(resp []byte, err pbmlib.ErrorInfo) {

		})

	})
	if err != nil {
		return nil, nil, err
	}

	//private sub
	privSub, err := nc.Subscribe(natsPrivateSubject, func(msg *nats.Msg) {
		data := msg.Data
		headers := map[string][]string(msg.Header)

		go postToPBM(pbm, data, headers, timeout, func(resp []byte, err pbmlib.ErrorInfo) {

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

func postToPBM(pbm pbmlib.PBM, data []byte, headers map[string][]string, timeout time.Duration, f func(resp []byte, err pbmlib.ErrorInfo)) {
	resp, Error := pbm.Post(data, headers, timeout)
	f(resp, Error)
}
