live/templ:
	go tool templ generate --watch --proxy="http://localhost:8080" --cmd="go run cmd/gluttony/main.go"
templ/watch:
	go tool templ generate --watch
release:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o gluttony cmd/main.go
