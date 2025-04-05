BUILD_DIR := build
TMP_DIR := tmp
TFSTATE_PATH := ../terraform/terraform.tfstate
LAMBROLL_FUNCTION := function_url.jsonnet
LAMBROLL_FUNCTION_URL := function.jsonnet
SSMWRAP_VERSION := 2.2.2

export SERVICE_NAME := otomo
export AWS_REGION := us-east-1
export ROLE_ARN := $(shell cat $(TFSTATE_PATH) | jq -r '.resources[] | select(.type == "aws_iam_role") | select(.name == "lambda") | .instances[0].attributes.arn')
export GOOS := linux
export GOARCH := arm64

dryrun: build
	lambroll deploy \
	    --src $(BUILD_DIR) \
		--function-url function_url.jsonnet \
        --dry-run

deploy: build
	lambroll deploy \
        --src $(BUILD_DIR) \
        --function-url function_url.jsonnet

diff:
	lambroll diff $(LAMBROLL_FLAGS)

destroy:
	lambroll delete

logs:
	lambroll logs --follow

build: clean
	mkdir -p $(BUILD_DIR)
	$(MAKE) $(BUILD_DIR)/app
	$(MAKE) $(BUILD_DIR)/bootstrap
	$(MAKE) $(BUILD_DIR)/config.toml
	$(MAKE) $(BUILD_DIR)/ssmwrap

clean:
	rm -rf $(BUILD_DIR) $(TMP_DIR)

$(BUILD_DIR)/%: %
		cp $< $@

$(BUILD_DIR)/app:
	go build -o $@ ../cmd/otomo/

$(BUILD_DIR)/config.toml: ../config.toml
	cp $< $@

$(BUILD_DIR)/ssmwrap: $(TMP_DIR)/ssmwrap_$(SSMWRAP_VERSION)/ssmwrap
	cp $< $@

$(TMP_DIR)/ssmwrap_$(SSMWRAP_VERSION)/ssmwrap:
	mkdir -p $(TMP_DIR)/ssmwrap_$(SSMWRAP_VERSION)
	cd $(TMP_DIR)/ssmwrap_$(SSMWRAP_VERSION) && \
		curl -sL https://github.com/handlename/ssmwrap/releases/download/v$(SSMWRAP_VERSION)/ssmwrap_$(SSMWRAP_VERSION)_linux_arm64.tar.gz \
		| tar xzf -
