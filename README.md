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
