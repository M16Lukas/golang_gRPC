package main

import (
	"bytes"
	"context"
	"fmt"
	"gRPC/config"
	"gRPC/pb"
	"io"
	"log"
	"net"
	"os"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedFileServiceServer
}

/*
-----------------------------------------------------------
* Unary RPC
-----------------------------------------------------------
*/
func (*server) ListFiles(ctx context.Context, req *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
	fmt.Println("ListFiles was invoked")

	dir := config.Config.StoragePath
	paths, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var filenames []string
	for _, path := range paths {
		if !path.IsDir() {
			filenames = append(filenames, path.Name())
		}
	}

	res := &pb.ListFilesResponse{
		Filenames: filenames,
	}

	return res, nil
}

/*
-----------------------------------------------------------
* Server Streaming RPC
-----------------------------------------------------------
*/
func (*server) Download(req *pb.DownloadRequest, stream pb.FileService_DownloadServer) error {
	fmt.Println("Download was invoked")

	filename := req.GetFilename()
	path := config.Config.StoragePath + filename

	// エラーハンドリング
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return status.Error(codes.NotFound, "file was not found")
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := make([]byte, 5)
	for {
		n, err := file.Read(buf)
		if n == 0 || err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		res := &pb.DownloadResponse{
			Data: buf[:n],
		}

		sendErr := stream.Send(res)

		if sendErr != nil {
			return sendErr
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}

/*
-----------------------------------------------------------
* Client Streaming RPC
-----------------------------------------------------------
*/
func (*server) Upload(stream pb.FileService_UploadServer) error {
	fmt.Println("Upload was invoked")

	var buf bytes.Buffer
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			res := &pb.UploadResponse{Size: int32(buf.Len())}
			return stream.SendAndClose(res)
		}

		if err != nil {
			return err
		}

		data := req.GetData()
		log.Printf("received data(bytes): %v", data)
		log.Printf("received data(string): %v", string(data))
		buf.Write(data)
	}
}

/*
-----------------------------------------------------------
* Bidirectional Streaming RPC
-----------------------------------------------------------
*/
func (*server) UploadAndNotifyProgess(stream pb.FileService_UploadAndNotifyProgressServer) error {
	fmt.Println("UploadAndNotifyProgess was invoked")

	size := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		data := req.GetData()
		log.Printf("received data: %v", data)
		size += len(data)

		res := &pb.UploadAndNotifyProgressResponse{
			Msg: fmt.Sprintf("received %v bytes", size),
		}

		err = stream.Send(res)
		if err != nil {
			return err
		}
	}
}

/*
-----------------------------------------------------------
* Interceptor - ロギング
-----------------------------------------------------------
*/
func myLogging() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {

		log.Printf("request data : %+v", req)

		resp, err = handler(ctx, req)
		if err != nil {
			return nil, err
		}

		log.Printf("response data : %+v", resp)

		return resp, nil
	}
}

/*
-----------------------------------------------------------
* Interceptor - 認証
-----------------------------------------------------------
*/
func authorize(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		return nil, err
	}

	if token != "test-token" {
		return nil, status.Error(codes.Unauthenticated, "token is invalid")
	}

	return ctx, nil
}

func main() {
	lis, err := net.Listen("tcp", "localhost:"+config.Config.Port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// SSL
	creds, err := credentials.NewServerTLSFromFile(
		config.Config.SSLPath+"localhost.pem",
		config.Config.SSLPath+"localhost-key.pem",
	)

	if err != nil {
		log.Fatalln(err)
	}

	// Interceptorの実装
	s := grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				myLogging(),
				grpc_auth.UnaryServerInterceptor(authorize),
			),
		),
	)
	pb.RegisterFileServiceServer(s, &server{})

	fmt.Println("server is runnings...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
