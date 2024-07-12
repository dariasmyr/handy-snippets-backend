<div style="text-align: center;">
    <img src="docs/logo.png"
         alt="Logo"
         style="width: 50%; height: auto;" />
</div>

# Handy Snippets Frontend

## About The App

[**Jolly Secrets**](https://jollysecrets.uxna.me) - an app for storing and sharing end-to-end encrypted notes online. Written on [**Go**](https://go.dev/), SQLite DB and [**GraphQL**](https://graphql.org/) API to interact with [**Handy Snippets Frontend**](https://github.com/dariasmyr/handy-snippets-frontend)

## How to Run
Сlone the repo:
```bash
go get github.com/dariasmyr/handy-snippets-backend
```

### Option 1: Run via Command Line
You can run the application directly from the command line. Do not forget to enter your PORT and FRONTEND_URL (for CORS policy) into env variables.
```bash
PORT=<your_port> FRONTEND_URL=<your_frontend_url>  go run cmd/server.go 
```

### Option 2: Run via Docker
Run the application using Docker by executing the following command. FO NOT forget to write your env variables in the `.env` file:
```bash
# Launch
$ docker-compose up -d

# Rebuild and launch
$ docker-compose up -d --build
```

## Database workflow
This app uses **SQLite** database for data size optimisation and embedding.

