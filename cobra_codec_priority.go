package main

import (
    "fmt"
    "strings"

    "github.com/spf13/cobra"
)

func init() {
    codecCmd := &cobra.Command{
        Use:   "codec-priority [show | <codec1,codec2,...>]",
        Short: "显示或设置运行时编码优先级",
        Args:  cobra.RangeArgs(0, 1),
        Run: func(cmd *cobra.Command, args []string) {
            if len(args) == 0 || args[0] == "show" {
                fmt.Printf("runtime=%v config=%v\n", RuntimeCodecPriority, Config.CodecPriority)
                return
            }
            list := strings.Split(args[0], ",")
            out := make([]string, 0, len(list))
            for _, c := range list {
                c = strings.TrimSpace(c)
                if c != "" {
                    out = append(out, c)
                }
            }
            RuntimeCodecPriority = out
            fmt.Printf("运行时编码优先级设置为: %v\n", RuntimeCodecPriority)
        },
    }
    rootCmd.AddCommand(codecCmd)
}