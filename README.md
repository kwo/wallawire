# Wallawire

## Local

### Build

    make

### Run

Once:

    mkdir -p walladata/db
    ./scripts/create-certs.sh
    cockroach start --host=0.0.0.0 --port=5432 --http-port=8080 --store=path=./walladata/db --certs-dir=./walladata/certs/dbserver 
    ./scripts/create-database.sh
    Ctrl-c

thereafter:

    cockroach start --host=0.0.0.0 --port=5432 --http-port=8080 --store=path=./walladata/db --certs-dir=./walladata/certs/dbserver 
    . .testenv
    go run main.go

optionally start the ui in dev mode

    cd ui
    yarn start


## Docker

### Build

Once:

    docker build -t ww/nodejs:latest deployments/nodejs
    docker build -t ww/golang:latest deployments/golang
    docker build -t ww/db:latest deployments/cockroach

thereafter:

    docker build -t ww/ui:latest --build-arg VERSION=$(git describe --always --tags --dirty="*") ui
    docker build -t wallawire:latest .

### Run

Once:

    mkdir -p walladata/db
    ./scripts/create-certs.sh
    docker-compose up db -d
    ./scripts/create-database.sh
    docker-compose down

thereafter:

    docker-compose up -d
    docker-compose logs -f
    ...
    docker-compose down
