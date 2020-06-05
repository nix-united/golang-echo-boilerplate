# Overview
This is an example of demo application based on the echo framework.
The application has a basic functional, it does register and authenticate with jwt users, create, 
update and delete posts, create and run migrations, validate requests, also have tests and swagger docs.

## Usage
1. Copy .env.dist to .env and set the environment variables.
2. Run your application using the command in the terminal:

    `docker-compose up -d`
3. Build the Swagger documentation using these commands in the terminal:
    
    `go get -v github.com/swaggo/swag/cmd/swag`
    
    `$GOPATH/bin/swag init`
4. Browse to {HOST}:{PORT}/swagger/index.html. You will see Swagger 2.0 API documents.
5. Using the API documentation, make requests to register a user (if necessary) and login.
6. After the successful login, copy a token from the response, then click "Authorize" and in a popup that opened, enter the value for "apiKey" in a form:
"Bearer {token}". For example:


    Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODk0NDA5NjYsIm9yaWdfaWF0IjoxNTg5NDM5OTY2LCJ1c2VyX2lkIjo1fQ.f8dSG3NxFLHwyA5-XIYALT5GtXm4eiH-motqtqAUBOI 

   
Then, click "Authorize" and close the popup.

## Directories
1. **/deploy** contains the container (Docker) package configuration and template(docker-compose) for project deployment.

2. **/development** includes Docker and docker-compose files for setup linter.

3. **/migrations** has files for run migrations.

4. **/server** is the main project folder. This folder contains the executable server.go.

5. **/server/builders** contains builders for initializing entities.

6. **/server/db** has seeders and method for connecting to the database.

7. **/server/handlers** contains request handlers.

8. **/server/models** includes structures describing data models.

9. **/server/repositories** contains methods for selecting entities from the database.

10. **/server/requests** has structures describing the parameters of incoming requests.

11. **/server/responses** includes structures describing the parameters of outgoing response.

12. **/server/routes** has a file for configuring routes.

13. **/server/services** contains methods for creating entities.

14. **/server/validation** contains the request validator.

15. **/tests**  includes tests and test data.

## Code quality
For control code quality we are use [golangci-lint](https://github.com/golangci/golangci-lint).
Golangci-lint is a linters aggregator.

Why we use linters? Linters help us:
1. Finding critical bugs
2. Finding bugs before they go live
3. Finding performance errors
4. To speed up the code review, because reviewers do not spend time searching for syntax errors and searching for
violations of generally accepted code style
5. The quality of the code is guaranteed at a fairly high level.

### How to use
Linter tool wrapped to docker-compose and first of all need to build container with linters

- `make lint-build`

Next you need to run linter to check bugs ant errors

- `make lint-check` - it will log to console what bugs and errors linters found

Finally, you need to fix all problems manually or using autofixing (if it's supported by the linter)

- `make lint-fix` 


## Swagger documentation


### Installation

1 Get the binary for github.com/swaggo/swag/cmd/swag:


    go get github.com/swaggo/swag/cmd/swag


2 In .env file set the value for the "HOST" variable. This is a host to which Swagger will make API requests. For example, for local development:


    HOST=localhost 
  
    
3 Run "swag init" in the project's root folder which contains the main.go file. This will parse your comments and generate the required files (docs folder and docs/docs.go).


     $GOPATH/bin/swag init 

    

Now, you are able to make requests which require authentication.

## Libraries
Migrations - https://github.com/ShkrutDenis/go-migrations

Jwt - https://github.com/dgrijalva/jwt-go

Swagger - https://github.com/swaggo/echo-swagger

Mocking db - https://github.com/selvatico/go-mocket

Orm - https://github.com/jinzhu/gorm