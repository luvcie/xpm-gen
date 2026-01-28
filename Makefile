# Name of the binary
NAME = xpm-gen

# Versioning
# Base version
BASE_VERSION = v1.0
# Count commits for patch level (default to 0 if not a git repo)
COMMIT_COUNT = $(shell git rev-list --count HEAD 2>/dev/null || echo 0)
# Full version string
VERSION = $(BASE_VERSION).$(COMMIT_COUNT)

# Go build flags
# -ldflags "-X main.Version=..." injects the variable into the binary
LDFLAGS = -ldflags "-X main.Version=$(VERSION)"

# Targets
.PHONY: all build clean run re

all: build

build:
	@echo "Building $(NAME) version $(VERSION)..."
	@go build $(LDFLAGS) -o $(NAME) main.go
	@echo "Done! Run ./$(NAME) -version to check."

run:
	@go run $(LDFLAGS) main.go

clean:
	@rm -f $(NAME)
	@rm -f *.xpm *.png

re: clean all
