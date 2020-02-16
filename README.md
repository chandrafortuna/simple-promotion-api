# Simple Promotion API

Promotion Service API. Create Promo, Calculate promo price, etc

Design database and flowchart in this folder `/_hotel-management-system-design`

Postman Collection: `Promotion API.postman_collection.json`

## Getting Started

First, Make sure you have set up \$GOPATH.

```bash
# Download this project
go get github.com/chandrafortuna/simple-promotion-api

```

### Running the application

Navigate to `simple-promotion-api` folder and build and run it:

```
cd $GOPATH/src/github.com/chandrafortuna/simple-promotion-api
go build -o promo-api
./promo-api
```

Application will running on http://localhost:8000/promo

The following table shows the HTTP methods and URLs that represent the action supported in the API.

| Request  | Description |
| ------------- | ------------- |
| `GET /promo`  | Show a list of available promo  |
| `POST /promo`  | Create new promo  |
| `POST /promo/apply`  | Apply promotion for list of room price |
| `POST /promo/distribute`  | Distribute promo quota |


For example request, please import postman collection in this repository

## Running the tests

```
go test ./test/
```
