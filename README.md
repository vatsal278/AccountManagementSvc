# AccountManagmentSvc

[![Build](https://github.com/vatsal278/AccountManagmentSvc/actions/workflows/build.yml/badge.svg)](https://github.com/vatsal278/AccountManagmentSvc/actions/workflows/build.yml) [![Test Cases](https://github.com/vatsal278/AccountManagmentSvc/actions/workflows/test.yml/badge.svg)](https://github.com/vatsal278/AccountManagmentSvc/actions/workflows/test.yml) [![Codecov](https://codecov.io/gh/vatsal278/AccountManagmentSvc/branch/main/graph/badge.svg)](https://codecov.io/gh/vatsal278/AccountManagmentSvc)

* This service was created using Golang.
* This service has used clean code principle and appropriate go directory structure.
* This service is completely unit tested and all the errors have been handled.
* This service utilises messageBroker service for communicating with other micro services.

## Starting the UserManagementService

* Start the Docker container for mysql with command :
```
docker run --rm --env MYSQL_ROOT_PASSWORD=pass --env MYSQL_DATABASE=accmgmt --publish 9085:3306 --name mysqlDb -d mysql
```
* Start the MsgBroker service using steps as described in the [link](https://github.com/vatsal278/msgbroker)


* Start the Api locally with command :
```
go run .\cmd\AccountManagmentSvc\main.go
```
### You can test the api using post man, just import the [Postman Collection](./docs/accountMgmtSvc.postman_collection.json) into your postman app.
### To check the code coverage
```
cd docs
go tool cover -html=coverage
```
## Account Management Service:

This application is split up into multiple components, each having a particular feature and use case. This will allow individual scale up/down and can be started up as micro-services.

HTTP calls are made across micro-services.

They are made asynchronous & de-coupled via pub-sub or messaging queues.

*For testing individual services, these can be via direct HTTP calls*


All requests & responses follow json encoding.
Requests are specific to the concerned endpoint while responses are of the following json format & specification:
>
>    Response Header: HTTP code
>
>    Response Body(json):
>    ```json
>    {
>       "status": <HTTP status code>,
>       "message": "<message>",
>       "data": {
>        // object to contain the appropriate response data per API
>       }
>    }
>    ```
# AccManagementSvc-Design

## AccManagementSvc Endpoints

## Create Account
This endpoint is used to create a new account in a Relational DB with `account_number` as the `primary key` and as a `int`,`user_id`as `string` and it should be `unique`, `income` and `spends` as `decimal`,`created_at` and `updated_at` as `timestamp` `active_services` and `inactive_services` as `json` and once its done we will notify the user mgmt svc through msg queue. This endpoint will get hit from usermanagement service through msg queue.

DbSchema:
```
+-------------------+---------------+------+-----+-------------------+-----------------------------------------------+
| Field             | Type          | Null | Key | Default           | Extra                                         |
+-------------------+---------------+------+-----+-------------------+-----------------------------------------------+
| user_id           | varchar(255)  | NO   | UNI | NULL              |                                               |
| account_num       | int           | NO   | PRI | NULL              | auto_increment                                |
| income            | decimal(18,2) | YES  |     | 0.00              |                                               |
| spends            | decimal(18,2) | YES  |     | 0.00              |                                               |
| created_at        | timestamp     | NO   |     | CURRENT_TIMESTAMP | DEFAULT_GENERATED                             |
| updated_at        | timestamp     | NO   |     | CURRENT_TIMESTAMP | DEFAULT_GENERATED on update CURRENT_TIMESTAMP |
| active_services   | json          | NO   |     | json_object()     | DEFAULT_GENERATED                             |
| inactive_services | json          | NO   |     | json_object()     | DEFAULT_GENERATED                             |
+-------------------+---------------+------+-----+-------------------+-----------------------------------------------+
```

## Account Summary
A user hits this endpoint in order to view the details of their account. They can view details like the services that are available and services that theywant to subscribe to.
There will be jwt token containing userid in cookie
#### Specification:
Method: `GET`

Path: `/account`

Request Body: `not required.`

Success to follow response as specified:

Response Header: HTTP 200

Response Body(json):
```json
{
  "status": 200,
  "message": "SUCCESS",
  "data": {
    "account_number": "<full account number>",
    "income": <income calculated based on all incoming transactions> as float,
    "spends": <spends calculated based on all outgoing transactions> as float,
    "active_services": ["<list of all services that user has subscribed to>",],
    "available_services": ["<list of all services that user has not subscribed to but are available for subscription>",]
  }
}
```

#### Specification:
Method: `POST`

Path: `/account`

Request Body:
```json
{
  "user_id": "<user_id for the record to activate>"
}
```

Success to follow response as specified:

Response Header: HTTP 200

Response Body(json):
```json
{
  "status": 201,
  "message": "SUCCESS",
  "data": nil
}
```

## Latest Transaction
This endpoint receives the details of latest transaction and it updates the income and spends column in db according to type of transaction
#### Specification:
Method: `PUT`

Path: `/account/update/transaction`

Request Body:
```json
{
  "account_number":"acc_no."
  "amount": "amount of the transaction",
  "type":"debit or credit"
}
```

Success to follow response as specified:

Response Header: HTTP 200

Response Body(json):
```json
{
  "status": 202,
  "message": "SUCCESS",
  "data": nil
}
```

## Update services
This endpoint updates the active service and available service column acc to query
#### Specification:
Method: `PUT`

Path: `/account/update/transaction`

Request Body:
```json
{
  "account_number":"acc_no.",
  "service_id": "id of service that needs to be removed or added",
  "type":"add or remove"
}
```

Success to follow response as specified:

Response Header: HTTP 200

Response Body(json):
```json
{
  "status": 202,
  "message": "SUCCESS",
  "data": nil
}
```

## AccManagementSvc Middlewares

1. ExtractUser: extracts the user_id from the cookie passed in the request and forwards it in the context for downstream processing.
2. ScreenRequest: allows requests only from the message queue to be passed downstream. The middleware checks the “`user-agent`” & request `URL` to identify requests originating from the message queue.
   *The URL(s) of the message queue(s) is passed as a configuration to the service to allow requests only from URLs in the list*.
3. Caching middleware
