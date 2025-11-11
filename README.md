English / [简体中文](./README-CN.md)

# Apple Music Downloader

A Go-based Apple Music downloader that can fetch songs, albums, playlists, and MVs.

## Prerequisites
- Install and add to `PATH`:
  - `MP4Box` (for tagging and MV muxing): https://gpac.io/downloads/gpac-nightly-builds/
  - `mp4decrypt` (required for MV decryption): https://www.bento4.com/downloads/
  - `ffmpeg` (optional, for animated artwork and post-download conversion)
- Ensure the decryption helper [wrapper-z](https://github.com/zesty-zesty/wrapper-z) is running before downloads.
- Go `1.23+` is required to build and run.

## Features
- Embed cover art, export/save LRC lyrics (requires `media-user-token`).
- Support segmented/translated lyrics (Beta), including pronunciation handling.
- Download via URL: albums, playlists, songs, artists, and stations.
- Wizard-style interaction: menu-driven `search`/`rip` and settings.
- Decrypt during download to reduce memory usage for large files.
- MV download and mux with `MP4Box` (audio + video).

## Command-line (non-interactive)
1. Build the project:
   - `go build .`
2. Start the decryption helper: [wrapper](https://github.com/WorldObservationLog/wrapper).
3. Run the program:
   - `./main` (Windows: `./main.exe`)
4. Download a single song/album/playlist:
   - Provide a direct link, e.g.: `./main https://music.apple.com/us/album/whenever-you-need-somebody-2022-remaster/1624945511`
5. Select tracks from an album:
   - Use `--select`, e.g.: `./main --select https://music.apple.com/us/album/whenever-you-need-somebody-2022-remaster/1624945511`
6. Specify quality (e.g., AAC):
   - Use `--aac`, e.g.: `./main --aac https://music.apple.com/us/album/1989-taylors-version-deluxe/1713845538`
   - All available qualities are documented in `config-example.yaml` and below.
7. Inspect available quality:
   - Use `--debug`, e.g.: `./main --debug https://music.apple.com/us/album/1989-taylors-version-deluxe/1713845538`

## Interactive (Cobra CLI + Wizard)
1. Build the project:
   - `go build .`
2. Start the decryption helper: [wrapper](https://github.com/WorldObservationLog/wrapper).
3. Run the program in interactive mode:
   - `./main wizard` (Windows: `./main.exe wizard`)

## Notes
- Retry behavior:
  - After each download, a summary of warnings/errors is printed.
  - If failures exist, the wizard or command mode will ask whether to retry failed items.
- Non-interactive mode prints a single-line summary.
- Wizard mode prints dynamic information only after a menu option is selected, avoiding noise during navigation.
- The REPL mode has been removed; corresponding code and documentation have been cleaned.

## Thanks
Special thanks to `chocomint` for `agent-arm64.js`.
Original script authors and all contributors.