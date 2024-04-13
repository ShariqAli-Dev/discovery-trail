turso:
	turso db shell discovery-trail
css:
	cd ./ui/static && npx tailwindcss -i ./src/input.css -o ./public/global.css --watch
static:
	cd ./ui/static && npx vite build 
build: static
	go build -o ./bin/web ./cmd/web
run: static
	go run ./cmd/template_formatter && go run ./cmd/web
