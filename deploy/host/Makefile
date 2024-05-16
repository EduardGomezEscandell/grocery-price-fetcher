PKG := grocery-price-fetcher
ETC := /etc/$(PKG)
DATA := /var/lib/$(PKG)
SERVICE := /etc/systemd/system/$(PKG).service
NOW := $$(date --rfc-3339 seconds)
DOCKER_USER ?= "NO-USER-DEFINED"

.PHONY: start stop purge

start:
	sudo mkdir -p "$(ETC)"
	sudo mkdir -p "$(DATA)"

	sudo cp ./service/$(PKG).service "$(SERVICE)"
	sudo cp ./compose.yaml "$(ETC)/compose.yaml"
	sudo cp ./data/* "$(DATA)"

	sudo systemctl daemon-reload
	sudo systemctl enable "$(PKG).service"
	sudo systemctl start "$(PKG).service"

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
.PHONY: package

VERSION := $(shell git rev-parse HEAD | cut -c 1-8)
package:
	@echo "Building package $(PKG) version $(VERSION)"

	sudo docker tag $(PKG) $(DOCKER_USER)/$(PKG):$(VERSION)
	sudo docker push $(DOCKER_USER)/$(PKG):$(VERSION)

	rm -r grocery-price-fetcher/ || true
	mkdir -p $(PKG)
	rsync -r ./* $(PKG) --exclude $(PKG) --exclude .gitignore

	# Remove package target from this very Makefile
	sed -i -n '/# DO NOT DEPLOY AFTER THIS LINE/q;p' $(PKG)/Makefile

	# Patch compose.yaml to pull image from docker hub
	sed -i 's#image: $(PKG)#image: docker.io/$(DOCKER_USER)/$(PKG):$(VERSION)#g' $(PKG)/compose.yaml

	echo $(VERSION) > "$(PKG)/version"
	tar -czf $(PKG).tar.gz $(PKG)