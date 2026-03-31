# Build Directory

The build directory is used to house all the build files and assets for your application.

The structure is:

* bin - Output directory
* darwin - macOS specific files
* appicon.png - Application icon

## macOS

The `darwin` directory holds files specific to macOS builds.
These may be customised and used as part of the build. To return these files to the default state, simply delete them and build with `wails build`.

The directory contains the following files:

- `Info.plist` - the main plist file used for macOS builds. It is used when building using `wails build`.
- `Info.dev.plist` - same as the main plist file but used when building using `wails dev`.