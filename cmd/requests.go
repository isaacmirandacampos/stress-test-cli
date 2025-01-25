package cmd

import (
	`fmt`
	`net/http`
	`sync/atomic`
)

func request(url string, statusCodeSuccess *int32) {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	if res.StatusCode == 200 {
		atomic.AddInt32(statusCodeSuccess, 1)
	}
	fmt.Printf("Status code %v \n", res.StatusCode)
}
