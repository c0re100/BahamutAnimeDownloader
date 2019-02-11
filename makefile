GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
BINARY_NAME=AniDownloader

all: build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)
deps:
	$(GOGET) github.com/grafov/m3u8
	$(GOGET) github.com/korovkin/limiter
	$(GOGET) github.com/MercuryEngineering/CookieMonster
	$(GOGET) gopkg.in/cheggaaa/pb.v1
