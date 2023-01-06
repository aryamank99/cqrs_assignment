Tenant Microservice Assignment:

Create a microservice using CQRS and Event Sourcing design pattern / architecture
Basic Tenant Attributes – Name (Should be unique), admin username / password, admin email, license type
(FULL / SELECTIVE), total users, flag to enable / disable, license start date and license end date
List of APIs to expose (https://restfulapi.net/resource-naming/):
  ● Get All tenants
  ● Get a particular tenant
  ● Create a tenant
  ● Update tenant – Separate APIs for
  ● Updating admin credentials
  ● Enabling / Disabling tenant
  ● Updating licensing information
  
Core language and libraries – Golang (Gorilla Mux, GORM)
Database – MySQL (RDBMS) / MongoDB (NoSQL)
Message Broker – RabbitMQ
IDE – Visual Studio Code
Tool to test API – Postman
