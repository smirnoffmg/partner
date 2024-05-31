
# Go Telegram Bot Template

This is a template for creating Telegram bots using Go. It includes a Dockerfile and docker-compose configuration for easy deployment.

## Features

- Simple and clean project structure
- Dockerized for easy deployment
- Example configuration file

## Prerequisites

- Go 1.22 or later
- Docker
- Docker Compose

## Getting Started

### Clone the Repository

```sh
git clone https://github.com/smirnoffmg/go-telegram-bot-template.git
cd go-telegram-bot-template
```

### Configuration

1. Copy `.env.example` to `.env` and update the values as needed.

### Building and Running with Docker

#### Build and Start the Application

```sh
docker-compose up --build -d
```

#### Stop the Application

```sh
docker-compose down
```

### Running Locally

#### Build the Application

```sh
go build -o main .
```

#### Run the Application

```sh
./main
```

## Project Structure

```plaintext
.
├── Dockerfile
├── README.md
├── config.yaml
├── docker-compose.yml
├── go.mod
├── go.sum
└── main.go
```

- `Dockerfile`: Docker configuration for building the Go application.
- `docker-compose.yml`: Docker Compose configuration for running the application.
- `config.yaml`: Configuration file for the application.
- `go.mod` and `go.sum`: Go module files managing project dependencies.
- `main.go`: Entry point for the Go application.
- `internal/`: Directory for logic - interfaces, usecases, etc.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## Contact

If you have any questions or need further assistance, feel free to contact the repository owner.
