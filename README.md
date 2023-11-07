# file-search-api

## Structure

- main.go - the entrypoint + REST API
- internal/logger - logger implementation
- internal/inputfile - input file loader
- internal/repository - the repo to search indices (in-memory file implementation added, more implementations with the same interface can be added, like DB or Redis)
- api - openapi documentation
- build - dockerfile

Please see code comments for more details.

> **Remember that code structure matters**
> 
> The folder structure loosely follows https://github.com/golang-standards/project-layout 
>
> The main.go is not placed inside cmd/ since this is a simple app and we only have one main.go, thus cmd folder adds complexity but not the value
> 
> main.go is responsible for initialization, router configuration and running of the app. For prod ready app I would
> create a kind of App struct to hold dependencies and to provide server running mechanism (to be able to run it from functional tests)
> and also I would extract HTTP handlers to a separate layer. But hopefully for this size of app / test task it is more than enough.

## Run the app

Put the `input.txt` into the project root.

Simply run
```shell
$ docker-compose up
```

Services exposed:

> REST API: http://localhost:8080/index/{value}

> Docs: http://localhost:8088

## Configuration file

The app is highly configurable via the `.env` file (if run via `docker compose`).

Also, it is configurable by setting `PORT`, `FILENAME`, `MARGIN_PERCENT` and `LOGGER_PREFIX` directly in the environment.

> I haven't added any dotenv file load/processing, nor any yaml configuration files, cause that would add more dependencies to the project,
> which can be easily avoided since our app is quite simple.

## Algorithm

Since the task requires to put the file into a slice (and not a tree, or a db), the best algo for such case is a binary search algorithm (O logN)

## Tests

`internal/` has 93.1% coverage. This can be improved, but 93.1 should be good enough for the purposes of the test task, because all the business logic is covered.

Run all tests:
```shell
$ make test
```

## Ideas

- pass the margin via the endpoint to be able to configure it per-request (e.g. GET /index/100?marginPercent=10)
- add caching decorator for the repo (to cache frequently requested values)
