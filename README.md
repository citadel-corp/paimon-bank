# citadel-corp/paimon-bank

API for simple Banking App service.

🔍 Tested with - TBA.

📝 Documentation - TBA.

🎵 Songs to test by - [playlist](https://open.spotify.com/album/1oVSp3g7ULNAHzFtdBvHEd?si=IVw3cdo6RUKDOdb1gYCJKQ).

## Getting Started

These instructions will give you a copy of the project up and running on
your local machine for development and testing purposes.

### Prerequisites

Requirements for the software and other tools to run and test
- [Go](https://go.dev/doc/install)
- [PostgreSQL](https://www.postgresql.org/download/)
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [k6](https://k6.io/docs/get-started/installation/) - to test

For monitoring system:
- [Docker](https://www.docker.com/get-started/)
- [docker-compose](https://docs.docker.com/compose/install/) - for orchestration
- [Prometeus](https://prometheus.io/docs/prometheus/latest/installation/#using-docker)
- [Grafana](https://grafana.com/docs/grafana/latest/setup-grafana/installation/docker/)


Note that we use [AWS S3 service](https://aws.amazon.com/s3/) to upload image,
[setup your own](https://docs.aws.amazon.com/AmazonS3/latest/userguide/GetStartedWithS3.html) S3 bucket if you want to test uploading image.

### Migrate the database

After [setting up your database locally](https://www.postgresql.org/docs/current/tutorial-createdb.html),
run this to migrate our database structure
```
$ migrate -source file://path/to/migrations -database postgres://localhost:5432/database up 2
```

### Running the service

A step by step series of that tell you how to get a development
environment running

Create file named .env and fill it with, example:
```
DB_URL = 'host=host.docker.internal port=5432 user=your-username password=your-password dbname=paimon-bank sslmode=disable'
DB_HOST =  host.docker.internal
DB_PORT =  5432
DB_USERNAME =  ${DB_USERNAME}
DB_PASSWORD =  ${DB_PASSWORD}
DB_NAME =  paimon-bank
DB_PARAMS =  '&sslmode=disable'

BCRYPT_SALT =  12
JWT_SECRET =  ${JWT_SECRET}

S3_ID =  ${S3_ID}
S3_SECRET_KEY =  ${S3_SECRET_KEY}
S3_BUCKET_NAME =  ${S3_BUCKET_NAME}
S3_REGION = ${S3_REGION}

ENV = development
```

Run the service

    $ docker-compose up -d

Now you can run the service.  
Service is running on port `8080`.  
Prometheus is running on port `9090`.  
Grafana is running on port `3000`.  

## Endpoints

Main service running on port `8080`.

- User
    - Register - `POST /v1/user/register`
    - Login - `POST /v1/user/login`
- Balance
    - Add - `POST /v1/balance`
    - List - `GET /v1/balance`
    - History - `GET /v1/balance/history`
- Transaction
    - Create - `POST /v1/transaction`
- Image
    - Upload - `POST /v1/image`
- Prometheus
    - Metrics - `/metrics`
    - Health - `/healthz`

## Monitoring system

Open the now available grafana dashboard http://localhost:3000/dashboards.

## Running the tests

TBA.

## Authors

The [Citadel Corp](https://github.com/citadel-corp) team:
  - [**TheSagab**](https://github.com/TheSagab)
  - [**Faye**](https://github.com/farolinar)

## License

This project is licensed under the [MIT License](https://github.com/citadel-corp/paimon-bank?tab=MIT-1-ov-file) - see the [LICENSE](https://github.com/citadel-corp/paimon-bank/blob/main/LICENSE) file for
details

## Acknowledgments

  - The Ramadhan ProjectSprint organizer and members
