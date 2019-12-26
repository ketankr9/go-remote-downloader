package main

import(
	"fmt"
	"net/rpc"
	"sync"
	"net/http"
	"strconv"
)

type Item struct {
	Url string
	Idx	int
	Start	int
	End	int
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
	client, err := rpc.DialHTTP("tcp", "localhost:6970")
	
	var reply int

	// client.Call("API.SetSem", 16, &reply)

	if err != nil {
		fmt.Println("Unable to connect ", err)
	}

	content_length := getContentLength(url)
	
	num_splits := 16
	chunks := content_length/num_splits
	
	// var chunks int
	// chunks = 10000000 // 1e6 bytes == 1 megabytes

	for i:=0; i<content_length/chunks; i++ {
		start := i*chunks
		end := (i+1)*chunks

		if i + 1 == content_length/chunks{
			end += content_length%chunks
		}

		item := Item{url, i, start, end-1, filename}
		
		wg.Add(1)
		go func(){
			client.Call("API.DownloadItem", item, &reply) // blocking without prefix `go` (waits for return value)
			wg.Done()
		}()

	}

	wOnDownload.Done()
	return nil
}  

func main(){

	for i := 0; i<1; i++ {
		wOnDownload.Add(1)
		go download("https://speed.hetzner.de/100MB.bin", "100MB.bin") // URL, FILENAME
	}

	wOnDownload.Wait()
	wg.Wait()

	fmt.Println("Downloaded")
}