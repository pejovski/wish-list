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

### Dependency 
- Event Bus - RabbitMQ [Common](https://github.com/pejovski/common)
- Catalog API - [Catalog](https://github.com/pejovski/catalog)

### Usage
- Make sure the shared RabbitMQ container is up and running [Common](https://github.com/pejovski/common)
- Make sure [Catalog API](http://localhost:8201) is active from [Catalog](https://github.com/pejovski/catalog)
```bash
docker-compose up -d
go run main.go
```
- open [Wish List API](http://localhost:8203)
- play!

## Swagger update
- use http://editor.swagger.io
- modify app/swagger/swagger.yaml
- run: statik -src=./app/swagger -dest=./app

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