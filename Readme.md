# Chat System API

This project is a chat system API developed in Go, using Cassandra as the database. The API allows users to register, log in, send messages, and retrieve conversation history between two users.

## Requirements

- [Go](https://golang.org/doc/install) 1.21 or higher
- [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/)
- [Swag](https://github.com/swaggo/swag) to generate Swagger documentation


## Defaults
  - Port: [8080]
  - Gin Gonic
  - docker y docker-compose.
  - go 1.21

## Endpoints

### Users

- **Register User**: `POST /users/register`
  - Registers a new user in the system.

- **Login**: `POST /users/login`
  - Verifies a user's credentials and allows them to log in.

- **List Users**: `GET /users`
  - Returns a list of all registered users.

### Messaging

- **Send Message**: `POST /messages/send`
  - Sends a message from one user to another. (Requires authentication)

- **Conversation History**: `GET /messages/history/:senderID/:receiverID`
  - Retrieves the message history between two users. (Requires authentication)

### System Health

- **Ping**: `GET /ping`
  - Checks the availability of the system.

## Usage

The Makefile included in this project provides convenient commands for quickly managing common tasks. Below is an explanation of each command available:

### `Start`

```bash
make start
```
This command starts all the services defined in the docker-compose.yml file in detached mode

### `Stop`

```bash
make start
```
This command stops and removes all the containers defined in the docker-compose.yml file. This is useful when you want to completely tear down the environment, including the removal of containers, networks, and volumes that were created.

### `Coverage`
```bash
make coverage
```
This command executes the coverage.sh script to generate a coverage report for the tests. The script typically runs the tests and then produces a report showing which parts of the code were exercised during the tests.

### Create/Update Swagger

1. Run the following command to install the latest version of Swagger:
    ```sh
    go install github.com/swaggo/swag/cmd/swag@latest
    ```

2. Execute the command:
    ```sh
    make build-doc
    ```

3. If you encounter the following error:
    ```
    /bin/sh: swag: command not found
    ```
   You should execute the following command and then repeat step 2:
    ```sh
    export PATH=$PATH:$(go env GOPATH)/bin
    ```

Once the Swagger documentation is generated, access it by opening a web browser and navigating to:
   `<host>/swagger/index.html`

For example, in a local environment, the URL would typically be:
   `http://localhost:8080/swagger/index.html`
