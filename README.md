# go-remote-downloader  

>A google RPC golang server and client program which downloads a file in parts parallelly over a remote server.  

**Client code**  
The client chooses the no of parts in which the file is to be downloaded.  Then it sends each chunk to the server over rpc to donwnload.  

**Server code**  
The server has a fixed no of download routines which can run in parallel, default to 16 here.  
The server accepts the dowload request from the client which comparises mainly of *url* and *byte-range* of the chunk and saves it in ```/tmp/```. Which can further be combined using cat. Refer ```verify.sh``` for more 

**On server side**
```
$ git clone https://github.com/ketankr9/go-remote-downloader.git
$ cd go-remote-downloader
$ go run server/main.go
```

**On Client Side**  
```
$ git clone https://github.com/ketankr9/go-remote-downloader.git
$ cd go-remote-downloader
$ go run client/main.go
```

**Verify reconstruction on server side**  
```
$ wget -o /tmp/100MB.bin https://speed.hetzner.de/100MB.bin
$ cd go-remote-downloader
$ chmod +x verify.sh
$ ./verify.sh
```

**Further Goals**  
*   Modify client code such that it distributes the file to be downloaded in chunks over multiple servers/clusters.
*   Handle various errors. Like http requests 