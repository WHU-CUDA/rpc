# WangxsRPC

#### Introduce
A High Performance RPC with Golang
#### 软件架构


#### Install

1.  go build main.go

#### Build the main

1. go run main.go


### How to add a ServiceMethod?
#### Step1 Start a registery service
```go
func startRegistry(wg *sync.WaitGroup) {
	l, _ := net.Listen("tcp", ":9999")
	registry.HandleHTTP()
	wg.Done()
	_ = http.Serve(l, nil)
}
```
The port :9999 stand for your service started at the port 9999
#### Step2 Start the rpc server
```go
func startServer(registryAddr string, wg *sync.WaitGroup) {
	var foo Foo
	l, _ := net.Listen("tcp", ":0")
	server := geerpc.NewServer()
	_ = server.Register(&foo)
	registry.Heartbeat(registryAddr, "tcp@"+l.Addr().String(), 0)
	wg.Done()
	server.Accept(l)
}
```
The registeryAddr is like this 
> registryAddr := "http://localhost:9999/_geerpc_/registry"

the struct is 
> http://ip:(your registery server port)/_geerpc_/registry
