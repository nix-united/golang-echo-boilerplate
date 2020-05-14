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

    
### Usage

1. Run your app, and browse to {HOST}:{PORT}/swagger/index.html. You will see Swagger 2.0 API documents.
2. Using the API documentation, make requests to register a user (if necessary) and login.
3. After the successful login, copy a token from the response, then click "Authorize" and in a popup that opened, enter the value for "apiKey" in a form:
"Bearer {token}". For example:


    Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODk0NDA5NjYsIm9yaWdfaWF0IjoxNTg5NDM5OTY2LCJ1c2VyX2lkIjo1fQ.f8dSG3NxFLHwyA5-XIYALT5GtXm4eiH-motqtqAUBOI 

   
Then, click "Authorize" and close the popup.
Now, you are able to make requests which require authentication.
