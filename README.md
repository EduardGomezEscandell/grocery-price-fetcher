# Grocery price fetcher

This repository contains a web application that allows you to pick the weekly meals, introduce the food you have in the pantry,
and compute your shopping list and its cost.

The purpose is for me to learn, so don't expect any high-quality commercial product.

## Architecture
> [!WARNING]
> There's no guarantee this section is up-to-date.

- The backend is a Go service serving static endpoints (for the frontend) and dynamic endpoints for the API.
- The database is a simple out-of-the box MySQL inside of a Docker volume.
- The frontend is a basic React application.

The product runs on three containers:
- `database` hosts the MySQL database.
- `prepopulator` runs on start-up only, it detects if the database is empty and pre-fills it with sample data.
- `grocery` runs the http server.

## Deployment
Check out [this guide](./deploy/README.md) to see how to deploy the server.
