简体中文 / [English](./README.md)

# Apple Music ALAC / Dolby Atmos 下载器（中文说明）

这是一个用Golang编写的小程序，支持下载 Apple Music 的单曲、专辑、MV、歌词与封面等信息，并提供交互式搜索与失败重试功能。

## 必备依赖
- 请先安装并将 `MP4Box` 加入环境变量：<https://gpac.io/downloads/gpac-nightly-builds/>
- 下载 MV 需要安装 `mp4decrypt`：<https://www.bento4.com/downloads/>
- 复制`config-example.yaml`为`config.yaml`，按需要更改配置文件。
- 若需获取 `aac-lc`、`MV 音频`、`歌词`等，需在 `config.yaml` 中填写你的 `media-user-token`（浏览器登录 Apple Music 后在开发者工具中获取cookie）。

## 新增/变更功能
- 失败项记录与“仅重试失败项”模式：
  - 程序运行时会记录失败曲目，按“实体ID（如专辑/歌单/电台）+曲目序号”保存。
  - 开启“仅重试失败项”（`retryOnly`）后，专辑/歌单/电台将优先只处理失败曲目；若该实体没有失败曲目记录，则回退到处理全部曲目。
  - 交互与批处理模式均支持一次任务结束后提示“是否重试失败项”，确认后自动开启仅重试失败模式并再次执行。
- 交互式 REPL：新增 `rip`、`search`、`concurrency`、`codec-priority`、`flags` 等命令，支持箭头选择搜索结果与动态调整下载并发、编解码优先级。
- 运行期编解码优先级：允许通过命令设置运行期优先级，下载过程将根据可用码流与优先级动态选择。

## 使用方法
1. 确认解密程序 [wrapper](https://github.com/zhaarey/wrapper) 已启动（用于下载时解密）。
2. 进入交互式模式（默认）：
   - `go run main.go`
   - 输入 `help` 查看命令；输入 `exit` 退出。
3. 下载专辑：
   - `go run main.go https://music.apple.com/us/album/whenever-you-need-somebody-2022-remaster/1624945511`
4. 下载单曲：
   - 方式一：`go run main.go --song https://music.apple.com/us/album/never-gonna-give-you-up-2022-remaster/1624945511?i=1624945512`
   - 方式二：`go run main.go https://music.apple.com/us/song/you-move-me-2022-remaster/1624945520`
5. 选择性下载专辑中的部分曲目：
   - `go run main.go --select https://music.apple.com/us/album/whenever-you-need-somebody-2022-remaster/1624945511`
   - 根据提示输入曲目编号（空格分隔）。
6. 下载歌单：
   - `go run main.go https://music.apple.com/us/playlist/taylor-swift-essentials/pl.3950454ced8c45a3b0cc693c2a7db97b`
7. 下载电台（需要 `media-user-token`）：
   - `go run main.go https://music.apple.com/us/station/…`
8. 下载艺人所有专辑：
   - `go run main.go https://music.apple.com/us/artist/taylor-swift/159260351 --all-album`
9. 下载 Dolby Atmos：
   - `go run main.go --atmos https://music.apple.com/us/album/1989-taylors-version-deluxe/1713845538`
10. 下载 AAC：
   - `go run main.go --aac https://music.apple.com/us/album/1989-taylors-version-deluxe/1713845538`
11. 查看清晰度信息（调试）：
   - `go run main.go --debug https://music.apple.com/us/album/1989-taylors-version-deluxe/1713845538`
12. 交互式搜索并选择：
   - `go run main.go --search [song|album|artist] "关键词"`

## 交互模式的失败重试
- 通过 `rip <url>` 或 `search <type> <kw>` 执行下载后，程序将输出详细的告警/错误汇总，并询问“是否重试失败项”。
- 选择重试后：
  - 程序会开启 `retryOnly=true`，仅重试记录中的失败曲目（无失败记录则回退全量）。
  - 本轮结束后自动重置为 `retryOnly=false` 并清空失败记录，以避免影响后续任务。

## 非交互批处理的失败重试
- 直接传入 URL 列表运行完成后（如 `go run main.go <url1> <url2> …`），若检测到错误，程序将询问是否“重试失败项”。
- 选择重试后：
  - 开启 `retryOnly=true` 进行第二轮，仅重试失败曲目；结束后重置为 `retryOnly=false` 并清空失败集合。

## 支持的编解码与资源
- 音频：`alac`（含采样率变体）、`ec-3 / ac-3`（Dolby Atmos）、`aac / aac-lc / aac-binaural / aac-downmix`
- MV：需 `mp4decrypt` 与 `media-user-token`（用于提取音轨与写入标签）。
- 歌词：支持逐字、逐句、翻译/发音（Beta，需要在 `config.yaml` 配置相应语言参数与 `media-user-token`）。

## 其它说明
- “实体级失败记录”（整张专辑/歌单/电台在进入阶段失败的情况）也会被记录，便于后续扩展；当前“仅重试失败项”以曲目为粒度进行。
- 交互模式下可通过 `concurrency <N>` 设置下载并发，`codec-priority <list>` 设置运行期编解码优先级，`flags` 查看当前标志位。
