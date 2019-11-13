# This file should be included into main Makefile !
ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

e2e:
	$(call log_info,Starting test environment:)
	cd $(ROOT_DIR) && docker-compose up -d --build nginx redis kiwi-iam
	$(call log_info,Test environment started!)
	make e2e/venom

e2e/nobuild: ## Run tests without rebuilding IAM
	$(call log_info,Starting test environment:)
	cd $(ROOT_DIR) && \
	docker-compose up -d --build nginx redis && \
	docker-compose up -d --no-build kiwi-iam
	$(call log_info,Test environment started!)
	make e2e/venom

e2e/venom:
	sh $(ROOT_DIR)/e2e-venom.sh

e2e/env-stop: ## Stop test environment
	cd $(ROOT_DIR) && docker-compose down -v
