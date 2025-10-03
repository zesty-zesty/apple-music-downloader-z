# Apple Music Downloader

一个使用 Go 编写的 Apple Music 下载工具，可下载歌曲、专辑、播放列表与MV等。

## 环境准备
- 安装并配置 `PATH`：
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
- MV下载并使用 `MP4Box` 混流（音频+视频）。

## 命令行（非交互）
1. 构建项目：
   - `go build .`
2. 启动解密助手：[wrapper](https://github.com/WorldObservationLog/wrapper)。
3. 运行项目：
   - `./main`（Windows:`./main.exe`）
4. 下载单曲/专辑/播放列表：
   - 直接输入链接，如：`./main https://music.apple.com/us/album/whenever-you-need-somebody-2022-remaster/1624945511`
5. 专辑选择曲目：
   - 使用`select`，如：`./main --select https://music.apple.com/us/album/whenever-you-need-somebody-2022-remaster/1624945511`
6. 指定质量下载（如：AAC）：
   - 使用`--aac`，如：`./main --aac https://music.apple.com/us/album/1989-taylors-version-deluxe/1713845538`
   - 所有可选的质量可见`config-example.yaml`或下文
7. 查看可用质量：
   - 使用`--debug`，如：`./main --debug https://music.apple.com/us/album/1989-taylors-version-deluxe/1713845538`

## 交互式（Cobra CLI + Wizard）
1. 构建项目：
   - `go build .`
2. 启动解密助手：[wrapper](https://github.com/WorldObservationLog/wrapper)。
3. 运行项目，以交互式模式运行：
   - `./main wizard`（Windows:`./main.exe wizard`）

## 说明
- 重试行为：
   - 每次下载后会输出告警/错误汇总。
   - 如存在失败项，会在向导或命令模式下询问是否重试失败项。
- 非交互模式会打印单行汇总。
- 向导模式在选中菜单项后才打印动态信息，避免导航时的干扰信息。
- 本项目已移除 REPL 模式，对应代码与文档均已清理。

## 致谢
特别感谢 `chocomint` 提供 `agent-arm64.js`。
脚本原作者和所有Contributors。