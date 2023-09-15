package routeHandler

import (
	"github.com/nats-io/nats.go"
	"ncpdpDestination/pkg/pbm"

	"time"
)

func HandleRoute(nc *nats.Conn, pp pbm.PBM, route, natsPublicSubject, natsPrivateSubject string, timeout time.Duration) (*nats.Subscription, *nats.Subscription, error) {

	sub, err := nc.Subscribe(natsPublicSubject, func(msg *nats.Msg) {
		data := msg.Data
		headers := map[string][]string(msg.Header)

		go postToPBM(pp, data, headers, timeout, func(resp []byte, err pbm.ErrorInfo) {

		})

	})
	if err != nil {
		return nil, nil, err
	}

	//private sub
	privSub, err := nc.Subscribe(natsPrivateSubject, func(msg *nats.Msg) {
		data := msg.Data
		headers := map[string][]string(msg.Header)

		go postToPBM(pp, data, headers, timeout, func(resp []byte, err pbm.ErrorInfo) {

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

func postToPBM(pbm pbm.PBM, data []byte, headers map[string][]string, timeout time.Duration, f func(resp []byte, err pbm.ErrorInfo)) {
	resp, Error := pbm.Post(data, headers, timeout)
	f(resp, Error)
}
