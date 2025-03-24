SHELL := /bin/bash

GOCMD=go
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install

BINARY_NAME=task-management-system

run-task:
	mkdir -p storage/
	go run ./cmd/task-ms task

run-authn:
	mkdir -p storage/
	go run ./cmd/task-ms authn

.PHONY: run-task run-authn