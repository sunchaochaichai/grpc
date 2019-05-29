package main

import (
	"golang.org/x/net/context"
	pb "grpc-example/proto/hello"
	"google.golang.org/grpc"
	"log"
)

const (
	Address = "127.0.0.1:50052"
	OpenTLS = false
)
var opts []grpc.DialOption

// customCredential 自定义认证
type customCredential struct{}

func (c customCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"appid":  "101010",
		"appkey": "i am key",
	}, nil
}

func (c customCredential) RequireTransportSecurity() bool {
	if OpenTLS {
		return true
	}

	return false
}
var CustomCredential = customCredential{}
func main() {
	//连接
	opts = append(opts,grpc.WithInsecure())
	opts = append(opts,grpc.WithPerRPCCredentials(CustomCredential))
	conn,err := grpc.Dial(Address,opts...)
	if err != nil {
		log.Fatalf("连接错误%v",err)
	}
	//关闭
	defer conn.Close()
	//初始化客户端
	c := pb.NewHelloClient(conn)
	//调用方法
	reqBody := new(pb.HelloRequest)
	reqBody.Name = "gRPC"
	r,err := c.SayHello(context.Background(),reqBody)
	if err != nil{
		log.Fatalf("调用错误%v",err)
	}
	log.Fatalf("显示回调信息%v",r.Message)
}

