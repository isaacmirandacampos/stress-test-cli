/*
Copyright © 2025 Isaac de Miranda Campos

*/
package cmd

import (
	`fmt`
	"os"
	`sync`
	`time`

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
		}
		requests, err := cmd.Flags().GetInt("requests")
		if err != nil {
			fmt.Println(err)
		}
		concurrency, err := cmd.Flags().GetInt("concurrency")
		if err != nil {
			fmt.Println(err)
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
				res, err := request(url)
				if err != nil {
					fmt.Println(err)
					return
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
