# kit4go

Kits for go

## Usage

```bash
go get -u github.com/xwi88/kit4go 
```

## Supported or TODO

* [x] json
    wrap: encoding/json & jsoniter
* [x] datetime util
* [x] [version util](https://github.com/xwi88/version)
* [x] mysql client
* [x] aerospike client
* [x] kafka producer: async & sync producer
* [x] kafka consumer: consumer & consumerGroup

## issue

- panic: non-positive interval for NewTicker
    - `sarama-cluster@v2.1.15+incompatible/consumer.go`
    - 1. `config.Consumer.Offsets.CommitInterval = config.Consumer.Offsets.AutoCommit.Interval`
    - 2. kafka `github.com/Shopify/sarama v1.26.4` -> `github.com/Shopify/sarama v1.24.1`
