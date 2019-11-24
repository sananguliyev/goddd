# GoDDD (Recipes API)
[![Build Status](https://travis-ci.org/SananGuliyev/goddd.svg?branch=master)](https://travis-ci.org/SananGuliyev/goddd)

This repository is sample recipe service which comes with CRUD and rating the recipes endpoints. This sample is just for giving an idea how to implement the DDD in Go.

## Implementation Notes
[Domain Driven Design (DDD)](https://en.wikipedia.org/wiki/Domain-driven_design) is an approach which is followed for this sample. It's open to discussion whether we need it or not because there is no super complexity in our domain but it's good approach for abstraction. 

In this example we have 3 layers:
* **Application** - Main entry point which contains endpoints (or commands), middleware, etc.
* **Domain** - Encapsulated business logic
* **Infrastructure** - Database, 3rd party sources, etc.

## Packages

* `github.com/vektra/mockery` (installing on `DockerfileTest`) which helps to auto-generate mock files based on the interfaces
* `github.com/go-pg/pg`: Entity manager
* `github.com/google/uuid`: Generating and inspecting UUIDs based on [RFC 4122](http://tools.ietf.org/html/rfc4122) and DCE 1.1
* `github.com/gorilla/handlers`: For logging HTTP requests in the [Apache Combined Log Format](http://httpd.apache.org/docs/2.2/logs.html#combined)
* `github.com/gorilla/mux`: HTTP router and URL matcher 
* `github.com/tbaud0n/go-rql-parser`: For translating [RQL (Resource Query Language)](https://dundalek.com/rql/draft-zyp-rql-00.html) queries to SQL 
* `github.com/stretchr/testify`: Assertions and mock generation

## Setup
You can just run `docker-compose up -d` and wait a couple of seconds that app became ready to handle the requests (or you can run `docker-compose up` and watch the logs to know when the app is ready)

## Test
For running the tests
```
docker-compose run test
```
Test service is using `DockerfileTest` which is generating mock files and running all tests with coverage. 

Furthermore, Do not forget to run `flyway` service up before running the tests because we need a database with tables for integration tests

![#f03c15](https://placehold.it/15/f03c15/000000?text=+) P.S. The mock files are not added to the pull request to keep the repository clean. You can use the commands for generating mock files:
```
$ go get github.com/vektra/mockery/.../
$ cd /path/to/project
$ mockery -dir domain/repository -name RecipeRepository
$ mockery -dir domain/repository -name RatingRepository
$ mockery -dir domain/interactor -name RecipeInteractor
$ mockery -dir domain/tool -name IdGenerator
```

## Database migration
[Flyway](https://flywaydb.org) is used as a database migration tool. It helps us to not stick on one language because whenever we have the decision to switch to another language then we do not need to care about database migration.

Furthermore, The 1st initial migration script is added to the repository which creates the tables and fills some dummy data.

## Authorization
For making `create`, `update`, `delete` endpoints protected I implemented just hardcoded `GoDDD` token but in current application design it is easy to switch to some other verification. E.g. [JWT](https://jwt.io/).

## Presenter
`gorilla/mux` is used for routing of the HTTP handlers. Due to application architecture, it's easy to implement other servers, e.g. [gRPC](https://grpc.io/).

## Filtering
[Resource Query Language (RQL)](https://dundalek.com/rql/draft-zyp-rql-00.html) is used for filtering. Supported operators:
* scalar operators
    * `=`: equal
    * `=ne=`: not equal
    * `=like=`: contain
    * `=match=`: contain (case-insensitive)
    * `=gt=`: greater
    * `=lt=`: less
    * `=ge=`: greater equal
    * `=le=`: less equal
* logic operators
    * `,` or `&`: and
    * `;` or `|`: or 
    
For `like` and `match` operators you should use `*` instead of `%`. E.g. `name=match=*goddd*` translates to `name ILIKE '%goddd%'`

## Examples

### Create
```
curl -X "POST" "http://localhost/recipes" \
     -H 'Authorization: GoDDD' \
     -d $'{
  "name": "Dolma",
  "difficulty": 4,
  "is_vegetarian": false,
  "prepare_time": 50
}'
```

### Update 
```
curl -X "PUT" "http://localhost/recipes/{recipe_id}" \
     -H 'Authorization: GoDDD' \
     -d $'{
  "name": "Dolma",
  "difficulty": 4,
  "is_vegetarian": false,
  "prepare_time": 45
}'
```

### Delete
```
curl -X "DELETE" "http://localhost/recipes/{recipe_id}" \
     -H 'Authorization: GoDDD'
```

### Get
```
curl "http://localhost/recipes/0d93bda0-f040-4564-967b-e59bf5571dcd"
```
P.S This endpoint returns one extra `rating` field which is an overall rating of the recipe. E.g. in case of 2 ratings (`5` & `4`) it will return `4.5`

### Rate
```
curl -X "POST" "http://localhost/recipes/{recipe_id}/rating" \
     -d $'{
  "value": 5
}'
```

### List
```
curl "http://localhost/recipes/?filter={conditions}&page={page_number}&limit={recipe_per_page}"
```

You can pass 3 (optional) query params to the recipe list endpoint
* `filter`: basic **RQL** filtering. E.g. `prepare_time==15;name=like=*kebab*` will return the recipes names which contains `%kebab%` and prepare time is `15`
* `limit`: for deciding the count of the recipes per page (default: `10`)
* `page`: current page number of pagination (default: `1`)

## Missing parts
Only important parts are covered by tests. Some of them are not fully covered but easy to cover.

* Not %100 covered by unit tests
* Not %100 covered by integration tests
