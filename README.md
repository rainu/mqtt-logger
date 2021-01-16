# mqtt-logger

A simple MQTT client written in go that subscribes to a configurable list of MQTT topics on the specified broker 
and logs the whole payload to stdout.

# Get the Binary
You can build it on your own (you will need [golang](https://golang.org/) installed):
```bash
go build -a -installsuffix cgo ./cmd/mqtt-logger/
```

Or you can download the release binaries: [here](https://github.com/rainu/mqtt-logger/releases/latest)

## Running

Configuration is taken from the environment or via arguments, for example:

```bash
export MQTT_BROKER_ADDRESS="tcp://localhost:1883"
export MQTT_USERNAME="foo"
export MQTT_PASSWORD="bar"

./mqtt-logger
```
or
```bash
./mqtt-logger --broker "tcp://localhost:1883" --user "foo" --password "bar"
```

### Available Environment/Arguments

|Argument|Environment|Default|Description|
|--------|-----------|-------|-----------|
|--topic|MQTT_TOPICS| # |List of the mqtt topics to subscribe to.|
|--blacklist|MQTT_TOPIC_BLACKLIST| |List of regular expression.|
|--broker|MQTT_BROKER_ADDRESS|  |The mqtt broker address.|
|--user|MQTT_USERNAME|  |The mqtt username.|
|--password|MQTT_PASSWORD|  |The mqtt password.|
|--timeout|MQTT_TIMEOUT| 1m |The timeout for the mqtt connection.|
|--ca-cert|MQTT_CA_CERT|  |The path of the CA-Cert file for secure mqtt connection.|
