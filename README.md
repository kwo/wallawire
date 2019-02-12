# Wallawire

## Build

### Local Build

    make

### Docker Build

Once:

    docker build -t ww/nodejs:latest deployments/nodejs
    docker build -t ww/golang:latest deployments/golang
    docker build -t ww/db:latest deployments/cockroach

thereafter:

    docker build -t ww/ui:latest -f Dockerfile.ui .
    docker build -t wallawire:latest .

## Run

Once:

    mkdir -p walladata/db
    ./scripts/create-certs.sh
    make up
    ./scripts/create-database.sh
    make down

thereafter:

    make up
    make logs
    ...
    make down
    
Teardown and reinitialize database:

    ./scripts/reinit-database.sh
