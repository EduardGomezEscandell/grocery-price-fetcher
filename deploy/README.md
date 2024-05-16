# Deployment
To deploy this project you need an Ubuntu machine.

**Local deployment**

You can deploy the service with one single command. Go to the repository root and run:
```
make full-start
```

**Remote deployment on Google Cloud Compute**

1. Build the latest version of the project. Start a shell in the repository root and run:
    ```
    make build-docker
    ```
    After this, the docker image will be built and called `grocery-price-fetcher`.
2. Install the Google Cloud SDK and log in. See more: https://cloud.google.com/sdk/docs/install
3. Log into Docker. See more: https://docs.docker.com/reference/cli/docker/login/
4. Push the image, build the package, and deploy to Google cloud in a single command, in the repository root, run:
    ```
    make deploy GCLOUD_VM="?" GCLOUD_ZONE="?" GCLOUD_PROJECT="?" DOCKER_USER="?"
    ```
    Replace the `?` with appropiate data. 
5. You're done

**Remote deployment on other tenants**

1. Build the latest version of the project. Start a shell in the repository root and run:
    ```
    make build-docker
    ```
    After this, the docker image will be built and called `grocery-price-fetcher`.
2. Log into Docker. See more: https://docs.docker.com/reference/cli/docker/login/   
3. Push the image and build the package. In the repository root, run:
    ```
    make package DOCKER_USER="YOUR_DOCKER_HUB_USERNAME"
    ```
    This will do two things:
    - Upload the image to dockerhub.
    - Create a tarball located at `.deploy/host/grocery-price-fetcher.tar.gz`.
          This tarball has the uploaded image URL hardcoded
4. Transfer the tarball to the deployment machine using the method of your choice.
5. Log into the deployment machine, go to the tarball, and run:
   ```
   tar -xvf grocery-price-fetcher.tar.gz
   cd grocery-price-fetcher
   make start
   ```

## Cleanup

> [!WARNING]
> The service's persistent data will be located at `/var/lib/grocery-price-fetcher`. Make sure to make a copy before cleaning up!

**Local deployment**

Go to the repository root and run `make clean`.

**Remote deployment on Google Cloud**

- In the build machine, run `make clean`
- Log into the development VM instance and run `cd ~/grocery-price-fetcher && make purge`

**Remote deployment on other tenants**

- In the build machine, run `make clean`
- In the deployment machine, go to the unpacked package and run `make purge`.
  
