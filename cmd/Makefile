GO=go 
GOFLAGS = -trimpath -tags release -buildmode exe
GOAPP = app.o

build:
	$(GO) build $(GOFLAGS) -o $(GOAPP)

run: build
	echo "clear; ./$(GOAPP)" | bash