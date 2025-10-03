# Apple Music Downloader

一个使用 Go 编写的 Apple Music 下载工具，可下载歌曲、专辑、播放列表与电台。

## 环境准备
- 安装并加入 `PATH`：
  - `MP4Box`（用于标签与 MV 混流）：https://gpac.io/downloads/gpac-nightly-builds/
  - `mp4decrypt`（MV 解密必需）：https://www.bento4.com/downloads/
  - `ffmpeg`（可选，用于动图封面与下载后转换）
- 在下载前确保解密助手 [wrapper](https://github.com/WorldObservationLog/wrapper) 已启动。
- 需要 Go `1.23+` 进行构建与运行。

## 功能
- 内嵌封面，导出/保存 LRC 歌词（需 `media-user-token`）。
- 支持分节/翻译歌词（Beta），含发音处理。
- 通过 URL 下载：专辑、播放列表、歌曲、艺人、电台。
- 向导式交互（Wizard）：菜单驱动的 `搜索`/`rip` 与设置。
- 下载过程中解密，减少大文件占用内存。
- MV 下载并使用 `MP4Box` 混流（音频+视频）。

## 命令行（非交互）
1. 启动解密助手：[wrapper](https://github.com/WorldObservationLog/wrapper)。
2. 下载专辑：
   - `go run main.go https://music.apple.com/us/album/whenever-you-need-somebody-2022-remaster/1624945511`
3. 下载单曲：
   - `go run main.go --song https://music.apple.com/us/album/never-gonna-give-you-up-2022-remaster/1624945511?i=1624945512`
   - 或 `go run main.go https://music.apple.com/us/song/you-move-me-2022-remaster/1624945520`
4. 专辑选择曲目：
   - `go run main.go --select https://music.apple.com/us/album/whenever-you-need-somebody-2022-remaster/1624945511`
5. 下载播放列表：
   - `go run main.go https://music.apple.com/us/playlist/taylor-swift-essentials/pl.3950454ced8c45a3b0cc693c2a7db97b`
6. Dolby Atmos：
   - `go run main.go --atmos https://music.apple.com/us/album/1989-taylors-version-deluxe/1713845538`
7. AAC：
   - `go run main.go --aac https://music.apple.com/us/album/1989-taylors-version-deluxe/1713845538`
8. 查看可用质量：
   - `go run main.go --debug https://music.apple.com/us/album/1989-taylors-version-deluxe/1713845538`

## 交互式（Cobra CLI + Wizard）
命令：
- `amd rip <url...>`：下载指定 URL（album|playlist|station|song）。支持多个 URL。
- `amd search <album|song|artist> <keywords>`：搜索并选择后下载。
- `amd concurrency <N>`：设置下载并发线程数（运行时）。
- `amd codec-priority show` 或 `amd codec-priority <alac,mp4a.40.2,ec-3>`：显示或设置编码优先级（运行时）。
- `amd flags`：显示当前下载相关标志与配置。

全局参数（可与任意命令一起使用）：
- `--atmos`、`--aac`、`--select`、`--song`、`--all-album`、`--debug`
- `--alac-max`、`--atmos-max`、`--aac-type`、`--mv-audio-type`、`--mv-max`、`--codec-priority`

向导入口与菜单（在 `amd` 顶层命令中）：
- 顶层菜单包含：`rip 单 URL`、`search 搜索下载`、`设置`、`帮助`、`退出`。
- `设置` 子菜单项使用静态名称，选中后才打印当前值：
  - 并发下载线程：选中后打印当前并发。
  - 编码优先级：菜单项静态；选中“查看当前优先级”后打印运行时与配置值。
  - 输出目录：选中后打印当前输出目录。

重试行为：
- 每次下载后会输出告警/错误汇总。
- 如存在失败项，会在向导或命令模式下询问是否重试失败项。

## 配置
复制 `config-example.yaml` 为 `config.yaml` 并编辑：
- `media-user-token`：歌词与 `aac-lc` 必需。
- `authorization-token`：自动获取，保持默认即可。
- `language`：店面语言；需符合支持列表。
- `output-folder`：所有下载的根目录。
- `download-concurrency`：并发下载上限。
- 封面/歌词嵌入选项：`cover-*`、`embed-*`。
- 下载后转换选项：`convert-*`、`ffmpeg-path`。

## 说明
- 非交互模式会打印单行汇总。
- 向导模式在选中菜单项后才打印动态信息，避免导航时的干扰信息。
- 本项目已移除 REPL 模式，对应代码与文档均已清理。

## 致谢
特别感谢 `chocomint` 提供 `agent-arm64.js`。