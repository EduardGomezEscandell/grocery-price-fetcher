# Deployment
To deploy this project you need an Ubuntu machine.

## Local deployment
Go to the repository root and run:
```
make full-start
```

## Remote deployment

1. Build the latest version of the project. Start a shell at the repository root and run:
    ```
    make build-docker
    ```
    After this, the docker image will be built and called `grocery-price-fetcher`.
2. Push the image and build the package
    1. Once again in the repository root, run:
        ```
        make package DOCKER_USER="YOUR_DOCKER_HUB_USERNAME"
        ```
        This will do two things:
        - Upload the image to dockerhub.
        - Create a tarball located at `.deploy/host/grocery-price-fetcher.tar.gz`.
              This tarball has the uploaded image URL hardcoded
3. Transfer the tarball to the deployment machine using the method of your choice.
4. Log into the deployment machine, go to the tarball, and run:
   ```
   tar -xvf grocery-price-fetcher.tar.gz
   cd grocery-price-fetcher
   make start
   ```

## Cleanup

```[!WARNING]
The service's persistent data will be located at `/var/lib/grocery-price-fetcher`. Make sure to make a copy before cleaning up!
```

## Local deployment
Go to the repository root and run `make clean`.

## Remote deployment
- In the build machine, run `make clean`
- In the deployment machine, go to unpacked package and run `make purge`.
  