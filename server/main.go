package main

import(
	"fmt"
	"net/rpc"
	"net"
	"net/http"
	"strconv"
	"io/ioutil"
	"log"
)

type Item struct {
	Url string
	Idx	int
	Start	int
	End	int
	Filename string
}

type API int
var semaphore chan int

func download(item Item) error {

	client := &http.Client {}
	req, _ := http.NewRequest("GET", item.Url, nil)  
	range_header := "bytes=" + strconv.Itoa(item.Start) +"-" + strconv.Itoa(item.End) // Add the data for the Range header of the form "bytes=0-100"
	req.Header.Add("Range", range_header)
	resp,_ := client.Do(req)
	defer resp.Body.Close()
	reader, _ := ioutil.ReadAll(resp.Body)
	body := string(reader)
	ioutil.WriteFile("/tmp/" + strconv.Itoa(item.Idx) +"_"+ item.Filename +".part", []byte(string(body)), 0x777) // Write to the file i as a byte array
	fmt.Printf(" %d", item.Idx)

	<-semaphore
	return nil
}

func (a *API)DownloadItem(item Item, reply *int) error {
	semaphore <- 1
	download(item) // go download(item) ---> if you want confirmation has the file been downloaded
	*reply = 200
	return nil
}

func (a* API)SetSem(parallel_connections int, reply *int) error {
	semaphore = make(chan int, parallel_connections)

	*reply = 201
	return nil
}

func main(){

	// Default
	semaphore = make(chan int, 16)
	
	var api = new(API)

	if  rpc.Register(api) != nil {
		log.Fatal("Error register rpc")
	}

	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", ":6970")
	
	if err != nil{
		log.Fatal("Error maybe port unavailable", err)
	}

	fmt.Printf("Serving on port %d\n", 6970)
	if  http.Serve(listener, nil) != nil {
		log.Fatal("Error http.server")
	}

}