# Logical Expression Evaluator

## Assumptions

The current implementation of the logical expression evaluator assumes that:

1. the given expression characters (and logical operators) must only be separated by spaces ` `;
2. the given expression must contain only lowercase variables, such as `x AND y`;
3. the given expression must be a valid logical expression (expression validation is not yet implemented -- see [TODOs section](#todos)).

## Running the application

### Required applications

- make
- docker:
  - docker-compose engine if running an older version of docker. You may also need to change the `docker compose` commands throughout the application to `docker-compose`.

Run `make server` to spin the database and server. The server runs at port 8080.

### Stopping the application

To stop the application and clean all resources, run `make server-down`.

## Running tests

There are three commands to run the application tests:

- `make unit-test`: runs only unit tests;
- `make integration-test`: runs only integration tests;
- `make test`: runs both unit and integration tests.

Integration tests will run in their own containers.

## API

The full API documentation can be seen through the swagger document at `internal/docs/generated/swagger.json`. You can also run directly from the application by running `make swag-serve`. It will install `go-swagger` through brew if you don't have it.

There is also a postman collection at `internal/docs/postman/exp-eval-postman-collection.json` for convenience with all the main requests.

There are some important mentions regarding the challenge's requirements that need to be made:

1. a `v1` prefix was used to version the endpoints, so, the endpoint `/evaluate/{expression_id}?x=1,y=0,z=1` described in the requirements is at `v1/evaluate/{expression_id}?x=1,y=0,z=1`. This rule applies to every endpoint;
2. the create and update endpoints were designed separately, instead of a POST at `v1/expressions` to perform both depending on the payload, we leverage the HTTP methods to differentiate between these two actions. `POST` to create and `PATCH` to update and expression.
3. some additional endpoints were added to the API, such as:
   1. GET `v1/expressions/{id}`: get expression by ID;
   2. DELETE `v1/expressions/{id}`: delete expression by ID;
   3. PATCH `v1/expressions`: updates an existing expression.
4. a

The API requires authentication for all endpoints. There are two mocked bearer tokens that can be used to run requests against the server:

```
token 1: 74edf612f393b4eb01fbc2c29dd96671
token 2: d88b4b1e77c70ba780b56032db1c259b
```

Each token is tied to a specific user, and there are operations in the system that can only be performed by the author of the expression e.g.: updating/deleting/getting by id.

### Some Sample Requests/Responses for each endpoint:

#### Create expression

Request:
```
curl --location --request POST 'http://localhost:8080/v1/expressions' \
--header 'Authorization: Bearer 74edf612f393b4eb01fbc2c29dd96671' \
--header 'Content-Type: application/json' \
--data-raw '{
    "expression": "((x OR y) AND (z OR k) OR j)"
}'
```

Response:
```
{
    "expressionID": "659f9c60-9056-4ba2-ae12-bb70533d8671",
    "expression": "((x OR y) AND (z OR k) OR j)",
    "username": "John Doe",
    "createdAt": "2023-02-05T18:03:59.41586Z",
    "updatedAt": "2023-02-05T18:03:59.41586Z"
}
```

#### Get expression by ID

Request:
```
curl --location --request GET 'http://localhost:8080/v1/expressions/659f9c60-9056-4ba2-ae12-bb70533d8671' \
--header 'Authorization: Bearer 74edf612f393b4eb01fbc2c29dd96671'
```

Response:
```
{
    "expressionID": "659f9c60-9056-4ba2-ae12-bb70533d8671",
    "expression": "((x OR y) AND (z OR k) OR j)",
    "username": "John Doe",
    "createdAt": "2023-02-05T18:03:59.41586Z",
    "updatedAt": "2023-02-05T18:03:59.41586Z"
}
```

#### List expressions

List expressions endpoint was built with support for pagination by sending `page_id` and `page_size` as query parameters. PageID must be greater than 0. Pagination information (page number, page items and total number of pages will be added to the header response)

Request:
```
curl --location --request GET 'http://localhost:8080/v1/expressions?page_id=1&page_size=2' \
--header 'Authorization: Bearer 74edf612f393b4eb01fbc2c29dd96671'
```

Response:
```
[
    {
        "rowID": 1,
        "expressionID": "659f9c60-9056-4ba2-ae12-bb70533d8671",
        "expression": "((x OR y) AND (z OR k) OR j)",
        "username": "John Doe",
        "createdAt": "2023-02-05T18:03:59.41586Z",
        "updatedAt": "2023-02-05T18:03:59.41586Z"
    },
    {
        "rowID": 2,
        "expressionID": "11625fa6-cb11-491d-97fa-089fa94d43b5",
        "expression": "((x OR y) AND (z OR k) OR j)",
        "username": "John Doe",
        "createdAt": "2023-02-05T18:05:44.774641Z",
        "updatedAt": "2023-02-05T18:05:44.774641Z"
    }
]
```

#### Delete expression

Request:
```
curl --location --request DELETE 'http://localhost:8080/v1/expressions/11625fa6-cb11-491d-97fa-089fa94d43b5' \
--header 'Authorization: Bearer 74edf612f393b4eb01fbc2c29dd96671'
```

Response: A success response will return a 204 (status no content) http status.

#### Evaluate expression

Request:
```
curl --location --request GET 'http://localhost:8080/v1/evaluate/11625fa6-cb11-491d-97fa-089fa94d43b5?X=1&y=0&z=1&k=0&j=1' \
--header 'Authorization: Bearer 74edf612f393b4eb01fbc2c29dd96671' \
--header 'Content-Type: application/json' \
--data-raw '{
    "expression": "((x OR y) AND (z OR k) OR j)"
}'
```

Response:
```
{
    "result": true
}
```


## TODOs

The logical expression evaluator method used is not complete, leading to some misleading results for some edge cases.
The following must still be implemented:

1. validate a given logical expression before saving/evaluating its value, so we know that we will always have reliable expressions; 
2. implement support for `!` (NOT) operator;
3. take the operators' order of precedence into account when performing an evaluation;
