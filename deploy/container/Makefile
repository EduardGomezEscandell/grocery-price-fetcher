ROOT := root
PKG := grocery-price-fetcher

.PHONY: build clean

build:
	mkdir -p "$(ROOT)/usr/bin"
	mkdir -p "$(ROOT)/etc/$(PKG)"
	mkdir -p "$(ROOT)/usr/share/$(PKG)/frontend"

	# Config
	cp "config.yaml"                 "$(ROOT)/etc/$(PKG)/config.yaml"

	# Artefacts
	cp    "../../bin/grocery-server" "$(ROOT)/usr/bin/$(PKG)"
	cp -R "../../frontend/build/."   "$(ROOT)/usr/share/$(PKG)/frontend/"

	# Build
	sudo docker build . -t grocery-price-fetcher

clean:
	rm -rf "$(ROOT)" || true
	sudo docker stop $$(sudo docker container ls -a | grep grocery-price-fetcher | cut -c 1-8) || true
	sudo docker rm $$(sudo docker container ls -a | grep grocery-price-fetcher | cut -c 1-8) || true
	sudo docker rmi grocery-price-fetcher || true
	exit 0

