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

### Migrate 
Migrations are applied using the migrate tool in `docker-compose.yml`. You can call the same command on any database you like to apply the migrations from the project.

Migrations are stored in `/migrations`. They are numbered sequentially and each migration has an up and a down. 

### Automated testing 
Automated API tests are run with the `go test` command. There's a test runner in the docker-compose file that runs the same command on a test database. As mentioned elsewhere, unit tests should run against a real database, so we don't need to mock dependencies (at least for now, there will likely be stuff we have to mock in the future).

## Architecture
The API is divided into a number of packages which are largely flat. Each core entity has its own package, e.g. Books, Ratings, Reviews, Authors, Users. Each of these packages contains the entities themselves, the persistence details (repos), orchestration (services), and delivery (route handlers). 

There are also a number of supplementary packages for config (which is read from .env) and delivery (setting up routes). 

Interfaces should only be used when there is likely to be variation. This probably won't happen often, so it's alright for e.g. services to depend on concrete repos. 

With the lack of interfaces, there won't be much mocking. Unit tests and integration tests can be run with a testing database via a test runner in docker (once we've set that up).