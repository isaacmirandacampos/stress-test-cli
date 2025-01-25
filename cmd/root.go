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

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "stress-test-cli",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
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
		var statusCodeSuccess int32
		for i := 0; i < requests; i++ {
			wg.Add(1)
			channel <- struct{}{}
			go func(url string) {
				defer wg.Done()
				defer func() { <-channel }()
				request(url, &statusCodeSuccess)
			}(url)
		}
		wg.Wait()
		elapsed := time.Since(start)
		fmt.Printf("Total execution time: %s\n", elapsed)
		fmt.Printf("Total requests: %v\n", requests)
		fmt.Printf("Total Status Code 200: %d\n", statusCodeSuccess)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
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
