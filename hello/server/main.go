package main

import ("context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	pb "grpc-example/proto/hello"
	"log"
	"net"
)


const (
	Address = "127.0.0.1:50052"
)
type helloServer struct {

}

func(this helloServer) SayHello(ctx context.Context,in *pb.HelloRequest) (*pb.HelloResponse,error) {
	resp := new(pb.HelloResponse)
	resp.Message = fmt.Sprintf("hello %s",in.Name)
	return resp,nil
}
var HelloServer = helloServer{}
func auth(ctx context.Context) error {
	md,ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return grpc.Errorf(codes.Unauthenticated,"无token认证信息")
	}
	var(
		appId string
		appKey string
	)

	if val, ok := md["appid"]; ok {
		appId = val[0]
	}

	if val, ok := md["appkey"]; ok {
		appKey = val[0]
	}

	if appId != "101010" || appKey != "i am key" {
		return grpc.Errorf(codes.Unauthenticated, "Token认证信息无效: appid=%s, appkey=%s", appId, appKey)
	}

	return nil
}

func main() {
	listen,err := net.Listen("tcp",Address)
	if err != nil {
		log.Fatalf("监听错误%v",err)
	}
	var opts []grpc.ServerOption
	// 注册interceptor
	var interceptor grpc.UnaryServerInterceptor
	interceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		err = auth(ctx)
		if err != nil {
			return
		}
		// 继续处理请求
		return handler(ctx, req)
	}
	opts = append(opts,grpc.UnaryInterceptor(interceptor))
	//实例化grpc
	s := grpc.NewServer(opts...)
	pb.RegisterHelloServer(s,HelloServer)
	s.Serve(listen)
}

