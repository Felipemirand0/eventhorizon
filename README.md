# EventHorizon

An open source [CloudEvents](https://cloudevents.io) implementation to allow handling a massive number of events.

## Introduction

EventHorizon provides a uniform way for handling events ([CloudEvents](https://cloudevents.io)) in a Kubernetes native way, by implementing the following resources:

- **[Singularity](samples/kubernetes/Singularity.yml)** - describes the instance itself, the transport method (currently supports HTTP or NATS).
- **[CloudEventHandler](samples/kubernetes/CloudEventHandler.yml)** - describes event handler itself, the encoder and the output.
- **[CloudEventValidator](samples/kubernetes/CloudEventValidator.yml)** - and optional validator, in case you want to enforce specific events.
- **[CloudEventOutput](samples/kubernetes/CloudEventOutput.yml)** - describes where to send the event output to (currently supports Fluentd).

## Current status

Under active development, any contributions are welcome.

## Installing

Clone this repository:

```shell
git clone https://github.com/acesso-io/eventhorizon.git

cd eventhorizon/
```

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

Download the binaries from the [GitHub releases page](https://github.com/acesso-io/eventhorizon/releases) to `dist/` folder.

Set those environment variables:

```shell
export EVENTHORIZON_MODE=standalone
export EVENTHORIZON_NAME=eventhorizon
export EVENTHORIZON_STANDALONE_CONFIG=$PWD/samples/standalone/stdout.yml
export EVENTHORIZON_LOGGING_LEVEL=info
export EVENTHORIZON_LOGGING_PRETTY=true
```

Start EventHorizon:

```shell
./dist/eventhorizon
```
