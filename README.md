# kvstore

## Overview

This project provides a key-value store service with an HTTP API interface. 
The service is built with Go and Tarantool, using JWT for authentication.

## Installation

1. **Clone the repository:**

```shell
git clone https://github.com/deadshvt/kvstore.git
```

2. **Go to the project directory:**

```shell
cd kvstore
```

3. **Install dependencies:**

```shell
go mod tidy
```

## Running the application

1. **Run containers**

```shell
make run
```

2. **Stop and remove containers**

```shell
make down
```

3. **Check logs**

```shell
make logs
```

## API

1. **Authentication**

Description: authenticate user and receive JWT token

Endpoint: `/api/login`

Method: `POST`

Headers:
- `Content-Type`: `application/json`

Request Body:

```shell
{
  "username": "your_username",
  "password": "your_password"
}
```

Response Codes:
- 200 OK
- 400 Bad Request
- 401 Unauthorized
- 500 Internal Server Error

Example:

```shell
curl -X POST http://localhost:8080/api/login \
-H "Content-Type: application/json" \
-d '{
  "username": "admin",
  "password": "presale"
}'
```

2. **Write data**

Description: write key-value pairs to the store

Endpoint: `/api/write`

Method: `POST`

Headers:
- `Content-Type`: `application/json`
- `Authorization`: `Bearer your_generated_token`

Request Body:

```shell
{
  "data": {
    "key1": "value1",
    "key2": "value2",
    "key3": 1
  }
}
```

Response Codes:
- 200 OK
- 400 Bad Request
- 401 Unauthorized
- 500 Internal Server Error

Example:

```shell
curl -X POST http://localhost:8080/api/write \
-H "Authorization: Bearer your_generated_token" \
-H "Content-Type: application/json" \
-d '{
  "data": {
    "key1": "value1",
    "key2": "value2",
    "key3": 1
  }
}'
```

3. **Read data**

Description: read key-value pairs from the store

Endpoint: `/api/read`

Method: `POST`

Headers:
- `Content-Type`: `application/json`
- `Authorization`: `Bearer your_generated_token`

Request Body:

```shell
{
  "keys": ["key1", "key2", "key3"]
}
```

Response Codes:
- 200 OK
- 400 Bad Request
- 401 Unauthorized
- 500 Internal Server Error

Example:

```shell
curl -X POST http://localhost:8080/api/read \
-H "Authorization: Bearer your_generated_token" \
-H "Content-Type: application/json" \
-d '{
  "keys": ["key1", "key2", "key3"]
}'
```
