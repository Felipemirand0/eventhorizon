# EventHorizon

An open source [CloudEvents](https://cloudevents.io) implementation to allow handling a massive number of events.

## Introduction

EventHorizon provides a uniform way for handling events ([CloudEvents](https://cloudevents.io)) in a Kubernetes native way, by implementing the following resources:

- **[Singularity](samples/kubernetes/Singularity.yml)** - describes the instance itself, the transport method (currently supports HTTP or NATS).
- **[CloudEventHandler](samples/kubernetes/CloudEventHandler.yml)** - describes event handler itself, the encoder and the output.
- **[CloudEventValidator](samples/kubernetes/CloudEventValidator.yml)** - an optional validator, in case you want to enforce specific events.
- **[CloudEventOutput](samples/kubernetes/CloudEventOutput.yml)** - describes where to send the event output to (currently supports Fluentd).

## Current status

Under active development, any contributions are welcome.

## Installing

Download latest release at [GitHub releases page](https://github.com/acesso-io/eventhorizon/releases), extract content and open a terminal inside the folder.

### Kubernetes

Create the custom resources:

```shell
kubectl apply -f install/kubernetes/crds/
```

Install the demo app:

```shell
kubectl apply -f install/kubernetes/demo.yml
```

### Standalone

Set those environment variables:

```shell
export EVENTHORIZON_MODE=standalone
export EVENTHORIZON_NAME=eventhorizon
export EVENTHORIZON_STANDALONE_CONFIG=samples/standalone/stdout.yml
export EVENTHORIZON_LOGGING_LEVEL=info
export EVENTHORIZON_LOGGING_PRETTY=true
```

Start EventHorizon:

```shell
./bin/eventhorizon
```

## Environment variables

| Variable                           | Default                                   | Options                                |
| ---------------------------------- | ----------------------------------------- | -------------------------------------- |
| EVENTHORIZON_MODE                  | kubernetes                                | kubernetes, standalone                 |
| EVENTHORIZON_NAME                  | kube-system/eventhorizon                  |                                        |
| EVENTHORIZON_KUBERNETES_INCLUSTER  | true                                      |                                        |
| EVENTHORIZON_KUBERNETES_KUBECONFIG |                                           |                                        |
| EVENTHORIZON_KUBERNETES_MASTERURL  |                                           |                                        |
| EVENTHORIZON_STANDALONE_CONFIG     | /opt/acesso/samples/standalone/stdout.yml |                                        |
| EVENTHORIZON_LOGGING_ENABLED       | true                                      |                                        |
| EVENTHORIZON_LOGGING_LEVEL         | info                                      | debug, info, warn, error, fatal, panic |
| EVENTHORIZON_LOGGING_PRETTY        | false                                     |                                        |

## Benchmarks

Basic comparison running on a MacBook Pro (Retina, 13-inch, Early 2015), 2,7 GHz Intel Core i5, 8 GB 1867 MHz DDR3, with Docker for Mac configured with 4 CPUs, 4 GiB of memory and 512 MiB of swap.

![fortio benchmark graphic 2019-08-21](benchmark/fortio/2019-08-27-fluentd_tcp_x_sock.png?raw=true "2019-08-21 fluentd tcp x fluentd sock")

Commands:

```shell
$ docker-compose -f compose.bench.yml up

$ fortio load \
    -c 4 -qps 50000 -t 30s -a -labels "eventhorizon stdout" \
    -content-type application/json \
    -H "Ce-Custom-A: Foo" \
    -H "Ce-Custom-B: Bar" \
    -H "Ce-Id: BenchmarkHTTPClient" \
    -H "Ce-Source: myapp" \
    -H "Ce-Specversion: 0.3" \
    -H "Ce-Subject: MyMethod.MyAction" \
    -H "Ce-Time: 2019-08-20T22:18:27.166904Z" \
    -H "Ce-Type: io.request.rpc" \
    http://localhost:1257

$ fortio load \
    -c 4 -qps 50000 -t 30s -a -labels "eventhorizon fluentd tcp" \
    -content-type application/json \
    -H "Ce-Custom-A: Foo" \
    -H "Ce-Custom-B: Bar" \
    -H "Ce-Id: BenchmarkHTTPClient" \
    -H "Ce-Source: myapp" \
    -H "Ce-Specversion: 0.3" \
    -H "Ce-Subject: MyMethod.MyAction" \
    -H "Ce-Time: 2019-08-20T22:18:27.166904Z" \
    -H "Ce-Type: io.request.rpc" \
    http://localhost:1247

$ fortio load \
    -c 4 -qps 50000 -t 30s -a -labels "eventhorizon fluentd sock" \
    -content-type application/json \
    -H "Ce-Custom-A: Foo" \
    -H "Ce-Custom-B: Bar" \
    -H "Ce-Id: BenchmarkHTTPClient" \
    -H "Ce-Source: myapp" \
    -H "Ce-Specversion: 0.3" \
    -H "Ce-Subject: MyMethod.MyAction" \
    -H "Ce-Time: 2019-08-20T22:18:27.166904Z" \
    -H "Ce-Type: io.request.rpc" \
    http://localhost:1237
```
