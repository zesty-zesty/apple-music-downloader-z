简体中文 / [English](./README.md)

# Apple Music 下载器（中文说明）

这是一个用Golang编写的小程序，支持下载 Apple Music 的单曲、专辑、MV、歌词与封面等信息，并提供交互式搜索与失败重试功能。

## 必备依赖
- 请先安装并将 `MP4Box` 加入环境变量：<https://gpac.io/downloads/gpac-nightly-builds/>
- 下载 MV 需要安装 `mp4decrypt`：<https://www.bento4.com/downloads/>
- 复制`config-example.yaml`为`config.yaml`，按需要更改配置文件。
- 若需获取 `aac-lc`、`MV 音频`、`歌词`等，需在 `config.yaml` 中填写你的 `media-user-token`（浏览器登录 Apple Music 后在开发者工具中获取cookie）。

## 支持功能
- 多线程下载&解密Apple Music的单曲、专辑、MV、歌词与封面等信息。
   - 支持的编解码与资源
      1. 音频：`alac`（含采样率变体）、`ec-3 / ac-3`（Dolby Atmos）、`aac / aac-lc / aac-binaural / aac-downmix`
      2. MV：需 `mp4decrypt` 与 `media-user-token`（用于提取音轨与写入标签）。
      3. 歌词：支持逐字、逐句、翻译/发音（Beta，需要在 `config.yaml` 配置相应语言参数与 `media-user-token`）。
- REPL交互式搜索与选择功能。
- 支持内嵌封面和LRC歌词（需要`media-user-token`）
- 支持获取逐词与未同步歌词

## 使用方法（命令行）
1. 确认解密程序 [wrapper](https://github.com/WorldObservationLog/wrapper) 已启动（用于下载时解密）。
2. 下载专辑：
   - `go run main.go https://music.apple.com/us/album/whenever-you-need-somebody-2022-remaster/1624945511`
3. 下载单曲：
   - 方式一：`go run main.go --song https://music.apple.com/us/album/never-gonna-give-you-up-2022-remaster/1624945511?i=1624945512`
   - 方式二：`go run main.go https://music.apple.com/us/song/you-move-me-2022-remaster/1624945520`
4. 选择性下载专辑中的部分曲目：
   - `go run main.go --select https://music.apple.com/us/album/whenever-you-need-somebody-2022-remaster/1624945511`
   - 根据提示输入曲目编号（空格分隔）。
5. 下载歌单：
   - `go run main.go https://music.apple.com/us/playlist/taylor-swift-essentials/pl.3950454ced8c45a3b0cc693c2a7db97b`
6. 下载电台（需要 `media-user-token`）：
   - `go run main.go https://music.apple.com/us/station/…`
7. 下载艺人所有专辑：
   - `go run main.go https://music.apple.com/us/artist/taylor-swift/159260351 --all-album`
8. 下载 Dolby Atmos：
   - `go run main.go --atmos https://music.apple.com/us/album/1989-taylors-version-deluxe/1713845538`
9. 下载 AAC：
   - `go run main.go --aac https://music.apple.com/us/album/1989-taylors-version-deluxe/1713845538`
10. 查看清晰度信息（调试）：
   - `go run main.go --debug https://music.apple.com/us/album/1989-taylors-version-deluxe/1713845538`

### 非交互批处理的失败重试
- 直接传入 URL 列表运行完成后（如 `go run main.go <url1> <url2> …`），若检测到错误，程序将询问是否“重试失败项”。
- 选择重试后：
  - 开启 `retryOnly=true` 进行第二轮，仅重试失败曲目；结束后重置为 `retryOnly=false` 并清空失败集合。

## 使用方法（交互模式）
1. 确认解密程序 [wrapper](https://github.com/WorldObservationLog/wrapper) 已启动（用于下载时解密）。
2. 进入交互式模式（默认）：
   - `go run main.go`
   - 输入 `help` 查看命令；输入 `exit` 或`quit`退出。

### 交互模式命令
- `rip <url>`：下载指定 URL 对应的曲目/专辑/歌单/电台。
- `search <type> <kw>`：交互式搜索并选择。
  - `<type>`：`song`、`album`、`artist` 等。
  - `<kw>`：搜索关键词。
- `concurrency <N>`：设置下载并发线程数
- `codec-priority <list>`：设置编码优先级，逗号分隔，如: alac,mp4a.40.2,ec-3
- `codec-priority show`：显示当前运行时与配置的编码优先级
- `flags`：显示当前下载相关标志
- `exit` 或 `quit`：退出交互式模式。

### 交互模式的失败重试
- 通过 `rip <url>` 或 `search <type> <kw>` 执行下载后，程序将输出详细的告警/错误汇总，并询问“是否重试失败项”。
- 选择重试后：
  - 程序会开启 `retryOnly=true`，仅重试记录中的失败曲目（无失败记录则回退全量）。
  - 本轮结束后自动重置为 `retryOnly=false` 并清空失败记录，以避免影响后续任务。

## 其它说明
- 整张专辑/歌单/电台在进入阶段失败的情况也会被记录，便于后续扩展；当前“仅重试失败项”以曲目为单位进行。
- 交互模式下可通过 `concurrency <N>` 设置下载并发，`codec-priority <list>` 设置运行期编解码优先级，`flags` 查看当前标志位。交互模式下的设置仅本次运行有效，不保存到`config.yaml`中。
- 特别感谢 chocomint 创建 agent-arm64.js