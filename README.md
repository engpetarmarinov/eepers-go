# Eepers Go
<p align=center>
  <img src="./assets/icon.png">
</p>

Enter a mysterious world inhabited by strange creatures called **Eepers**. 

Originally inspired by [eepers](https://github.com/tsoding/eepers), rewritten in Go with raylib.

## Features

- **Keyboard & Gamepad Support** - Play with arrow keys/WASD or any standard gamepad (Xbox, PlayStation, etc.)
- **Mystical Creatures** - Encounter Guardian Eepers, Mother Eepers, Gnome Eepers, and the Father
- **Stealth & Strategy** - Outsmart patrol patterns and use bombs to clear your path

## Build and Run

### Development
```console
go run ./cmd/eepers-go/main.go
```

### Build for Distribution

Build for all platforms:
```console
make package-all
```

Or build for specific platforms:

**macOS (ARM64)**
```console
make package-darwin
```
Creates `builds/eepers-macos-arm64.zip` containing Eepers.app bundle

**Windows (AMD64)**
```console
make package-windows
```
Creates `builds/eepers-windows-amd64.zip` with executable and dependencies

**Linux (AMD64)**
```console
make package-linux
```
Creates `builds/eepers-linux-amd64.tar.gz` with executable and dependencies

Distribution packages will be created in the `builds/` directory.
