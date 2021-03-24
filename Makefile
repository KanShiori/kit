SHELL := /bin/bash
GO := go
DOCKER := docker
PWD := $(shell pwd)


.PHONY: lint
lint:
	golangci-lint run