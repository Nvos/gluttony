live/templ:
	go tool templ generate --watch --proxy="http://localhost:8080" --cmd="go run cmd/gluttony/main.go"
templ/watch:
	go tool templ generate --watch