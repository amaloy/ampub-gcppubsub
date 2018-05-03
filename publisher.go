package main

import (
	"context"

	"cloud.google.com/go/pubsub"
)

// pubSubPublisher - Implementation of ampub.Publisher that publishes to GCP Pub/Sub
type pubSubPublisher struct {
	client *pubsub.Client
}

func newPubSubPublisher(projectID string) (*pubSubPublisher, error) {
	// Create new GCP Pub/Sub Client
	ctx, cancel := context.WithTimeout(context.Background(), pubsubTimeout)
	defer cancel()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return &pubSubPublisher{
		client: client,
	}, nil
}

func (p *pubSubPublisher) Publish(ctx context.Context, topicName string, key string, data []byte) (err error) {
	// Get reference to topic
	topic := p.client.Topic(topicName)
	defer topic.Stop()

	// Construct Pub/Sub Message
	message := &pubsub.Message{
		Data: data,
	}
	if len(key) > 0 {
		message.Attributes["key"] = key
	}

	// Publish to Pub/Sub
	r := topic.Publish(ctx, message)

	// Block for feedback from Pub/Sub about the publishing
	tctx, cancel := context.WithTimeout(ctx, pubsubTimeout)
	defer cancel()
	_, err = r.Get(tctx)

	return
}
