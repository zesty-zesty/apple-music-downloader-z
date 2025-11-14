package main

import (
    "fmt"
    "os"
    "runtime/pprof"
    "strings"
    "log/slog"

    "github.com/spf13/cobra"
    "main/utils/ampapi"
)

var (
    // 提前初始化 rootCmd，确保在各子命令的 init() 执行前可用
    rootCmd  = &cobra.Command{
        Use:   "amd",
        Short: "Apple Music Downloader CLI",
        PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
            // 将 CLI 标志应用到运行时配置
            if cmd.Flags().Changed("atmos") {
                v, _ := cmd.Flags().GetBool("atmos")
                dl_atmos = v
            }
            if cmd.Flags().Changed("aac") {
                v, _ := cmd.Flags().GetBool("aac")
                dl_aac = v
            }
            if cmd.Flags().Changed("select") {
                v, _ := cmd.Flags().GetBool("select")
                dl_select = v
            }
            if cmd.Flags().Changed("song") {
                v, _ := cmd.Flags().GetBool("song")
                dl_song = v
            }
            if cmd.Flags().Changed("all-album") {
                v, _ := cmd.Flags().GetBool("all-album")
                artist_select = v
            }
            if cmd.Flags().Changed("debug") {
                v, _ := cmd.Flags().GetBool("debug")
                debug_mode = v
            }

            // 初始化日志
            var logLevel, logFormat, logFile string
            var noColor bool
            if f := cmd.Flags().Lookup("log-level"); f != nil {
                logLevel, _ = cmd.Flags().GetString("log-level")
            }
            if f := cmd.Flags().Lookup("log-format"); f != nil {
                logFormat, _ = cmd.Flags().GetString("log-format")
            }
            if f := cmd.Flags().Lookup("log-file"); f != nil {
                logFile, _ = cmd.Flags().GetString("log-file")
            }
            if f := cmd.Flags().Lookup("no-color"); f != nil {
                noColor, _ = cmd.Flags().GetBool("no-color")
            }

            // 如果开启 debug，则提升到 debug 日志级别
            if debug_mode && (logLevel == "" || logLevel == "info") {
                logLevel = "debug"
            }

            Logger = SetupLogger(logLevel, logFormat, logFile, noColor)
            if Logger != nil {
                Logger.Info("logger initialized", slog.String("level", logLevel), slog.String("format", logFormat), slog.String("file", logFile))
            }

            if cmd.Flags().Changed("alac-max") {
                v, _ := cmd.Flags().GetInt("alac-max")
                Config.AlacMax = v
            }
            if cmd.Flags().Changed("atmos-max") {
                v, _ := cmd.Flags().GetInt("atmos-max")
                Config.AtmosMax = v
            }
            if cmd.Flags().Changed("aac-type") {
                v, _ := cmd.Flags().GetString("aac-type")
                Config.AacType = v
            }
            if cmd.Flags().Changed("mv-audio-type") {
                v, _ := cmd.Flags().GetString("mv-audio-type")
                Config.MVAudioType = v
            }
            if cmd.Flags().Changed("mv-max") {
                v, _ := cmd.Flags().GetInt("mv-max")
                Config.MVMax = v
            }
            if cmd.Flags().Changed("codec-priority") {
                list, _ := cmd.Flags().GetString("codec-priority")
                parts := strings.Split(list, ",")
                out := make([]string, 0, len(parts))
                for _, p := range parts {
                    p = strings.TrimSpace(p)
                    if p != "" {
                        out = append(out, p)
                    }
                }
                Config.CodecPriority = out
            }

            // 启用 CPU Profiling（用于 PGO）
            if cmd.Flags().Changed("profile-cpu") {
                cpuProfilePath, _ = cmd.Flags().GetString("profile-cpu")
                if cpuProfilePath != "" && !cpuProfileActive {
                    f, err := os.Create(cpuProfilePath)
                    if err != nil {
                        return fmt.Errorf("failed to create profile file %s: %w", cpuProfilePath, err)
                    }
                    if err := pprof.StartCPUProfile(f); err != nil {
                        _ = f.Close()
                        return fmt.Errorf("failed to start CPU profile: %w", err)
                    }
                    cpuProfileFile = f
                    cpuProfileActive = true
                }
            }

            // 获取 token
            var err error
            cliToken, err = ampapi.GetToken()
            if err != nil {
                if Config.AuthorizationToken != "" && Config.AuthorizationToken != "your-authorization-token" {
                    cliToken = strings.Replace(Config.AuthorizationToken, "Bearer ", "", -1)
                } else {
                    return fmt.Errorf("failed to get token: %w", err)
                }
            }
            return nil
        },
        PersistentPostRun: func(cmd *cobra.Command, args []string) {
            // 结束 CPU Profiling
            if cpuProfileActive {
                pprof.StopCPUProfile()
                if cpuProfileFile != nil {
                    _ = cpuProfileFile.Close()
                }
                cpuProfileActive = false
            }
        },
    }
    cliToken string
)

func init() {
    // 持久化标志（全局）
    rootCmd.PersistentFlags().BoolVar(&dl_atmos, "atmos", false, "Enable atmos download mode")
    rootCmd.PersistentFlags().BoolVar(&dl_aac, "aac", false, "Enable adm-aac download mode")
    rootCmd.PersistentFlags().BoolVar(&dl_select, "select", false, "Enable selective download")
    rootCmd.PersistentFlags().BoolVar(&dl_song, "song", false, "Enable single song download mode")
    rootCmd.PersistentFlags().BoolVar(&artist_select, "all-album", false, "Download all artist albums")
    rootCmd.PersistentFlags().BoolVar(&debug_mode, "debug", false, "Enable debug mode to show audio quality information")
    // 日志控制
    rootCmd.PersistentFlags().String("log-level", "info", "Log level: debug, info, warn, error")
    rootCmd.PersistentFlags().String("log-format", "text", "Log format: text, json, auto")
    rootCmd.PersistentFlags().String("log-file", "", "Log file path (enable file logging)")
    rootCmd.PersistentFlags().Bool("no-color", false, "Disable color in console output")
    rootCmd.PersistentFlags().IntVar(&Config.AlacMax, "alac-max", Config.AlacMax, "Specify the max quality for download alac")
    rootCmd.PersistentFlags().IntVar(&Config.AtmosMax, "atmos-max", Config.AtmosMax, "Specify the max quality for download atmos")
    rootCmd.PersistentFlags().StringVar(&Config.AacType, "aac-type", Config.AacType, "Select AAC type, aac aac-binaural aac-downmix aac-lc")
    rootCmd.PersistentFlags().StringVar(&Config.MVAudioType, "mv-audio-type", Config.MVAudioType, "Select MV audio type, atmos ac3 aac")
    rootCmd.PersistentFlags().IntVar(&Config.MVMax, "mv-max", Config.MVMax, "Specify the max quality for download MV")
    rootCmd.PersistentFlags().String("codec-priority", strings.Join(Config.CodecPriority, ","), "Specify codec priority, comma separated")
    rootCmd.PersistentFlags().StringVar(&cpuProfilePath, "profile-cpu", "", "生成 CPU Profile（pprof），用于 PGO，例如 default.pgo")

    // 绑定 aac_type 指针到配置，避免 setDlFlags() 写入空指针
    aac_type = &Config.AacType
}

// Execute 作为 Cobra 入口
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        if Logger != nil {
            Logger.Error("command execution failed", slog.Any("error", err))
        } else {
            fmt.Println(err)
        }
    }
}