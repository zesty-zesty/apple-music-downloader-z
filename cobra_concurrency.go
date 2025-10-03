package main

import (
    "fmt"
    "strconv"

    "github.com/spf13/cobra"
)

func init() {
    concCmd := &cobra.Command{
        Use:   "concurrency <N>",
        Short: "设置下载并发线程数",
        Args:  cobra.ExactArgs(1),
        Run: func(cmd *cobra.Command, args []string) {
            n, err := strconv.Atoi(args[0])
            if err != nil || n <= 0 {
                fmt.Println("Invalid N.")
                return
            }
            DownloadConcurrency = n
            fmt.Printf("并发线程数设置为: %d\n", DownloadConcurrency)
        },
    }
    rootCmd.AddCommand(concCmd)
}