# API
Hi, and welcome to the bs-books API. The API is written in Go. 

## Technologies
The API depends on a number of technologies: 

### Docker
The site can be run locally via Docker. Each deployable (in the cmd directory) should have a `Dockerfile`. `docker-compose.yml` starts all the containers needed for the API to run locally. 

### Swagger
We use Swagger via swaggo to document the API. Docs can be accessed at `localhost:8080/api/swagger/index.html` when running the API locally. 

The `swag` command can be installed with the following command:

`go install github.com/swaggo/swag/cmd/swag@latest`

Run the following command after making any changes to the API, including any new or modified endpoints: 

`swag init -g .\cmd\server\main.go`

Note: all commands like the above should be run from the api directory

When you run `swag init`, it might create a left and right delim property in docs.go that cause an error. These are safe to be deleted, I'm not sure why they show up. One for the tech debt backlog.

#### Annotation
For endpoints to appear in the Swagger UI, they have to be annotated and the docs need to be regenerated. See books/handler for an example of annotation. 

While all endpoints should return the same shape of data (defined in delivery/response), Swagger doesn't handle go generics, so each package should have a swagger.go file to define the documented concrete type that gets returned by an endpoint. These don't get used in practice, they are purely for documentation. 

#### Auth endpoints 
Some endpoints rely on the auth middleware to check the logged in user. To mark this in swagger, add

`// @Security BearerAuth`

to the endpoint godoc (in the handler file). This will require the user to have set their authorization in swagger to call the endpoint. 

When logging in via the swagger UI, the format for the token you put in is `Bearer <token>`, rather than simply `<token>`. It's a bit annoying, but it does the job. You can generate a bearer token by calling the log in endpoint with valid credentials (which can be created by calling the register endpoint).

### Migrate 
Migrations are applied using the migrate tool in `docker-compose.yml`. You can call the same command on any database you like to apply the migrations from the project.

Migrations are stored in `/migrations`. They are numbered sequentially and each migration has an up and a down. 

### Automated testing 
Automated API tests are run with the `go test` command. There's a test runner in the docker-compose file that runs the same command on a test database. As mentioned elsewhere, unit tests should run against a real database, so we don't need to mock dependencies (at least for now, there will likely be stuff we have to mock in the future).

Unit and simple integration tests live next to the code they're testing. This makes it really easy to see what's getting tested where, and means you don't have to worry about mirroring the project structure in a test package elsewhere. 

Complex, multi-stage tests live in `internal_test`. For example, the first test writtne here tested the end-to-end auth flow of registering a user, logging in as that user, and calling a protected endpoint with their JWT. 

#### Future improvement 
We may choose to use go testcontainers in the future to run our automated tests in isolated containers. For now, a db in docker should be sufficient for our needs. 

## Architecture
The API is divided into a number of packages which are largely flat. Each core entity has its own package, e.g. Books, Ratings, Reviews, Authors, Users. Each of these packages contains the entities themselves, the persistence details (repos), orchestration (services), and delivery (route handlers). 

There are also a number of supplementary packages for config (which is read from .env) and delivery (setting up routes). 

Interfaces should only be used when there is likely to be variation. This probably won't happen often, so it's alright for e.g. services to depend on concrete repos. 

With the lack of interfaces, there won't be much mocking. Unit tests and integration tests can be run with a testing database via a test runner in docker (once we've set that up).

### Folder structure 
Below is the canonical folder structure for an aggregate root (e.g. books, users, ratings), but also for 

#### Handlers 
These are the entry points to the server. Their responsibility is ensuring that requests are parsed correctly and responses are formatted correctly. They delegate all business logic conerns to services. A handler should ideally only depend on one service. 

A package should define its own errors. The handler is responsible for setting the http code for these errors (as the service shouldn't know about http). Auth has an example of using a switch statement to determine the codes for concrete errors that the service can throw. 

#### Services 
Services orchestrate business logic. Services can depend on other services, e.g. an auth service may have a function to register a user. Because it's a security concern, /auth should be the entry point, but actually creating a user entity and persisting it is not an auth concern, that should be handled by the user service. So, the auth service can depend on the user service, the registration endpoint can call the user service to create the actual user (including ensuring uniqueness of emails and so on), but auth should handle security concerns such as password hashing and generating tokens. 

#### Repos
Repos define how data is persisted. This will usually be in our Postgres DB, but that detail shouldn't matter to the rest of the application. Repos should be private, services outside of a package shouldn't be calling repos directly (we can make sure of this by setting the first letter of a repo struct to lowercase, e.g. ratingRepo). 

##### A note about db
We have an interface called DBTX which represents the shape of a query both in and out of the context of a transaction. What this means is that a service can coordinate transactions while repos always just execute queries. Repos shouldn't know about transactions. 

In terms of di, handlers should not know about db as that's a persistence detail and handlers are a transport detail. Services need to know about db in order to start, commit and roll back transactions, but they never call queries themselves. Repos (and readers) are the only things that execute queries. Repos must be passed a db from the service in case they're called from within a transaction, whereas readers own their own db connection as they never run transactions. 

#### Entities
Entities are the things the application reasons about. Business logic is carried out on entities. Each entity has a UUID as its identifier. Note: the way that data is persisted may not match one-to-one with the actual entity definitions. 

#### Queries 
Queries are organisationally separated from domain services and repos. The reason for this is that the system is likely going to be doing a great deal many more reads than writes. Reads should be fast and can take data from across bounded contexts, whereas writes require heavy validation. If we get the writing right, then we don't need to validate data in reads. This is basically CQRS but without all the faff that generally goes into it. 

#### DTOs
DTOs (data transfer objects) are specific views of data that are useful for consumers of the API but don't match the entity definitions exactly. For example, you might have a DTO for a review. This may contain the scores for the rating associated with the review and the username of the user who posted it, both of which aren't in the Review entity definition but are aggregated for the use of consumers (mainly the frontend). 

#### Side effects
An endpoint will either retrieve data or mutate state, but the use case might have side effects. For example, when a rating is created, the service will ensure that the rating entity is created and persisted, but we'll also want the average score of the book to be recalculated. That mutation is a side effect and shouldn't be handled in the main flow. Rather, an event should be raised that gets processed in the background. 

An example of how we might do this is a transactional outbox. We have an events table in the db that gets polled by a worker. Events have a processed flag, so if the worker finds any unprocessed events, it will process that event and set the processed flag in a transaction. That ensures eventual consistency and means we don't lose data if for some reason the worker can't commit in progress changes. 

There may be other cases where we need side effects handled in a timely manner. In those cases, we may use a message queue so instead of relying on polling, the event gets picked up from the queue. 