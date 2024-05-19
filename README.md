# Grocery price fetcher

This repository contains a web application that allows you to pick the weekly meals, introduce the food you have in the pantry,
and compute your shopping list and its cost.

The purpose is for me to learn, so don't expect any high-quality commercial product.

## Architecture
> [!WARNING]
> There's no guarantee this section is up-to-date.

- The backend is a Go service serving static endpoints (for the frontend) and dynamic endpoints for the API.
- The database is just a bunch of JSON files inside of a Docker volume. I'll make it into a proper database eventually.
- The frontend is a basic React application.

## Deployment
Check out [this guide](./deploy/README.md) to see how to deploy the server.
