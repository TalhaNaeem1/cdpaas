# pipelineService

pipelineService lets you easily extract and load data from and to databases, APIs and file formats. Ready-to-use connectors are provided that adjust automatically as schemas and APIs evolve. Each of your pipelines can be readily monitored, refreshed as needed, and schedule the updates.



**Install golangci-linter**

`go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.40.1`

**Install Swag Tool**

`go get -u github.com/swaggo/swag/cmd/swag@v1.6.7`

**Update Swagger Documentation**

From pipelineService folder execute

`swag init`

**Access Swagger Documentation**

`http://<SERVER_HOST>:<SERVER_PORT>/api/v1/swagger/index.html`

**Execute Unit Tests**

From pipelineService folder execute

`go test ./...`

**Environment Variables**
```
DB_USERNAME=<DB_USERNAME>
DB_PASSWORD=<DB_PASSWORD>
DB_NAME=<DB_NAME>
DB_HOST=<DB_HOST>
DB_PORT=<DB_PORT>
BUILD_ENV=dev
SERVER_PORT=<SERVER_PORT>
AIRBYTE_HOST=<AIRBYTE_HOST>
AIRBYTE_PORT=<AIRBYTE_PORT>
```