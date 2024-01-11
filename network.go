package main

import (
	"net/http"
	"time"
)

const (
	TestUrl = "https://drive.proton.me"
)

func CheckNetwork() bool {
	client := http.Client{
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	_, err := client.Head(TestUrl)
	return err == nil
}

func WaitNetwork() {
	for !CheckNetwork() {
		time.Sleep(time.Second * 10)
	}
}
