NAME=homevision-exercise

.PHONY: build
## build: Compile the application.
build:
	@go build -o $(NAME)

.PHONY: run
## run: Build and Run.
run: build
	@./$(NAME)

.PHONY: clean
## clean: Clean project and previous builds.
clean:
	@rm -f $(NAME)
	@rm -f *.jpg

.PHONY: deps
## deps: Download modules
deps:
	@go mod download

.PHONY: test
## test: Run tests
test:
	@go test ./...

.PHONY: lint
## lint: Run golangci-lint (must be installed separately - https://golangci-lint.run/usage/install/)
lint:
	@golangci-lint run


.PHONY: help
all: help
# help: show this help message
help: Makefile
	@echo
	@echo " Choose a command to run in "$(APP_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
