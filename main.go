package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/amaloy/ampub"
)

const version = "0"
const defaultPubsubTimeout = 10 * time.Second
const envVarPrefix = "AMPUB_PUBSUB_"

var pubsubTimeout = defaultPubsubTimeout

func main() {
	log.Printf("ampub-gcppubsub version %s", version)

	// Create GCP Pub/Sub Publisher
	p, err := newPubSubPublisher(getEnvVar("PROJECTID", true))
	if err != nil {
		log.Fatal(err)
	}

	// Assign new value to pubsubTimeout if it is being overridden
	assignTimeoutOverride()

	// Ensure any requested topics are created
	err = ensureTopics(p.client)
	if err != nil {
		log.Fatal(err)
	}

	// Start AmPub API server
	ampub := new(ampub.AmPub)
	ampub.Run(p)
}

func getEnvVar(name string, required bool) string {
	name = envVarPrefix + name
	value, found := os.LookupEnv(name)
	if !found && required {
		log.Fatalf("Environment variable required but not found: %s", name)
	}
	log.Printf("env: %s=%s", name, value)
	return value
}

func assignTimeoutOverride() {
	envValue := getEnvVar("TIMEOUTMS", false)
	intValue, err := strconv.ParseUint(envValue, 10, 0)
	if err == nil {
		pubsubTimeout = time.Duration(intValue) * time.Millisecond
	}
}

func ensureTopics(client *pubsub.Client) (err error) {
	envValue := getEnvVar("TOPICS", false)
	if len(envValue) == 0 {
		return
	}

	topics := strings.Split(envValue, ",")
	for _, t := range topics {
		log.Printf("Ensuring topic exists: %s", t)

		ctx, cancel := context.WithTimeout(context.Background(), pubsubTimeout)
		defer cancel()

		// Check if topic exists
		topic := client.Topic(t)
		exists, err := topic.Exists(ctx)
		if err != nil {
			break
		}

		// Create topic if necessary
		if !exists {
			_, err = client.CreateTopic(ctx, t)
			if err != nil {
				break
			}
		}
	}
	return
}
