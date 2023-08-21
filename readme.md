# Mini Wallet Project for JULO

Mini Wallet is a Go microservice to deposit and withdraw money from an account. Made for the JULO recruitment process for the Back End Developer postition.

## Installation

This project needs Go and Sqlite. To get its dependencies, install the packages from the go.mod

```bash
go mod download
```

This project requires the CGO_ENABLED variable from the GOENV to be true
```bash
go env -w CGO_ENABLED=1
```
## Usage
to start the server, run this command:
```bash
go run main.go
```
the port defaults at [http://localhost:9000](http://localhost:9000)
## Requests
All the endpoints are tested using Postman. You can export the collection [here](https://api.postman.com/collections/19615785-a41a178b-5c81-42e4-9d5d-805d781d058b?access_key=PMAT-01H8C3SQKEW0NN5YQ6HS3378C1)