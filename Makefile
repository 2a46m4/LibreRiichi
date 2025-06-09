all: frontend go

frontend:
	yarn run build

go:
	go build cmd/main.go