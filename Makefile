#########################################
## Preamble
SHELL := $(shell which bash)
.SHELLFLAGS := -eu -o pipefail -c
.ONESHELL:             ;   # Recipes execute in same shell
.NOTPARALLEL:          ;   # Wait for this target to finish
.SILENT:               ; 	 # No need for @
.EXPORT_ALL_VARIABLES: ;   # Export variables to child processes.
.DELETE_ON_ERROR:

# Modify the block character to be `-\t` instead of `\t`
ifeq ($(origin .RECIPEPREFIX), undefined)
  $(error This Make does not support .RECIPEPREFIX. Please use GNU Make 4.0 or later)
endif
.RECIPEPREFIX = -

MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules


#########################################
## Setup
# Assign the default goal
.DEFAULT_GOAL := help
default: $(.DEFAULT_GOAL)
all: help


#########################################
##  Logging
# Enable Logging
#exec 1> >(tee -a make.log) 2>&1
# ifeq ($(.DEFAULT_GOAL),)
#   $(warning no default goal is set)
# endif
# $(warning default goal is $(.DEFAULT_GOAL))


#########################################
## Variables
PROJECT_DIR := $(shell pwd)
PYTHON_DIR 	:= $(PROJECT_DIR)/platform/python
COCKPIT_DIR := $(PROJECT_DIR)/platform/cockpit
GO_DIR 			:= $(PROJECT_DIR)/src
BUILD_DIR 	:= $(PROJECT_DIR)/build

OSARCH := "linux/amd64 linux/386 windows/amd64 windows/386 darwin/amd64 darwin/386"
ENV = /usr/bin/env

GO ?= go
PYTHON ?= /usr/bin/env python3 -u -B


#########################################
## Help Command
.PHONY: help
help: ## List commands
-	echo -e "USAGE: make \033[36m[COMMAND]\033[0m\n"
-	echo "Available commands:"
-	awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\t\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)


#########################################
##
PHONY: pyhton
python:    ## Run the Sensible Pyhton Module
-	echo -e "\033[36mRunning Sensible...\033[0m"
- $(PYTHON) $(PYTHON_DIR)/sensible.py --dir=$(PROJECT_DIR)/tests


#########################################
##
PHONY: go
go:    ## Run the Sensible Go Module (go run)
-	echo -e "\033[36mRunning Sensible...\033[0m"
- cd $(GO_DIR)
# - go mod tidy
# -	$(GO) get
# -	$(GO) run *.go
- $(GO) run ./

PHONY: dep
dep:    ## Get build dependencies
-	echo -e "\033[36mGrabbing Dependencies...\033[0m"
- cd $(GO_DIR)
-	go get

PHONY: build
build:    ## Build the app
-	echo -e "\033[36mBuilding Sensible...\033[0m"
- cd $(GO_DIR)
- dep ensure && go build -o $(BUILD_DIR)/ ./main.go

# cross-build:   ## Build the app for multiple os/arch
# -	echo -e "\033[36mBuilding Sensible...\033[0m"
# - cd $(GO_DIR)
# -	gox -osarch=$(OSARCH) -output "bin/blackbeard_{{.OS}}_{{.Arch}}"


#########################################
##
.PHONY: yaml-lint
yaml-lint: ## Run YAML Lint, Ansible Lint, and remove DOS line endings.
-	echo "Starting the Linting Process"
-	echo "Fixing line endings..."
-	find . -type f -exec sed -i "s|\r$\||" {} \;
-	echo "Running Linter..."
-	find $(ANSIBLE_DIR) -type f -name "*.yml" -exec ansible-lint --force-color -p {} \;
-	find $(ANSIBLE_DIR) -type f -name "*.yml" -exec yamllint -f colored {} \;


#########################################
##
.PHONY: clean
clean: ## Clean up
-	echo "Cleaning..."
-	find ./ -name '*.pyc' -exec rm -f {} \;
-	find ./ -name 'Thumbs.db' -exec rm -f {} \;
-	find ./ -name '*~' -exec rm -f {} \;
