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

## Architecture