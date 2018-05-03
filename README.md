# ampub-gcppubsub
A Google Cloud Pub/Sub AmPub publisher for [AmPub](https://github.com/amaloy/ampub)

# Message Format
The AmPub API provides a topic string, request body bytes, and optionally a key string.
* topic string = The Pub/Sub topic
* request body bytes = The message data
* key string = A message attribute where `key={key string}`

# Environment Variables
* `AMPUB_PUBSUB_PROJECTID` - required, the project that the Pub/Sub client will be scoped to
* `AMPUB_PUBSUB_TIMEOUTMS` - optional, the timeout to use when communicating with Pub/Sub, default is 10000 ms
* `AMPUB_PUBSUB_TOPICS` - optional, comma-separated list of topics to ensure exist on startup
* All variables used by the [Google Cloud Pub/Sub Go client](https://godoc.org/cloud.google.com/go/pubsub) are also valid

# Docker Example
```
docker build -t amaloy/ampub-gcppubsub .

docker run -e AMPUB_PUBSUB_PROJECTID=myproject -e AMPUB_PUBSUB_TOPICS=mytopic amaloy/ampub-gcppubsub
```