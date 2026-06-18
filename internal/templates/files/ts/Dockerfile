FROM <TEST_HARNESS_IMAGE_HERE>:latest

# Install Node.js and npm
RUN apt-get update && apt-get install -y nodejs npm

WORKDIR /app
COPY . .

RUN npm install
RUN make build

CMD ["/app/harness/bin/app"]
