PKG := grocery-price-fetcher
ETC := /etc/$(PKG)
SERVICE := /etc/systemd/system/$(PKG).service

NOW := $(shell date +"%Y-%m-%d %H:%M:%S")
NOW_UTC := $(shell date +"%Y-%m-%d %H:%M:%S" --utc)

.PHONY: install start stop purge help

help: ## Show this help message
	@echo "Grocery Price Fetcher Deployment Makefile"
	@echo ""
	@echo "COMMANDS"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install the service
	sudo mkdir -p "$(ETC)"

	sudo cp ./service/$(PKG).service "$(SERVICE)"
	sudo cp ./compose.yaml "$(ETC)/compose.yaml"
	sudo cp ./config.yaml "$(ETC)/config.yaml"
	sudo cp -r ./sample-data "$(ETC)/sample-data"

	[ -f "$(ETC)/db_root_password.txt" ]                              \
		&& echo "Root password file already exists"                   \
		|| (                                                          \
			sudo touch         "$(ETC)/db_root_password.txt"          \
			&& sudo chmod 0600 "$(ETC)/db_root_password.txt"          \
			&& echo "$$(openssl rand -hex 32)"                        \
			   | sudo tee "$(ETC)/db_root_password.txt" > /dev/null   \
			&& echo "Root password file created"                      \
		)
	
	[ -d "$(ETC)/certificates" ]                              		  \
		&& echo "Certificate directory already exists"                \
		|| (                                                          \
			sudo mkdir -p      "$(ETC)/certificates"                  \
			&& sudo chmod 0600 "$(ETC)/certificates"                  \
			&& sudo cp ./certificates/* "$(ETC)/certificates"         \
			&& echo "Self-signed certificates used"                   \
		)

	sudo systemctl daemon-reload

	@echo ""
	@echo "Service installed. To start:"
	@echo "  1. Set your password in $(ETC)/db_root_password.txt"
	@echo "  2. Set your certificates in $(ETC)/certificates"
	@echo "  3. Validate or change the configuration in $(ETC)/config.yaml and $(ETC)/compose.yaml"
	@echo "     You may also do this after starting the service, but changes will not take effect until the service is restarted"
	@echo "  4. Run command 'make start' to start the service"
	@echo ""

start: ## Start the service
	sudo systemctl daemon-reload
	sudo systemctl enable "$(PKG).service"
	sudo systemctl restart "$(PKG).service"
	systemctl status "$(PKG).service" --no-pager -l
	@echo To view logs run:
	@echo
	@echo journalctl -u \"$(PKG).service\" --since \"$(NOW)\"
	@echo

db-shell: ## Open a shell in the database container
	sudo docker exec \
		-it $(PKG)-database-1 \
		bash -c 'mysql -u root -p$$(cat /run/secrets/db_root_password)'

stop: ## Stop the service
	sudo systemctl stop "$(PKG).service"    || true
	sudo systemctl disable "$(PKG).service" || true

purge: stop ## Stop and remove all service files and data
	sudo docker compose rm -fsv          || true
	sudo docker container rm -f $(PKG)-grocery-1   || true
	sudo docker container rm -f $(PKG)-database-1   || true
	sudo docker container rm -f $(PKG)-prepopulate-1   || true
	sudo docker volume rm -f $(PKG)_data || true
	sudo rm -rf "$(SERVICE)"             || true
	sudo rm -rf "$(ETC)" 	             || true

# DO NOT DEPLOY AFTER THIS LINE
.PHONY: package deploy

VERSION := $(shell git describe --tags --always --dirty)
package: ## Build a deployable package (requires DOCKER_USER)
	@echo "Building package $(PKG) version $(VERSION)"
	if [ -z "$(DOCKER_USER)" ]; then echo "Error: Environment variable DOCKER_USER not set" >&2; exit 99; fi

	sudo docker tag $(PKG) $(DOCKER_USER)/$(PKG):$(VERSION)
	sudo docker push $(DOCKER_USER)/$(PKG):$(VERSION)

	rm -r $(PKG) / || true
	mkdir -p $(PKG)
	rsync -r ./* $(PKG) --exclude $(PKG) --exclude .gitignore

	# Remove package target from this very Makefile
	sed -i -n '/# DO NOT DEPLOY AFTER THIS LINE/q;p' $(PKG)/Makefile

	# Patch compose.yaml to pull image from docker hub
	sed -i 's#image: $(PKG)#image: docker.io/$(DOCKER_USER)/$(PKG):$(VERSION)#g' $(PKG)/compose.yaml

	echo $(VERSION) > "$(PKG)/version"
	tar -czf $(PKG).tar.gz $(PKG)

deploy: package ## Deploy the package to a remote server (requires GCLOUD_VM, GCLOUD_ZONE, GCLOUD_PROJECT)
	if [ -z "$(GCLOUD_VM)"      ]; then echo "Error: Environment variable GCLOUD_VM not set" >&2;      exit 99; fi
	if [ -z "$(GCLOUD_ZONE)"    ]; then echo "Error: Environment variable GCLOUD_ZONE not set" >&2;    exit 99; fi
	if [ -z "$(GCLOUD_PROJECT)" ]; then echo "Error: Environment variable GCLOUD_PROJECT not set" >&2; exit 99; fi

	gcloud compute scp                  \
		$(PKG).tar.gz                   \
		$(GCLOUD_VM):~/$(PKG).tar.gz     \
		--project='$(GCLOUD_PROJECT)'   \
		--zone='$(GCLOUD_ZONE)'

	gcloud compute ssh                  \
		--project='$(GCLOUD_PROJECT)'   \
		--zone='$(GCLOUD_ZONE)'         \
		$(GCLOUD_VM)                    \
		--command="cd ~ && tar -xzf $(PKG).tar.gz && cd $(PKG) && make install && make start"

	@echo "Allowing server to start..."
	sleep 10

	@echo "Deployment logs:"
	gcloud compute ssh                  \
		--project='$(GCLOUD_PROJECT)'   \
		--zone='$(GCLOUD_ZONE)'         \
		$(GCLOUD_VM)                    \
		--command="journalctl -u '$(PKG).service' --since '$(NOW_UTC)' --no-pager -l"
	
	@echo "Deployment sanity-test..."

	URL="https://$$(                      \
		gcloud compute instances describe \
		  $(GCLOUD_VM)                    \
		  --project=$(GCLOUD_PROJECT)     \
		  --zone=$(GCLOUD_ZONE)           \
		  | grep natIP                    \
		  | sed 's#^.*natIP: \(.*\)#\1#'  \
		)/api/version";                   \
	curl --retry 20 --retry-delay 5 --retry-all-errors -k "$$URL"

	
