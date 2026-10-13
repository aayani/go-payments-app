.DEFAULT_GOAL := bootstrap

TF = docker run --rm -v "$(CURDIR)/infra:/workspace" -w /workspace hashicorp/terraform:latest
TF_VARS = -var="aws_endpoint=http://host.docker.internal:4566"

.PHONY: bootstrap up down

bootstrap: up

up:
	@test -f .env || cp .env.example .env
	docker compose -f compose-infra.yaml up -d --wait
	docker build -t go-payments-app:local .
	$(TF) init -input=false
	@if ! $(TF) state list 2>/dev/null | grep -qx 'aws_db_instance.payments'; then \
		instance=$$(docker compose -f compose-infra.yaml exec -T ministack aws --endpoint-url http://localhost:4566 rds describe-db-instances --query "DBInstances[?DBInstanceIdentifier=='payments'].DBInstanceIdentifier" --output text) || exit 1; \
		if [ "$$instance" = "payments" ]; then \
			$(TF) import -input=false $(TF_VARS) aws_db_instance.payments payments || exit 1; \
		fi; \
	fi
	$(TF) apply -input=false -auto-approve $(TF_VARS)

down:
	docker compose -f compose-infra.yaml down