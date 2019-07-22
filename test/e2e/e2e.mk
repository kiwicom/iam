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
	$(call log_info,Starting venom tests:)
	cd $(ROOT_DIR) && docker-compose build venom
	$(eval VENOM_CONTAINER_NAME = venom-container-$(shell date +%s))
	cd $(ROOT_DIR) && docker-compose run -d --name $(VENOM_CONTAINER_NAME) venom tail -f /dev/null
	docker exec $(VENOM_CONTAINER_NAME) venom run --var-from-file variables.yml --parallel 5 --format=xml --output-dir="." --strict ; \
	let "EXIT_CODE=$$?" ; \
	cd $(ROOT_DIR) && docker exec  $(VENOM_CONTAINER_NAME) cat test_results.xml > test_results.xml ; \
	docker-compose rm --force --stop -v venom || true ; \
	exit $$EXIT_CODE
	$(call log_success,Venom integration tests completed!)

e2e/env-stop: ## Stop test environment
	cd $(ROOT_DIR) && docker-compose down -v
