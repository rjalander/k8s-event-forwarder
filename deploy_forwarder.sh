#!/bin/sh

docker build -t localhost:5000/cdevents/k8s-event-forwarder:latest .
docker push localhost:5000/cdevents/k8s-event-forwarder:latest

kubectl apply -f deployment.yaml

kubectl apply -f event-serviceaccount.yaml
kubectl apply -f event-role.yaml
kubectl apply -f event-rolebinding.yaml

