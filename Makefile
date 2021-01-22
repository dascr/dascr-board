include .env
include frontend/.env
export

.PHONY:

build-linux_64: generate download
	@echo "[*] Building for linux x64"
	@CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/linux_amd64/dascr-board
	@echo "[OK] App binary was created!"
	@echo "[OK] Your backend binary is at ./dist/<os>/"

build-linux_386: generate download
	@echo "[*] Building for linux i386"
	@CGO_ENABLED=1 GOOS=linux GOARCH=386 go build -ldflags="-s -w" -o dist/linux_386/dascr-board
	@echo "[OK] App binary was created!"
	@echo "[OK] Your backend binary is at ./dist/<os>/"

build-mac: generate download
	@echo "[*] Building for mac"
	@CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/darwin_amd64/dascr-board
	@echo "[OK] App binary was created!"
	@echo "[OK] Your backend binary is at ./dist/<os>/"

build-armv5: generate download
	@echo "[*] Building for armv5"
	@CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=5 go build -ldflags="-s -w" -o dist/arm_5/dascr-board
	@echo "[OK] App binary was created!"
	@echo "[OK] Your backend binary is at ./dist/<os>/"

build-armv6: generate download
	@echo "[*] Building for armv6"
	@CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=6 go build -ldflags="-s -w" -o dist/arm_6/dascr-board
	@echo "[OK] App binary was created!"
	@echo "[OK] Your backend binary is at ./dist/<os>/"

build-armv7: generate download
	@echo "[*] Building for armv7"
	@CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags="-s -w" -o dist/arm_7/dascr-board
	@echo "[OK] App binary was created!"
	@echo "[OK] Your backend binary is at ./dist/<os>/"

build-armv8_64: generate download
	@echo "[*] Building for armv8_64"
	@CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o dist/arm64_8/dascr-board
	@echo "[OK] App binary was created!"
	@echo "[OK] Your backend binary is at ./dist/<os>/"

build-frontend: clean-frontend
	@echo "[*] Building Svelte App"
	@cd frontend && yarn install
	@NODE_ENV=production yarn --cwd "./frontend" run build
	@rm ./frontend/public/build/*.map
	@echo "[OK] Svelte App was built"
	@echo "[OK] Serve content of ./frontend/public via a webserver"

generate:
	@echo "[*] Embedding via parcello"
	@PARCELLO_RESOURCE_DIR=./static go generate ./...
	@echo "[OK] Done bundeling things"

download:
	@echo "[*] go mod dowload"
	@go mod download

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
