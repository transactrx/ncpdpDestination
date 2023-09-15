package natshelper

import (
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"time"
)

func setupConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 8760 * time.Hour
	reconnectDelay := time.Second
	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		log.Printf("%s nats.DisconnectErrHandler disconnected due to: %s, will attempt reconnects for %.0fm", time.Now(), err, totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("%s nats.ReconnectHandler reconnected [%s]", time.Now(), nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Printf("%s nats.ClosedHandlerExiting: %v", time.Now(), nc.LastError())
		os.Exit(-1)
	}))

	return opts
}

func CreateNatsClient(jwt, key, url string) (*nats.Conn, error) {
	opts := []nats.Option{nats.UserJWTAndSeed(jwt, key)}
	opts = setupConnOptions(opts)
	return nats.Connect(url, opts...)

}
