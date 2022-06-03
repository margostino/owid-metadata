#!make
SHELL = /bin/sh
.DEFAULT: help

-include .env .env.local .env.*.local

# Defaults
DOCKER_COMPOSE = USERID=$(shell id -u):$(shell id -g) docker-compose ${compose-files}
ALL_ENVS := local ci
env ?= local

.PHONY: help
help:
	@echo "Lumos pipeline"
	@echo ""
	@echo "Usage:"
	@echo "  vercel.dev                     - Replicate the Vercel deployment environment locally, allowing you to test your Serverless Functions, without requiring you to deploy each time a change is made"
	@echo "  run                            - Start main runner for local testing"
	@echo ""
	@echo "Project-level environment variables are set in .env file:"
	@echo "  VERCEL=1"
	@echo "  VERCEL_ENV=development"
	@echo ""
	@echo "Note: Store protected environment variables in .env.local or .env.*.local"
	@echo ""

.PHONY: clean.metadata
clean.metadata:
	rm ./metadata/*

