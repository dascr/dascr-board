include .env
include frontend/.env
export

.PHONY: build

generate:
	@echo "[*] Embedding via parcello"
	@PARCELLO_RESOURCE_DIR=./static go generate ./...
	@echo "[OK] Done bundeling things"

build-backend: clean-backend generate
	@echo "[*] go mod dowload"
	@go mod download
	@echo "[*] Building for linux"
	@GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/linux_amd64/dascr-board
	@GOOS=linux GOARCH=386 go build -ldflags="-s -w" -o dist/linux_386/dascr-board
	@echo "[*] Building for mac"
	@GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/darwin_amd64/dascr-board
	@echo "[OK] App binary was created!"
	@echo "[OK] Your backend binary is at ./dist/<os>/"

build-frontend: clean-frontend
	@echo "[*] Building Svelte App"
	@NODE_ENV=production yarn --cwd "./frontend" run build
	@rm ./frontend/public/build/*.map
	@echo "[OK] Svelte App was built"
	@echo "[OK] Serve content of ./frontend/public via a webserver"

build-all: build-backend build-frontend

clean-backend:
	@echo "[*] Cleanup Backend App"
	@rm -f *.db
	@rm -rf uploads
	@rm -rf ./dist
	@echo "[OK] Cleanup done"

clean-frontend:
	@echo "[*] Cleanup Svelte App"
	@rm -rf ./frontend/public/build
	@echo "[OK] Cleanup done"

clean-all: clean-backend clean-frontend

run-dev-backend:
	@echo "[*] Starting Backend Development"
	@go run main.go

run-dev-frontend:
	@echo "[*] Starting Frontend Development, listening to 127.0.0.1"
	@yarn --cwd "./frontend" run dev

run-dev-frontend-exp:
	@echo "[*] Starting Frontend Development, listening to 0.0.0.0"
	@yarn --cwd "./frontend" run dev-expose
