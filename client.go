package main

import (
	"fmt"
	"net/http"
	"net/rpc"
	"strconv"
	"sync"
)

type Item struct {
	Url      string
	Idx      int
	Start    int
	End      int
	Filename string
}

var wg sync.WaitGroup
var wOnDownload sync.WaitGroup

func getContentLength(url string) int {
	res, _ := http.Head(url)
	length, _ := strconv.Atoi(res.Header["Content-Length"][0]) // Get the content length from the header request
	return length
}

func download(url string, filename string) error {
	var sem1, sem2 chan int

	client1, err := rpc.DialHTTP("tcp", "10.2.60.42:6969")
	if err == nil {
		sem1 = make(chan int, 16) // maximum parallel tasks assigned to localhost:6970
	} else {
		// failed to connect to the server
		fmt.Println("Failed to connect to ..42")
		sem1 = make(chan int, 1)
		sem1 <- 1
	}

	client2, err := rpc.DialHTTP("tcp", "localhost:6969")
	if err == nil {
		sem2 = make(chan int, 16) // maximum parallel tasks assigned to localhost:6969
	} else {
		// failed to connect to the server
		sem2 = make(chan int, 1)
		sem2 <- 1
	}

	var reply int
	// client.Call("API.SetSem", 16, &reply)

	content_length := getContentLength(url)

	num_splits := 32
	chunks := content_length / num_splits

	// var chunks int
	// chunks = 10000000 // 1e6 bytes == 1 megabytes

	for i := 0; i < content_length/chunks; i++ {
		start := i * chunks
		end := (i + 1) * chunks

		if i+1 == content_length/chunks {
			end += content_length % chunks
		}

		item := Item{url, i, start, end - 1, filename}

		select {
		case sem1 <- 1:
			wg.Add(1)
			go func() {
				client1.Call("API.DownloadItem", item, &reply) // blocking without prefix `go` (waits for return value)
				<-sem1
				wg.Done()
			}()
		case sem2 <- 1:
			wg.Add(1)
			go func() {
				client2.Call("API.DownloadItem", item, &reply) // blocking without prefix `go` (waits for return value)
				<-sem2
				wg.Done()
			}()
		}
	}

	wOnDownload.Done()
	return nil
}

func main() {
	fmt.Println("Starting...")
	for i := 0; i < 1; i++ {
		wOnDownload.Add(1)
		go download("https://speed.hetzner.de/100MB.bin", "100MB.bin") // URL, FILENAME
	}

	wOnDownload.Wait()
	wg.Wait()

	fmt.Println("Downloaded")
}
