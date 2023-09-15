package main

import (
	"log"
	"ncpdpDestination/pkg/dummypbm"
	"ncpdpDestination/pkg/natshelper"
	"ncpdpDestination/pkg/routeHandler"
	"os"
	"strings"
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

func main() {
	d := dummypbm.DummyPBM{}
	cfg := readConfiguration()
	nc, err := natshelper.CreateNatsClient(cfg.NatsJWT, cfg.NatsKey, cfg.NatsURL)
	if err != nil {
		log.Panicf("error while connecting to nats: %v", err)
	}

	for _, route := range cfg.Routes {
		if strings.Trim(route, "") != "" {
			go routeHandler.HandleRoute(nc, &d, route, cfg.NatsPublicSubject, cfg.NatsPrivateSubjectPrefix, time.Second*20)
		}

	}

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
