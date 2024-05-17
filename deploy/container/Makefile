ROOT := root
PKG := grocery-price-fetcher

.PHONY: help build clean

help: ## Show this help message
	@echo "Grocery Price Fetcher Container Build Makefile"
	@echo ""
	@echo "COMMANDS"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Create directories, copy files, and build Docker image
	mkdir -p "$(ROOT)/usr/bin"
	mkdir -p "$(ROOT)/etc/$(PKG)"
	mkdir -p "$(ROOT)/usr/share/$(PKG)/frontend"

	# Config
	cp "config.yaml"                 "$(ROOT)/etc/$(PKG)/config.yaml"

	# Artefacts
	cp    "../../bin/grocery-server" "$(ROOT)/usr/bin/$(PKG)"
	cp -R "../../frontend/build/."   "$(ROOT)/usr/share/$(PKG)/frontend/"

	# Build
	sudo docker build . -t $(PKG)

clean: ## Remove build output, stop & remove the Docker container, and remove the image
	rm -rf "$(ROOT)" || true
	sudo docker stop $$(sudo docker container ls -a | grep $(PKG) | cut -c 1-8) || true
	sudo docker rm $$(sudo docker container ls -a | grep $(PKG) | cut -c 1-8) || true
	sudo docker rmi $(PKG) || true
	exit 0

