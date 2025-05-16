lint_docker_compose_file = "./development/golangci_lint/docker-compose.yml"

lint-build:
	@echo "🌀 ️container is building..."
	@docker-compose --file=$(lint_docker_compose_file) build -q
	@echo "✔  ️container built"

lint-check:
	@echo "🌀️ code linting..."
	@docker-compose --file=$(lint_docker_compose_file) run --rm echo-golinter golangci-lint version && golangci-lint run \
 		&& echo "✔️  checked without errors" \
 		|| echo "☢️  code style issues found"

lint-fix:
	@echo "🌀 ️code fixing..."
	@docker-compose --file=$(lint_docker_compose_file) run --rm echo-golinter golangci-lint run --fix \
		&& echo "✔️  fixed without errors" \
		|| (echo "⚠️️  you need to fix above issues manually" && exit 1)
	@echo "⚠️️ run \"make lint-check\" again to check what did not fix yet"

organize-imports:
	@gci write --custom-order -s standard -s "prefix(github.com/nix-united/golang-echo-boilerplate)" -s default --skip-generated --skip-vendor .

lint:
	go tool golangci-lint run ./...

goose-create: # Creates goose migration. Example: NAME=migration_name make goose-create
	goose -dir migrations create $(NAME) sql
