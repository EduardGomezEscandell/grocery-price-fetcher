PKG := grocery-price-fetcher
ETC := /etc/$(PKG)
DATA := /var/lib/$(PKG)
SERVICE := /etc/systemd/system/$(PKG).service
NOW := $$(date --rfc-3339 seconds)

.PHONY: start stop purge

start:
	sudo mkdir -p "$(ETC)"
	sudo mkdir -p "$(DATA)"

	sudo cp ./service/$(PKG).service "$(SERVICE)"
	sudo cp ./compose.yaml "$(ETC)/compose.yaml"
	sudo cp ./data/* "$(DATA)"

	sudo systemctl daemon-reload
	sudo systemctl enable "$(PKG).service"
	sudo systemctl restart "$(PKG).service"

	systemctl status $(PKG).service --no-pager -l
	printf "\\nTo view logs run:\\n$ journalctl -u $(PKG).service --since '$(NOW)'\\n"

stop:
	sudo systemctl stop "$(PKG).service"    || true
	sudo systemctl disable "$(PKG).service" || true

purge: stop
	sudo rm -rf "$(SERVICE)"  || true
	sudo rm -rf "$(ETC)" || true
	sudo rm -rf "$(DATA)" || true

# DO NOT DEPLOY AFTER THIS LINE
.PHONY: package deploy

VERSION := $(shell git rev-parse HEAD | cut -c 1-8)
package:
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

deploy: package
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
		$(GCLOUD_VM)					\
		--command="cd ~ && tar -xzf $(PKG).tar.gz && cd $(PKG) && make start"

	@echo "Allowing server to start..."
	sleep 10
	
	curl --fail http://$$(                \
		gcloud compute instances describe \
		  $(GCLOUD_VM)                    \
		  --project=$(GCLOUD_PROJECT)     \
		  --zone=$(GCLOUD_ZONE)           \
		  | grep natIP                    \
		  | sed 's#^.*natIP: \(.*\)#\1#'  \
		)/api/helloworld
