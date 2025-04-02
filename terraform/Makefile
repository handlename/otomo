.PHONY: plan
plan: init
	bash -c 'set -o pipefail; terraform plan -parallelism=90 | egrep -v "(Refreshing state...|Reading...|Read complete after)"'

.PHONY: apply
apply: init
	terraform apply -parallelism=90

.PHONY: init
init:
	terraform init

.PHONY: lock
lock:
	terraform init -upgrade
	terraform providers lock \
		-platform=linux_amd64 \
		-platform=linux_arm64 \
		-platform=darwin_amd64 \
		-platform=darwin_arm64
