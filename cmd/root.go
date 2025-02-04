/*
Copyright © 2025 Isaac de Miranda Campos
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "stress-test-cli",
	Short: "A cli to do stress test",
	Long:  `A cli to do a stress test. Setting up the --url <target> to work`,
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		url, err := cmd.Flags().GetString("url")
		if err != nil {
			fmt.Println(err)
			return
		}
		if url == "" {
			fmt.Println("url is required")
			return
		}
		requests, err := cmd.Flags().GetInt("requests")
		if err != nil {
			fmt.Println(err)
			return
		}
		concurrency, err := cmd.Flags().GetInt("concurrency")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("url: %s, requests: %d, concurrency: %d\n", url, requests, concurrency)

		channel := make(chan struct{}, concurrency)
		var wg sync.WaitGroup
		statusCodes := sync.Map{}
		for i := 0; i < requests; i++ {
			wg.Add(1)
			channel <- struct{}{}
			go func(url string) {
				defer wg.Done()
				defer func() { <-channel }()
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				res, err := request(&ctx, url)
				if err != nil {
					if errors.Is(ctx.Err(), context.DeadlineExceeded) {
						summaryStatusCode(&statusCodes, 408)
						return
					}
					summaryStatusCode(&statusCodes, 408)
				}
				summaryStatusCode(&statusCodes, res.StatusCode)
			}(url)
		}
		wg.Wait()
		elapsed := time.Since(start)
		report(elapsed, requests, &statusCodes)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("url", "u", "", "Define the target url")
	rootCmd.Flags().IntP("requests", "r", 100, "Define the number of requests")
	rootCmd.Flags().IntP("concurrency", "c", 1, "Define the number of concurrent requests")
}
