# WangxsRPC

![](https://github.com/WHU-CUDA/rpc/workflows/Go/badge.svg)

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
> registryAddr := "http://localhost:9999/_wangxsrpc_/registry"

the struct is 
> http://ip:(your registery server port)/_wangxsrpc_/registry
#### Step4 Call the Method with different way
##### 1. Broadcast
```go
func broadcast(registry string) {
	d := xclient.NewWangxsRegistryDiscovery(registry, time.Second)
	xc := xclient.NewXClient(d, xclient.RandomSelect, nil)
	result := 0
	xc.Broadcast(context.Background(), "Foo.Sum", &Args{Num1: 1, Num2: 2}, &result)
	log.Printf("1 + 2 = %d\n", result)
}
```
> xclient is the high performance client for rpc

> xclient.RandomSelect is the load balance algorithm

##### 2.Call
```go
func call(registry string) {
	d := xclient.NewWangxsRegistryDiscovery(registry, 0)
	xc := xclient.NewXClient(d, xclient.RandomSelect, nil)
	defer func() { _ = xc.Close() }()
	// send request & receive response
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			result := 0
			xc.Call(context.Background(), "Foo.Sum", &Args{Num1: i, Num2: i * i}, &result)
		}(i)
	}
	wg.Wait()
}
```

#### ext Start Debug HTML Page
```go
wangxsrpc.HandleHTTP()
```

The page url is http://yourip:port/debug/wangxsrpc

default is http://localhost:9999/debug/wangxsrpc
