.PHONY: build-win-amd64 build-linux-amd64 build-darwin-arm64 clear

BUILDS_DIR := ./builds
RAYLIB_DIR := ./raylib
RAYLIB_VERSION := 5.5

build-win-amd64:
	mkdir -p $(BUILDS_DIR)/windows-amd64/assets && cp -r ./assets $(BUILDS_DIR)/windows-amd64/
	cp $(RAYLIB_DIR)/raylib-$(RAYLIB_VERSION)_win64_mingw-w64/raylib.dll $(BUILDS_DIR)/windows-amd64/raylib.dll
	GOOS=windows GOARCH=amd64 go build -ldflags="-H=windowsgui" -o $(BUILDS_DIR)/windows-amd64/eepers.exe cmd/eepers-go/main.go

build-linux-amd64:
	mkdir -p $(BUILDS_DIR)/linux-amd64/assets && cp -r ./assets $(BUILDS_DIR)/linux-amd64/
	cp $(RAYLIB_DIR)/raylib-$(RAYLIB_VERSION)_linux_amd64/libraylib.so $(BUILDS_DIR)/linux-amd64/libraylib.so
	GOOS=linux GOARCH=amd64 go build -o $(BUILDS_DIR)/linux-amd64/eepers cmd/eepers-go/main.go

build-darwin-arm64:
	mkdir -p $(BUILDS_DIR)/darwin-arm64/assets && cp -r ./assets $(BUILDS_DIR)/darwin-arm64/
	cp $(RAYLIB_DIR)/raylib-$(RAYLIB_VERSION)_macos/libraylib.dylib $(BUILDS_DIR)/darwin-arm64/libraylib.dylib
	GOOS=darwin GOARCH=arm64 go build -o $(BUILDS_DIR)/darwin-arm64/eepers cmd/eepers-go/main.go

clear:
	rm -rf $(BUILDS_DIR)
