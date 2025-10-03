package main

import (
    "fmt"
    "strings"

    "github.com/spf13/cobra"
)

func init() {
    searchCmd := &cobra.Command{
        Use:   "search <album|song|artist> <keywords>",
        Short: "搜索并选择后下载",
        Args:  cobra.MinimumNArgs(2),
        Run: func(cmd *cobra.Command, args []string) {
            st := strings.ToLower(args[0])
            kw := args[1:]
            selectedUrl, err := handleSearch(st, kw, cliToken)
            if err != nil {
                fmt.Println("Search error:", err)
                return
            }
            if selectedUrl == "" {
                fmt.Println("No selection.")
                return
            }
            clearIssues()
            handleSingleURL(selectedUrl, cliToken)
            // 完成后显示详细告警/错误信息并支持重试
            printIssuesSummary()
            for counter.Error > 0 {
                if !askYesNo("是否重试失败项? (y/N) ") {
                    break
                }
                clearIssues()
                retryOnly = true
                handleSingleURL(selectedUrl, cliToken)
                retryOnly = false
                printIssuesSummary()
            }
            clearFail()
            fmt.Printf("=======  [\u2714 ] Completed: %d/%d  |  [\u26A0 ] Warnings: %d  |  [\u2716 ] Errors: %d  =======\n", counter.Success, counter.Total, counter.Unavailable+counter.NotSong, counter.Error)
        },
    }
    rootCmd.AddCommand(searchCmd)
}