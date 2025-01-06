binary_name = neba
main_package_path = ./cmd/neba
build_dir = build

# default target
all: clean tidy prod

# clean the build directory
clean:
	@if [ -d $(build_dir) ]; then \
		rm -rf $(build_dir); \
	fi

# tidy modfiles and format .go files
.PHONY: tidy
tidy:
	go mod tidy -v
	go fmt ./...

# build the application
.PHONY: build
build:
	go build -o /tmp/bin/${binary_name} ${main_package_path}

# run the application
.PHONY: run
run: build
	/tmp/bin/${binary_name}

# build the application for production
.PHONY: prod
prod:
	GOOS=windows GOARCH=amd64 go build -ldflags='-s' -o build/${binary_name} ${main_package_path}
