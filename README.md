# Linnworks x Square Integration

This is a purpose built tool (for a good friend of mine) that allows inventory and order syncing between Linnworks and Square. It also provides a simple dashboard for extra visibility.

As the name suggests, `go-lsi` has been built with Golang.

## Features
* Syncing of Inventory and Services from Linnworks to Square,
* Syncing of Orders from Square to Linnworks,
* A recurring job that starts syncing every 10',
* Simple dashboard for admins to trigger syncing and view inventory and orders,
* A RESTful API for programmatic access,

## API
These are the RESTful endpoints that the app exposes for programmatic access. Endpoints denoted with 🔒, are protected and they require a `Bearer` token.
```
🔒 GET /api/v1/ping
🔒 GET api/v1/inventory - Params (at least one required): ?sku=string&barcode=string
🔒 GET /api/v1/orders
🔒 GET /api/v1/sync/status
🔒 POST /api/v1/sync/recent
🔒 POST /api/v1/sync - Body: {"start": "datetime", "end": "datetime"}
POST /api/v1/auth - Body: {"username": "string", "password": "string"}
```

A Postman collection is provided at `files/postman`.

## Usage

### Environement variables
The following environment variables are required to run the app
```
PORT=int
DB=path-to-db
SIGNING_KEY=a-long-string-to-sign-jwts
LINNWORKS_APP_ID=string
LINNWORKS_APP_SECRET=string
LINNWORKS_APP_TOKEN=string
SQUARE_ACCESS_TOKEN=string
SQUARE_HOST=string
SQUARE_API_VERSION=string
SQUARE_LOCATION_ID=string
SQUARE_TEAM_MEMBER_IDS=comma,separated,list,of,user,ids
SYNC_INTERVAL=non-zero-int
```

### Run with docker
A `Dockerfile` and a `Makefile` are provided for easy start. The `make run` command expects a `.env` file to be present with the aforementioned variables.

```bash
make build
make run
```

### Run from source
Export the required environment variables and the you can use `go run` to start the app:
```
export PORT=8080
...
go run main.go
```

## License
Since this built to address specific needs and business cases, it might not be useful to others. 

However, feel free to use any parts of this if you need to 😄.
