NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

XC_GOOS = $(shell go env GOOS)
XC_GOARCH = $(shell go env GOARCH)

all: test build

build: deps
	@echo "$(OK_COLOR)==> Building $(XC_GOOS)/$(XC_GOARCH)$(NO_COLOR)"
	@gox -parallel=1 -os "$(XC_GOOS)" -arch "$(XC_GOARCH)" -output "bin/{{.Dir}}" ./...

test: deps
	@echo "$(OK_COLOR)==> Running Tests...$(NO_COLOR)"
	@go test -cover ./...

deps:
	@echo "$(OK_COLOR)==> Installing dependencies$(NO_COLOR)"
	@go get -d -v -t ./...

release: clean-pkg deps
	@echo "$(OK_COLOR)==> Releasing$(NO_COLOR)"
	@gox -os "darwin linux windows" -arch "amd64" -output "pkg/{{.OS}}_{{.Arch}}/{{.Dir}}" ./...

clean: clean-bin clean-pkg

clean-bin:
	@rm -rf inductor inductor.exe bin/

clean-pkg:
	@rm -rf  pkg/

.PHONY: all build test deps clean release
