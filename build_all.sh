#!/bin/bash

GOOS=linux GOARCH=arm go build -o bin/gofindsub-linux-arm cmd/main.go
GOOS=linux GOARCH=amd64 go build -o bin/gofindsub-linux-amd64 cmd/main.go
GOOS=linux GOARCH=386 go build -o bin/gofindsub-linux-368 cmd/main.go

GOOS=windows GOARCH=arm go build -o bin/gofindsub-win-arm cmd/main.go
GOOS=windows GOARCH=amd64 go build -o bin/gofindsub-win-amd64 cmd/main.go
GOOS=windows GOARCH=386 go build -o bin/gofindsub-win-368 cmd/main.go

GOOS=darwin GOARCH=amd64 go build -o bin/gofindsub-darwin-amd64 cmd/main.go
GOOS=darwin GOARCH=arm64 go build -o bin/gofindsub-darwin-amr64 cmd/main.go