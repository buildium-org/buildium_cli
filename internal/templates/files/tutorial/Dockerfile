FROM ubuntu:latest

# Set environment variables to non-interactive for apt
ENV DEBIAN_FRONTEND=noninteractive

# Install git and golang if required, you can add more dependencies as needed
RUN apt-get update && apt-get install -y git

# Install dependencies
RUN apt-get install -y build-essential golang

# Make working directory and copy local repo contents
WORKDIR /app/harness
COPY . .

RUN make build