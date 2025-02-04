package cmd

import (
	"fmt"
	"sync"
	"time"
)

func report(executionTime time.Duration, numberOfRequests int, statusCodes *sync.Map) {
	fmt.Printf("=========================================\n")
	fmt.Printf("Total execution time: %s\n", executionTime)
	fmt.Printf("Total requests: %v\n", numberOfRequests)
	fmt.Printf("\nSummary: \n")
	statusCodes.Range(func(key, value interface{}) bool {
		count := value.(int32)
		percentage := float64(count) * 100 / float64(numberOfRequests)
		fmt.Printf("Status codes %d | Count %d (%.2f%%)\n", key, count, percentage)
		return true
	})
}
