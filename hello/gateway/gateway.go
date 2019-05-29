package main

import (
	"flag"
	"net/http"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	gw "grpc-example/proto/hello"
)

var (
	echoEndpoint = flag.String("echo_endpoint", "127.0.0.1:50052", "endpoint of YourService")
)

const(
	OpenTLS = false
)

// customCredential 自定义认证
//type customCredential struct{}
//
//func (c customCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
//	return map[string]string{
//		"appid":  "1010101",
//		"appkey": "i am key",
//	}, nil
//}
//func (c customCredential) RequireTransportSecurity() bool {
//	if OpenTLS {
//		return true
//	}
//
//	return false
//}
//
//var CustomCredential = customCredential{}
var opts []grpc.DialOption

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts = append(opts,grpc.WithInsecure())
	//opts = append(opts,grpc.WithPerRPCCredentials(CustomCredential))
	err := gw.RegisterHelloHandlerFromEndpoint(ctx, mux, *echoEndpoint, opts)
	if err != nil {
		return err
	}
	return http.ListenAndServe(":8080", mux)
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}