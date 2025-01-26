package cmd

import (
	`sync`
)

func summaryStatusCode(statusCodes *sync.Map, key int) {
	existentCount, ok := statusCodes.Load(key)
	var statusCount int32
	if ok == true {
		statusCount = existentCount.(int32) + 1
	} else {
		statusCount = 1
	}
	statusCodes.Store(key, statusCount)
}
