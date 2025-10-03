English / [简体中文](./README-CN.md)

# Apple Music Downloader
A command-line tool to download songs, albums, playlists, and stations from Apple Music writting in Go.

## Prerequisites
- Install and add to `PATH`:
  - `MP4Box` (used for tagging and MV muxing): https://gpac.io/downloads/gpac-nightly-builds/
  - `mp4decrypt` (required for MV decryption): https://www.bento4.com/downloads/
  - `ffmpeg` (optional, used for animated artwork and post-download conversion)
- Ensure the decryption helper [wrapper](https://github.com/WorldObservationLog/wrapper) is running before downloads.
- Go `1.23+` to build and run.

## Features
- Embed cover art and export/save LRC lyrics (requires `media-user-token`).
- Support syllable/translation lyrics (beta), including pronunciation handling.
- Download by URL: album, playlist, song, artist, station.
- Interactive wizard with menu-driven `search`/`rip` and settings.
- Decrypt while downloading to reduce memory usage with large files.
- MV download and mux with `MP4Box` (audio+video).

## Supported audio types
- `alac` (audio-alac-stereo)
- `ec3` (audio-atmos / audio-ec3)
- `aac` (audio-stereo)
- `aac-lc` (audio-stereo)
- `aac-binaural` (audio-stereo-binaural)
- `aac-downmix` (audio-stereo-downmix)

For `aac-lc`, MV and lyrics, you must configure a valid `media-user-token`.

## Usage (non-interactive)
1. Start the decryption helper: [wrapper](https://github.com/WorldObservationLog/wrapper).
2. Download an album:
   - `go run main.go https://music.apple.com/us/album/whenever-you-need-somebody-2022-remaster/1624945511`
3. Download a single song:
   - `go run main.go --song https://music.apple.com/us/album/never-gonna-give-you-up-2022-remaster/1624945511?i=1624945512`
   - or `go run main.go https://music.apple.com/us/song/you-move-me-2022-remaster/1624945520`
4. Select tracks from an album:
   - `go run main.go --select https://music.apple.com/us/album/whenever-you-need-somebody-2022-remaster/1624945511`
5. Download a playlist:
   - `go run main.go https://music.apple.com/us/playlist/taylor-swift-essentials/pl.3950454ced8c45a3b0cc693c2a7db97b`
6. Dolby Atmos:
   - `go run main.go --atmos https://music.apple.com/us/album/1989-taylors-version-deluxe/1713845538`
7. AAC:
   - `go run main.go --aac https://music.apple.com/us/album/1989-taylors-version-deluxe/1713845538`
8. Inspect available quality:
   - `go run main.go --debug https://music.apple.com/us/album/1989-taylors-version-deluxe/1713845538`

## Usage (CLI with Cobra)
Commands:
- `amd rip <url...>`: 下载指定 URL (album|playlist|station|song)。支持多个 URL。
- `amd search <album|song|artist> <keywords>`: 搜索并选择后下载。
- `amd concurrency <N>`: 设置下载并发线程数。
- `amd codec-priority show` 或 `amd codec-priority <alac,mp4a.40.2,ec-3>`: 显示或设置编码优先级。
- `amd flags`: 显示当前下载相关标志与配置。

Global flags (可与任意命令一起使用):
- `--atmos`、`--aac`、`--select`、`--song`、`--all-album`、`--debug`
- `--alac-max`、`--atmos-max`、`--aac-type`、`--mv-audio-type`、`--mv-max`、`--codec-priority`

Examples:
- `go run . rip https://music.apple.com/us/album/.../1624945511`
- `go run . search album "1989 taylor"`
- `go run . rip --atmos https://music.apple.com/us/album/...`
- `go run . concurrency 8`

Retry behavior:
- 任务结束后会输出告警/错误汇总。
- 如存在失败项，会询问是否重试失败项。

## Configuration
Copy `config-example.yaml` to `config.yaml` and edit:
- `media-user-token`: required for lyrics and `aac-lc`.
- `authorization-token`: auto-obtained; leave default unless needed.
- `language`: storefront language; see supported languages list.
- `output-folder`: root output directory for all downloads.
- `download-concurrency`: maximum parallel downloads.
- `cover-*`, `embed-*`: artwork/lyrics embedding options.
- Conversion options (`convert-*`, `ffmpeg-path`) for post-download conversion.

## Lyrics: translation and pronunciation (Beta)
1. Open [Apple Music](https://beta.music.apple.com) and log in.
2. In DevTools, open `Network` and search a song supporting translation/pronunciation.
3. Refresh (`Ctrl+R`), play the song, and click the lyric button.
4. Find `syllable-lyrics` in requests and open it.
5. Copy the `l=` language values from the request: `.../syllable-lyrics?l=<values>&extend=ttmlLocalizations`.
6. Paste into `config.yaml`. If pronunciation is not needed, remove `%5D=<...>` from the value.

## Notes
- Non-interactive mode prints a one-line summary.
- Interactive wizard prints dynamic info (e.g., current values) only after choosing a menu item.
- After each download, the wizard summarizes warnings/errors and can retry failed items.

## Thanks
Special thanks to `chocomint` for creating `agent-arm64.js`.
