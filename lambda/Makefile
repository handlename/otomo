BUILD_DIR := build
TFSTATE_PATH := ../terraform/terraform.tfstate
LAMBROLL_FUNCTION := function_url.jsonnet
LAMBROLL_FUNCTION_URL := function.jsonnet

export ROLE_ARN := $(shell cat $(TFSTATE_PATH) | jq -r '.resources[] | select(.type == "aws_iam_role") | select(.name == "lambda") | .instances[0].attributes.arn')
export GOOS := linux
export GOARCH := arm64

dryrun: build
	lambroll deploy \
	    --src $(BUILD_DIR) \
        --dry-run

deploy: build
	lambroll deploy $\
        --src $(BUILD_DIR)

diff:
	lambroll diff $(LAMBROLL_FLAGS)

destroy:
	lambroll delete

build: clean
	mkdir -p $(BUILD_DIR)
	$(MAKE) $(BUILD_DIR)/app
	$(MAKE) $(BUILD_DIR)/bootstrap
	$(MAKE) $(BUILD_DIR)/config.yaml

clean:
	rm -rf $(BUILD_DIR)

$(BUILD_DIR)/app:
	go build -o $@ .

$(BUILD_DIR)/%: %
	cp $< $@
