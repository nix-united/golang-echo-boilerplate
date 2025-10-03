<h1 align="center">Welcome to <span style="color:mediumseagreen">Echo boilerplate</span> 👋</h1>

<p align="center">
  <img src="https://img.shields.io/github/go-mod/go-version/nix-united/golang-echo-boilerplate" height="25"/>
  <img src="https://img.shields.io/badge/gorm-v1.25-green" height="25"/>
  <img src="https://img.shields.io/badge/swagger-v1.16-orange" height="25"/>
  <img src="https://img.shields.io/badge/echo-v4.13-yellow" height="25"/>
  <img src="https://img.shields.io/badge/docker-support-darkgreeen" height="25"/>
  <img src="https://goreportcard.com/badge/github.com/nix-united/golang-echo-boilerplate?style=flat" height="25" />
  <img src="https://img.shields.io/github/actions/workflow/status/nix-united/golang-echo-boilerplate/ci.yml" height="25"/>
  <img src="https://img.shields.io/github/checks-status/nix-united/golang-echo-boilerplate/master?label=CI checks" height="25"/>
  <img src="https://img.shields.io/github/license/nix-united/golang-echo-boilerplate?style=flat&color=blue" height="25" />
</p>

It's an API Skeleton project based on Echo framework.
Our aim is reducing development time on default features that you can meet very often when your work on API.
There is a useful set of tools that described below. Feel free to contribute!

## What's inside:

- Registration
- Authentication with JWT
- CRUD API for posts
- Migrations
- Request validation
- Swagger docs
- Environment configuration
- Docker development environment

## Usage
1. Copy .env.dist to .env and set the environment variables. There are examples for all the environment variables except COMPOSE_USER_ID, COMPOSE_GROUP_ID which are used by the linter. To get the current user ID, run in terminal:
    
    `echo $UID`
    
    In the .env file set these variables:

    `COMPOSE_USER_ID="username in current system"` - your username in system

    `COMPOSE_GROUP_ID="user uid"` - the user ID which you got in the terminal

2. Run your application using the command in the terminal:

    `docker-compose up`
3. Browse to {HOST}:{PORT}/swagger/index.html. You will see Swagger 2.0 API documents.
4. Using the API documentation, make requests to register a user (if necessary) and login.
5. After the successful login, copy a token from the response, then click "Authorize" and in a popup that opened, enter the value for "apiKey" in a form:
"Bearer {token}". For example:


    Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODk0NDA5NjYsIm9yaWdfaWF0IjoxNTg5NDM5OTY2LCJ1c2VyX2lkIjo1fQ.f8dSG3NxFLHwyA5-XIYALT5GtXm4eiH-motqtqAUBOI 

   
Then, click "Authorize" and close the popup.
Now, you are able to make requests which require authentication.

## Directories
1. **/cmd** entry points.

2. **/config** has structures which contains service config.

3. **/db** has seeders and method for connecting to the database.

4. **/deploy** contains the container (Docker) package configuration and template(docker-compose) for project deployment.

5. **/development** includes Docker and docker-compose files for setup linter.

6. **/migrations** has files for run migrations.

7. **/models** includes structures describing data models.

8. **/repositories** contains methods for selecting entities from the database.

9. **/requests** has structures describing the parameters of incoming requests, and validator.

10. **/responses** includes structures describing the parameters of outgoing response.

11. **/server** is the main project folder. This folder contains the executable server.go.

12. **/server/builders** contains builders for initializing entities.

13. **/server/handlers** contains request handlers.

14. **/server/routes** has a file for configuring routes.

15. **/services** contains methods for creating entities.

16. **/tests**  includes tests and test data.

## Code quality
For control code quality we use [golangci-lint](https://github.com/golangci/golangci-lint).
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


## Libraries
Migrations - https://github.com/ShkrutDenis/go-migrations

Jwt - https://github.com/golang-jwt/jwt/v5

Swagger - https://github.com/swaggo/echo-swagger

Mocking db - https://github.com/selvatico/go-mocket

Orm - https://gorm.io/gorm

## License

The project is developed by [NIX][1] and distributed under [MIT LICENSE][2]

[1]: https://nixs.com/
[2]: https://github.com/nixsolutions/golang-echo-boilerplate/blob/master/LICENSE
