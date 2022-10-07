include .env
include frontend/.env
export

.PHONY:

build-linux_64: download
	@echo "[*] Building for linux x64"
	@CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/linux_amd64/dascr-board
	@echo "[OK] App binary was created!"
	@echo "[OK] Your backend binary is at ./dist/<os>/"

build-linux_386: download
	@echo "[*] Building for linux i386"
	@CGO_ENABLED=1 GOOS=linux GOARCH=386 go build -ldflags="-s -w" -o dist/linux_386/dascr-board
	@echo "[OK] App binary was created!"
	@echo "[OK] Your backend binary is at ./dist/<os>/"

build-mac: download
	@echo "[*] Building for mac"
	@CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/darwin_amd64/dascr-board
	@echo "[OK] App binary was created!"
	@echo "[OK] Your backend binary is at ./dist/<os>/"

build-armv5: download
	@echo "[*] Building for armv5"
	@CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=5 go build -ldflags="-s -w" -o dist/arm_5/dascr-board
	@echo "[OK] App binary was created!"
	@echo "[OK] Your backend binary is at ./dist/<os>/"

build-armv6: download
	@echo "[*] Building for armv6"
	@CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=6 go build -ldflags="-s -w" -o dist/arm_6/dascr-board
	@echo "[OK] App binary was created!"
	@echo "[OK] Your backend binary is at ./dist/<os>/"

build-armv7: download
	@echo "[*] Building for armv7"
	@CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags="-s -w" -o dist/arm_7/dascr-board
	@echo "[OK] App binary was created!"
	@echo "[OK] Your backend binary is at ./dist/<os>/"

build-armv8_64: download
	@echo "[*] Building for armv8_64"
	@CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o dist/arm64_8/dascr-board
	@echo "[OK] App binary was created!"
	@echo "[OK] Your backend binary is at ./dist/<os>/"

build-darwin-apple-silicon: download
	@echo "[*] Building for darwin_arm64 (apple silicon)"
	@CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/darwin_arm64/dascr-board
	@echo "[OK] App binary was created!"
	@echo "[OK] Your backend binary is at ./dist/<os>/"

build-frontend: clean-frontend
	@echo "[*] Building SvelteKit App"
	@cd frontend &&  pnpm install --silent
	@NODE_ENV=production pnpm --dir "./frontend" run --silent build
	@echo "[OK] SvelteKit App was built"
	@echo "[OK] Serve content of ./frontend/build via a webserver"

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
	@echo "[*] Cleanup SvelteKit App"
	@rm -rf ./frontend/build
	@echo "[OK] Cleanup done"

clean-all: clean-backend clean-frontend

run-dev-backend:
	@echo "[*] Starting Backend Development"
	@go run main.go

run-dev-frontend:
	@echo "[*] Starting Frontend Development, listening to 127.0.0.1"
	@pnpm --dir "./frontend" run dev