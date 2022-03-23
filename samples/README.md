# Chat Samples

This directory contains multiple samples on how to build and use this chat app in multiple scenarios

## Pre Requisites
### Install RabbitMQ

Have [the main prerequisites]() ready
Install the [RabbitMQ Cluster Operator]() and [Topology Operator]() to build this samples, this because they are using RabbitMQ's Brokers and Triggers.

To install simply run:
```bash
kubectl apply -f https://github.com/rabbitmq/cluster-operator/releases/latest/download/cluster-operator.yml
kubectl apply -f https://github.com/jetstack/cert-manager/releases/latest/download/cert-manager.yaml
kubectl wait --for=condition=Ready pods --all -n cert-manager
kubectl apply -f https://github.com/rabbitmq/messaging-topology-operator/releases/latest/download/messaging-topology-operator-with-certmanager.yaml
```

## Run samples

To run this samples simply run:
```bash
ko apply -Rf samples/$SAMPLE_NAME
```
In the case of the multi-namespace-topology you need to wait for both of the brokers to be ready before starting the chat services for the events to be sent correctly 


## Cleanup

To clean any of the examples just run:

```bash
ko delete -Rf samples/$SAMPLE_NAME
```
