package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func main() {
	cs := flag.String("c", "1", "Number of threads to use")
	url := flag.String("url", "http://localhost:8000/api/products", "Url to do gets")
	ns := flag.String("nc", "1", "Transactions per thread")
	flag.Parse()
	c, err := strconv.Atoi(*cs)
	if err != nil {
		log.Fatal(err)
	}
	n, err := strconv.Atoi(*ns)
	if err != nil {
		log.Fatal(err)
	}
	wg := sync.WaitGroup{}
	client := &http.Client{
		Transport: &http.Transport{
			MaxConnsPerHost:     c,
			MaxIdleConns:        c,
			MaxIdleConnsPerHost: c,
		},
	}
	t0 := time.Now()
	for i := 0; i < c; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < n; j++ {
				req, err := http.NewRequest(
					http.MethodGet,
					*url,
					nil,
				)
				if err != nil {
					log.Fatal(err)
				}

				resp, err := client.Do(req)
				if err != nil {
					log.Println(err)
					continue
				}
				defer resp.Body.Close()

				_, err = io.ReadAll(resp.Body)
				if err != nil {
					log.Fatalf("Error reading response body: %v", err)
				}
			}
		}()
	}
	wg.Wait()
	tf := time.Since(t0)
	tps := float64(c) * float64(n) / tf.Seconds()
	fmt.Println("SEC :", tf.Seconds())
	fmt.Println("TPS:", tps)
}
