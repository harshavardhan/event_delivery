## Intro

- This repo implements a POC for event delivery with the following components:
    - Runs a http server on port 8090. Exposes `receive_event` endpoint to receive event data
    - Redis as datastore (assumes that redis is up and running on localhost at port 6379. If different, use ``REDIS_URL`` env variable)
    - Starts a mock consumer processing events from the datastore
    - Also initiates a mock producer producing some initial events and sending them to receive events endpoint

## Running instructions

- To build an executable binary run ``go build -o main``
- To execute the binary run ```./main``` (Ensure that redis is up and running before this)

## Docker compose

- If you wish to run both application and redis as docker containers in a single command, run ``docker-compose up``

## Tests

- Tests are written only for the critical redis package which takes care of storing and processing incoming events. Run tests using ``go test -v -cover ./redis
``