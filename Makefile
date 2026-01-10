.PHONY: build-win-amd64 build-linux-amd64 build-darwin-arm64 clear package-darwin package-windows package-linux package-all

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
	mkdir -p $(BUILDS_DIR)/darwin-arm64/Eepers.app/Contents/MacOS
	mkdir -p $(BUILDS_DIR)/darwin-arm64/Eepers.app/Contents/Resources
	cp -r ./assets $(BUILDS_DIR)/darwin-arm64/Eepers.app/Contents/Resources/
	cp $(RAYLIB_DIR)/raylib-$(RAYLIB_VERSION)_macos/libraylib.dylib $(BUILDS_DIR)/darwin-arm64/Eepers.app/Contents/MacOS/libraylib.dylib
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILDS_DIR)/darwin-arm64/Eepers.app/Contents/MacOS/eepers cmd/eepers-go/main.go
	@echo "Creating app icon..."
	@if command -v sips >/dev/null 2>&1 && command -v iconutil >/dev/null 2>&1; then \
		mkdir -p $(BUILDS_DIR)/darwin-arm64/icon.iconset; \
		sips -z 16 16     ./assets/icon.png --out $(BUILDS_DIR)/darwin-arm64/icon.iconset/icon_16x16.png >/dev/null 2>&1; \
		sips -z 32 32     ./assets/icon.png --out $(BUILDS_DIR)/darwin-arm64/icon.iconset/icon_16x16@2x.png >/dev/null 2>&1; \
		sips -z 32 32     ./assets/icon.png --out $(BUILDS_DIR)/darwin-arm64/icon.iconset/icon_32x32.png >/dev/null 2>&1; \
		sips -z 64 64     ./assets/icon.png --out $(BUILDS_DIR)/darwin-arm64/icon.iconset/icon_32x32@2x.png >/dev/null 2>&1; \
		sips -z 128 128   ./assets/icon.png --out $(BUILDS_DIR)/darwin-arm64/icon.iconset/icon_128x128.png >/dev/null 2>&1; \
		sips -z 256 256   ./assets/icon.png --out $(BUILDS_DIR)/darwin-arm64/icon.iconset/icon_128x128@2x.png >/dev/null 2>&1; \
		sips -z 256 256   ./assets/icon.png --out $(BUILDS_DIR)/darwin-arm64/icon.iconset/icon_256x256.png >/dev/null 2>&1; \
		sips -z 512 512   ./assets/icon.png --out $(BUILDS_DIR)/darwin-arm64/icon.iconset/icon_256x256@2x.png >/dev/null 2>&1; \
		sips -z 512 512   ./assets/icon.png --out $(BUILDS_DIR)/darwin-arm64/icon.iconset/icon_512x512.png >/dev/null 2>&1; \
		sips -z 1024 1024 ./assets/icon.png --out $(BUILDS_DIR)/darwin-arm64/icon.iconset/icon_512x512@2x.png >/dev/null 2>&1; \
		iconutil -c icns $(BUILDS_DIR)/darwin-arm64/icon.iconset -o $(BUILDS_DIR)/darwin-arm64/Eepers.app/Contents/Resources/AppIcon.icns; \
		rm -rf $(BUILDS_DIR)/darwin-arm64/icon.iconset; \
		echo "Icon created successfully"; \
	else \
		echo "Warning: sips or iconutil not found, skipping icon creation"; \
	fi
	echo '<?xml version="1.0" encoding="UTF-8"?>\n\
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">\n\
<plist version="1.0">\n\
<dict>\n\
	<key>CFBundleExecutable</key>\n\
	<string>eepers</string>\n\
	<key>CFBundleIconFile</key>\n\
	<string>AppIcon</string>\n\
	<key>CFBundleIdentifier</key>\n\
	<string>com.engpetarmarinov.eepers</string>\n\
	<key>CFBundleName</key>\n\
	<string>Eepers</string>\n\
	<key>CFBundleVersion</key>\n\
	<string>1.0</string>\n\
	<key>CFBundleShortVersionString</key>\n\
	<string>1.0</string>\n\
	<key>CFBundlePackageType</key>\n\
	<string>APPL</string>\n\
	<key>LSMinimumSystemVersion</key>\n\
	<string>10.13</string>\n\
	<key>NSHighResolutionCapable</key>\n\
	<true/>\n\
</dict>\n\
</plist>' > $(BUILDS_DIR)/darwin-arm64/Eepers.app/Contents/Info.plist

clear:
	rm -rf $(BUILDS_DIR)

package-darwin: build-darwin-arm64
	@echo "Creating macOS distribution package..."
	cd $(BUILDS_DIR)/darwin-arm64 && zip -r -q ../eepers-macos-arm64.zip Eepers.app
	@echo "Created: $(BUILDS_DIR)/eepers-macos-arm64.zip"

package-windows: build-win-amd64
	@echo "Creating Windows distribution package..."
	cd $(BUILDS_DIR)/windows-amd64 && zip -r -q ../eepers-windows-amd64.zip .
	@echo "Created: $(BUILDS_DIR)/eepers-windows-amd64.zip"

package-linux: build-linux-amd64
	@echo "Creating Linux distribution package..."
	cd $(BUILDS_DIR)/linux-amd64 && tar -czf ../eepers-linux-amd64.tar.gz .
	@echo "Created: $(BUILDS_DIR)/eepers-linux-amd64.tar.gz"

package-all: package-darwin package-windows package-linux
	@echo "All distribution packages created successfully!"
	@echo "Files ready for GitHub release:"
	@ls -lh $(BUILDS_DIR)/*.zip $(BUILDS_DIR)/*.tar.gz 2>/dev/null || true

