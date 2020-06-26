#! /bin/bash

# build web services

cd ./api
env GOOS=linux GOARCH=amd64 go build -o ../bin/api

cd ..

cd ./scheduler
env GOOS=linux GOARCH=amd64 go build -o ../bin/scheduler

cd ..

cd ./streamserver
env GOOS=linux GOARCH=amd64 go build -o ../bin/streamserver

cd ..

cd ./web
env GOOS=linux GOARCH=amd64 go build -o ../bin/web

cd ..


