# MarketWatcher Server

The MarketWatcher Server is a gRPC server that streams real-time stock prices.

The application is written in Go and is meant to illustrate some interesting features of gRPC and other frameworks, while demonstrating some best practices. The application is only a demo, meant primarily as a learning tool.

## Pre-requisites

You will want to install the [MarketWatcher App](https://github.com/rfield/MarketWatcherApp), and Android application that exercises some of the server's capabilities. 

Additionally, you will need to install Go.

The server also depends on a very simple Postgres database for it's list of users and their stock holdings. So you should install Postgres as well. The DDL to create the tables is in the ./sql directory.

You will want to create an account with the MarketData.App API, or find another source for stock prices. Access to the MarketData.app APIs relies on a token associated with your account, usually pulled from the MARKETDATA_TOKEN environment variable. See References below for additional information.

## How to Run

Do this:

```bash
make
bin/mkt_server
```

## References

[MarketData.App APIs for Real-time stock prices](https://www.marketdata.app/sdk/go/golang-stock-api/real-time-stock-api/)

[Google API Design Guide](https://docs.cloud.google.com/apis/design)

[Google API Improvement Proposals](https://google.aip.dev/)

[gRPC](https://grpc.io/docs/what-is-grpc/introduction/)

## Google Protobuf Dependencies

[resource.proto in Google's Github](https://github.com/googleapis/googleapis/blob/master/google/api/resource.proto)

[descriptor.proto in Google's Github](https://github.com/protocolbuffers/protobuf/blob/main/src/google/protobuf/descriptor.proto)
