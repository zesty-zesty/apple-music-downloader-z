package main

import (
    "fmt"
    "os"
    "strings"

    survey "github.com/AlecAivazis/survey/v2"
    "github.com/manifoldco/promptui"
    "github.com/spf13/cobra"
    yaml "gopkg.in/yaml.v2"
)

func init() {
    wizardCmd := &cobra.Command{
        Use:   "wizard",
        Short: "交互式向导（PromptUI）",
        Run: func(cmd *cobra.Command, args []string) {
            token := cliToken
            for {
                actions := []string{"rip 单 URL", "search 搜索下载", "设置", "帮助", "退出"}
                var action string
                if err := survey.AskOne(&survey.Select{Message: "选择操作", Options: actions}, &action); err != nil {
                    return
                }

                switch action {
                case "rip 单 URL":
                    urlPrompt := promptui.Prompt{Label: "请输入 Apple Music URL", Templates: &promptui.PromptTemplates{Success: ""}}
                    url, err := urlPrompt.Run()
                    if err != nil || strings.TrimSpace(url) == "" {
                        fmt.Println("未输入 URL，已取消")
                        continue
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

                case "search 搜索下载":
                    types := []string{"album", "song", "artist"}
                    var st string
                    if err := survey.AskOne(&survey.Select{Message: "选择搜索类型", Options: types}, &st); err != nil {
                        continue
                    }

                    kwPrompt := promptui.Prompt{Label: "输入关键词", Templates: &promptui.PromptTemplates{Success: ""}}
                    kwStr, err := kwPrompt.Run()
                    if err != nil || strings.TrimSpace(kwStr) == "" {
                        fmt.Println("未输入关键词，已取消")
                        continue
                    }

                    selectedUrl, err := handleSearch(st, strings.Fields(kwStr), token)
                    if err != nil {
                        fmt.Println("Search error:", err)
                        continue
                    }
                    if selectedUrl == "" {
                        fmt.Println("未选择内容")
                        continue
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

                case "设置":
                    runSettingsMenu()

                case "帮助":
                    printWizardHelp()

                case "退出":
                    return
                }
            }
        },
    }
    rootCmd.AddCommand(wizardCmd)
}

// 设置一级菜单
func runSettingsMenu() {
    for {
        items := []string{
            "并发下载线程",
            "编码优先级 codec-priority",
            "输出目录 output-folder",
            "返回",
        }
        var choose string
        if err := survey.AskOne(&survey.Select{Message: "设置", Options: items}, &choose); err != nil {
            return
        }
        switch choose {
        case items[0]:
            fmt.Printf("当前并发: %d\n", DownloadConcurrency)
            concPrompt := promptui.Prompt{Label: "请输入并发数(整数)", Default: fmt.Sprintf("%d", DownloadConcurrency), Templates: &promptui.PromptTemplates{Success: ""}}
            concStr, err := concPrompt.Run()
            if err == nil {
                concStr = strings.TrimSpace(concStr)
                var v int
                _, _ = fmt.Sscanf(concStr, "%d", &v)
                if v > 0 {
                    DownloadConcurrency = v
                    fmt.Println("并发数已更新为:", v)
                }
            }
        case items[1]:
            runCodecPriorityMenu()
        case items[2]:
            fmt.Printf("当前输出目录: %s\n", strings.TrimSpace(OutputFolder))
            outPrompt := promptui.Prompt{Label: "设置输出目录(路径)", Default: strings.TrimSpace(OutputFolder), Templates: &promptui.PromptTemplates{Success: ""}}
            outStr, err := outPrompt.Run()
            if err == nil && strings.TrimSpace(outStr) != "" {
                OutputFolder = strings.TrimSpace(outStr)
                fmt.Println("输出目录(运行时)已更新为:", OutputFolder)
                if askYesNo("是否写入 config.yaml 的 output-folder? (y/N) ") {
                    if err := writeConfigKey("output-folder", OutputFolder); err != nil {
                        fmt.Println("写入配置失败:", err)
                    } else {
                        fmt.Println("config.yaml 已更新 output-folder")
                    }
                }
            }
        case items[3]:
            return
        }
    }
}

// 编码优先级菜单
func runCodecPriorityMenu() {
    for {
        items := []string{
            "查看当前优先级",
            "设置运行时优先级(逗号分隔)",
            "清除运行时覆盖(恢复配置)",
            "写入配置文件(codec-priority)",
            "返回",
        }
        var choose string
        if err := survey.AskOne(&survey.Select{Message: "编码优先级", Options: items}, &choose); err != nil {
            return
        }
        switch choose {
        case items[0]:
            fmt.Println("运行时优先级:", strings.Join(currentCodecPriority(), ","))
            fmt.Println("配置优先级:", strings.Join(Config.CodecPriority, ","))
        case items[1]:
            prompt := promptui.Prompt{Label: "输入优先级(例如: alac,mp4a.40.2,ec-3)", Templates: &promptui.PromptTemplates{Success: ""}}
            s, err := prompt.Run()
            if err == nil {
                parts := strings.Split(s, ",")
                var list []string
                for _, p := range parts {
                    t := strings.TrimSpace(p)
                    if t != "" {
                        list = append(list, t)
                    }
                }
                if len(list) > 0 {
                    RuntimeCodecPriority = list
                    fmt.Println("运行时编码优先级已更新:", strings.Join(list, ","))
                }
            }
        case items[2]:
            RuntimeCodecPriority = nil
            fmt.Println("已清除运行时覆盖，恢复为配置优先级。")
        case items[3]:
            src := currentCodecPriority()
            if len(src) == 0 {
                fmt.Println("当前优先级为空，未写入。")
                continue
            }
            if err := writeConfigKey("codec-priority", src); err != nil {
                fmt.Println("写入配置失败:", err)
            } else {
                Config.CodecPriority = src
                fmt.Println("config.yaml 已更新 codec-priority")
            }
        case items[4]:
            return
        }
    }
}

// 帮助
func printWizardHelp() {
    fmt.Println("向导命令概览：")
    fmt.Println("- rip 单 URL: 输入 Apple Music URL 进行下载")
    fmt.Println("- search 搜索下载: 选择类型与关键词后下载")
    fmt.Println("- 设置: 管理并发、编码优先级、输出目录等")
    fmt.Println("  - 并发下载线程: 设置下载并发数")
    fmt.Println("  - 编码优先级: 查看/设置运行时优先级，写入配置")
    fmt.Println("  - 输出目录: 设置运行时输出目录，支持写入配置")
    fmt.Println("- 退出: 结束向导")
}

// 写入配置文件指定键
func writeConfigKey(key string, value interface{}) error {
    data, err := os.ReadFile("config.yaml")
    if err != nil {
        return err
    }
    var m map[string]interface{}
    if err := yaml.Unmarshal(data, &m); err != nil {
        return err
    }
    m[key] = value
    out, err := yaml.Marshal(m)
    if err != nil {
        return err
    }
    return os.WriteFile("config.yaml", out, 0644)
}