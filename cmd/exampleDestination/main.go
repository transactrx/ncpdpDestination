package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/transactrx/ncpdpDestination/pkg/dummypbm"
	"github.com/transactrx/ncpdpDestination/pkg/natshelper"
	"github.com/transactrx/ncpdpDestination/pkg/routeHandler"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Config struct {
	NatsURL                  string
	NatsJWT                  string
	NatsKey                  string
	NatsPrivateSubjectPrefix string
	NatsPublicSubject        string
	NatsQueue                string
	Routes                   []string
}

var publicSubscriptions map[string]*nats.Subscription = make(map[string]*nats.Subscription)
var privateSubscriptions map[string]*nats.Subscription = make(map[string]*nats.Subscription)

func main() {
	dummyPBM := dummypbm.DummyPBM{}
	cfg := readConfiguration()

	nc, err := natshelper.CreateNatsClient(cfg.NatsJWT, cfg.NatsKey, cfg.NatsURL)
	if err != nil {
		log.Panicf("error while connecting to nats: %v", err)
	}

	for i, route := range cfg.Routes {
		if strings.Trim(route, "") != "" {
			go func() {
				publicSub, privateSub, err := routeHandler.HandleRoute(nc, &dummyPBM, route, cfg.NatsPublicSubject, cfg.NatsPrivateSubjectPrefix, cfg.NatsQueue, time.Second*20)
				if err != nil {
					log.Panicf("error while handling route %s: %v", route, err)
				}

				routeK := fmt.Sprintf("%d", i)
				publicSubscriptions[routeK] = publicSub
				privateSubscriptions[routeK] = privateSub
			}()
		}

	}
	dummyPBM.Start()

	DrainSubscriptions()

}

func DrainSubscriptions() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range sigChan {
			switch sig {
			case syscall.SIGINT:
				for route, sub := range privateSubscriptions {
					log.Printf("Draining private subscription for route %s", route)
					_ = sub.Drain()
				}
				for route, sub := range publicSubscriptions {
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

func readConfiguration() Config {
	cfg := Config{}
	cfg.NatsURL = getEnvironmentVariableOrPanic("NATS_URL")
	cfg.NatsJWT = getEnvironmentVariableOrPanic("NATS_JWT")
	cfg.NatsKey = getEnvironmentVariableOrPanic("NATS_KEY")
	cfg.NatsPrivateSubjectPrefix = getEnvironmentVariableOrPanic("NATS_PRIVATE_SUBJECT_PREFIX")
	cfg.NatsPublicSubject = getEnvironmentVariableOrPanic("NATS_PUBLIC_SUBJECT")
	cfg.NatsQueue = getEnvironmentVariableOrDefault("NATS_QUEUE", "EXAMPLE_DEST")
	cfg.Routes = strings.Split(getEnvironmentVariableOrPanic("ROUTES"), ",")
	return cfg
}

func getEnvironmentVariableOrPanic(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Panicf("Environment variable %s is missing", key)
	}
	return value
}
func getEnvironmentVariableOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
