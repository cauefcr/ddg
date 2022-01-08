package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"websearch"
	"websearch/provider"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ddg [query]")
		return
	}
	web := websearch.New(provider.NewUnofficialDuckDuckGo())
	csig := make(chan os.Signal, 1)
	signal.Notify(csig, syscall.SIGPIPE)
	go func() {
		<-csig
		os.Exit(0)
	}()
	backoff := 1 * time.Second
	for {
		res, err := web.Search(strings.Join(os.Args[1:], " "), 25)
		if err != nil {
			backoff *= 2
			time.Sleep(backoff)
			continue
		}
		for _, r := range res {
			fmt.Printf("%v|%v|%v\n", r.Title, r.Link.String(), r.Description)
		}
		time.Sleep(backoff)
	}
}
