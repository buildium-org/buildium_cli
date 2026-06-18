# Golang Starter Template

This repository provides starter files to help you run your code against the Buildium test harness.

## What's Included

- **Dockerfile** — A base image that extends the test harness image. Copy your app binary into it and run the harness against your code.
- **docker-compose.yml** — Lets you define and manage not only your app service, but also any additional services required by the tutorial, such as Postgres, Redis, or Elasticsearch.
- **Makefile** — Handles building your project binary, copying over `meta.json`, and running the Docker image to test your solution.

## Getting Started

1. Replace the `FROM` directive in the Dockerfile with the base test harness image provided by the tutorial author.
2. In the `Dockerfile`, update the image name to your desired image name.
3. Update the `docker-compose.yml` to use your locally built image.
4. Update `meta.json` with your project ID from the Buildium website. When you're ready, set the `stage` value to the stage you want to test against.
5. Remove the current Git remote and add a new remote pointing to your desired GitHub repository.
