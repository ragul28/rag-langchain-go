build:
	go build -ldflags="-s -w"

install:
	go install
	
run:
	go build && ./rag-langchain-go

mod:
	go mod tidy
	go mod verify
	go mod vendor
