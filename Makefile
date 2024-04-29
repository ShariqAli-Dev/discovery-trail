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
docker-build-install:
	go install github.com/a-h/templ/cmd/templ@latest
	go get ./...
	go mod tidy
	go mod download 
	templ generate
	cd ./ui/static && npm install
	cd ./ui/static && npx tailwindcss -i ./src/input.css -o ./public/global.css
	cd ./ui/static && npm run build 
	go run ./cmd/template_formatter
	go build -o ./bin/web ./cmd/web