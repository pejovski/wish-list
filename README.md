# Wish List

### Requirements

```bash
go version 1.13
docker
docker-compose
golang-statik: sudo apt install golang-statik
```

### Installation

Use [go mod](https://blog.golang.org/using-go-modules) to install dependencies.

```bash
go mod tidy
```

Run [docker-compose](https://docs.docker.com/compose/) to build docker images and run necessary containers.

```bash
docker-compose up -d
```

### Swagger
- use http://editor.swagger.io
- modify app/swagger/swagger.yaml
- run: statik -src=./app/swagger -dest=./app

### Usage

```bash
go run main.go -race
```

# Architecture and Design

The project code follows the design principles from the resources bellow

### Microsoft Micro-Services

https://docs.microsoft.com/en-us/dotnet/architecture/microservices/index

### Uber Go Code Structure

https://www.youtube.com/watch?v=nLskCRJOdxM
- extended with receivers and emitters for working with events

### Rest HTTP Server by Go veteran

https://www.youtube.com/watch?v=rWBSMsLG8po

## License
[MIT](https://choosealicense.com/licenses/mit/)