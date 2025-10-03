package main

import (
    "fmt"
    "strings"

    "github.com/manifoldco/promptui"
    "github.com/spf13/cobra"
)

func init() {
    wizardCmd := &cobra.Command{
        Use:   "wizard",
        Short: "交互式向导（PromptUI）",
        Run: func(cmd *cobra.Command, args []string) {
            token := cliToken

            // 操作选择：rip / search
            actions := []string{"rip 单 URL", "search 搜索下载"}
            sel := promptui.Select{Label: "选择操作", Items: actions}
            actionIdx, _, err := sel.Run()
            if err != nil {
                return
            }

            // 并发度设置
            concItems := []string{fmt.Sprintf("保持当前并发: %d", DownloadConcurrency), "修改并发数"}
            concSel := promptui.Select{Label: "并发下载线程", Items: concItems}
            concIdx, _, err := concSel.Run()
            if err == nil && concIdx == 1 {
                concPrompt := promptui.Prompt{Label: "请输入并发数(整数)", Default: fmt.Sprintf("%d", DownloadConcurrency)}
                concStr, err := concPrompt.Run()
                if err == nil {
                    concStr = strings.TrimSpace(concStr)
                    var v int
                    _, _ = fmt.Sscanf(concStr, "%d", &v)
                    if v > 0 {
                        DownloadConcurrency = v
                    }
                }
            }

            switch actionIdx {
            case 0: // rip 单 URL
                urlPrompt := promptui.Prompt{Label: "请输入 Apple Music URL"}
                url, err := urlPrompt.Run()
                if err != nil || strings.TrimSpace(url) == "" {
                    fmt.Println("未输入 URL，已取消")
                    return
                }
                url = strings.TrimSpace(url)

                clearIssues()
                handleSingleURL(url, token)
                printIssuesSummary()

                for counter.Error > 0 {
                    if !askYesNo("是否重试失败项? (y/N) ") {
                        break
                    }
                    clearIssues()
                    retryOnly = true
                    handleSingleURL(url, token)
                    retryOnly = false
                    printIssuesSummary()
                }
                clearFail()
                fmt.Printf("=======  [✔ ] Completed: %d/%d  |  [⚠ ] Warnings: %d  |  [✘ ] Errors: %d  =======\n", counter.Success, counter.Total, counter.Unavailable+counter.NotSong, counter.Error)

            case 1: // search 搜索下载
                types := []string{"album", "song", "artist"}
                typeSel := promptui.Select{Label: "选择搜索类型", Items: types}
                typeIdx, _, err := typeSel.Run()
                if err != nil {
                    return
                }
                st := types[typeIdx]

                kwPrompt := promptui.Prompt{Label: "输入关键词"}
                kwStr, err := kwPrompt.Run()
                if err != nil || strings.TrimSpace(kwStr) == "" {
                    fmt.Println("未输入关键词，已取消")
                    return
                }

                selectedUrl, err := handleSearch(st, strings.Fields(kwStr), token)
                if err != nil {
                    fmt.Println("Search error:", err)
                    return
                }
                if selectedUrl == "" {
                    fmt.Println("未选择内容")
                    return
                }

                clearIssues()
                handleSingleURL(selectedUrl, token)
                printIssuesSummary()

                for counter.Error > 0 {
                    if !askYesNo("是否重试失败项? (y/N) ") {
                        break
                    }
                    clearIssues()
                    retryOnly = true
                    handleSingleURL(selectedUrl, token)
                    retryOnly = false
                    printIssuesSummary()
                }
                clearFail()
                fmt.Printf("=======  [✔ ] Completed: %d/%d  |  [⚠ ] Warnings: %d  |  [✘ ] Errors: %d  =======\n", counter.Success, counter.Total, counter.Unavailable+counter.NotSong, counter.Error)
            }
        },
    }
    rootCmd.AddCommand(wizardCmd)
}