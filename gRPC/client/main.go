package main

import (
	"context"
	"fmt"
	"gRPC/config"
	"gRPC/pb"
	"io"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func main() {
	// SSL
	certFile := config.Config.RootCAPath + "rootCA.pem"
	creds, err := credentials.NewClientTLSFromFile(certFile, "")

	// サーバーとの接続
	conn, err := grpc.Dial("localhost:"+config.Config.Port, grpc.WithTransportCredentials(creds))

	if err != nil {
		log.Fatalf("Failed to connet: %v", err)
	}
	defer conn.Close()

	client := pb.NewFileServiceClient(conn)
	/* -----------------------------------------------------------
	* Unary RPC
	-----------------------------------------------------------*/
	// callListFiles(client)

	/* -----------------------------------------------------------
	* Server Streaming RPC
	-----------------------------------------------------------*/
	callDownload(client)

	/* -----------------------------------------------------------
	* Client Streaming RPC
	-----------------------------------------------------------*/
	// CallUpload(client)

	/* -----------------------------------------------------------
	* Bidirectional Streaming RPC
	-----------------------------------------------------------*/
	// CallUploadAndNotifyProgess(client)
}

/*
-----------------------------------------------------------
* Unary RPC
-----------------------------------------------------------
*/
func callListFiles(client pb.FileServiceClient) {
	// contextに定義したいメタデータを定義
	md := metadata.New(map[string]string{
		"authorization": "Bearer bad-token",
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	res, err := client.ListFiles(ctx, &pb.ListFilesRequest{})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(res.GetFilenames())
}

/*
-----------------------------------------------------------
* Server Streaming RPC
-----------------------------------------------------------
*/
func callDownload(client pb.FileServiceClient) {
	// deadlines(= timeout, 5sec)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	req := &pb.DownloadRequest{Filename: "name.txt"}
	stream, err := client.Download(ctx, req)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			resErr, ok := status.FromError(err)
			if ok {
				if resErr.Code() == codes.NotFound {
					log.Fatalf("Error code : %v, Error message : %v", resErr.Code(), resErr.Message())
				} else if resErr.Code() == codes.DeadlineExceeded {
					log.Fatalln("Deadline Exceeded")
				} else {
					log.Fatalln("unknown grpc error")
				}
			} else {
				log.Fatalln(err)
			}
		}

		log.Printf("Response from Download(bytes): %v", res.GetData())
		log.Printf("Response from Download(string): %v", string(res.GetData()))
	}
}

/*
-----------------------------------------------------------
* Client Streaming RPC
-----------------------------------------------------------
*/
func CallUpload(client pb.FileServiceClient) {
	filename := "sports.txt"
	path := config.Config.StoragePath + filename

	file, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	stream, err := client.Upload(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	buf := make([]byte, 5)
	for {
		n, err := file.Read(buf)
		if n == 0 || err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalln(err)
		}

		req := &pb.UploadRequest{
			Data: buf[:n],
		}
		sendErr := stream.Send(req)
		if sendErr != nil {
			log.Fatalln(sendErr)
		}

		time.Sleep(1 * time.Second)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalln(err)
	}

	log.Panicf("received data size: %v", res.GetSize())
}

/*
-----------------------------------------------------------
* Bidirectional Streaming RPC
-----------------------------------------------------------
*/
func CallUploadAndNotifyProgess(client pb.FileServiceClient) {
	filename := "sports.txt"
	path := config.Config.StoragePath + filename

	file, err := os.Open(path)
	if err != nil {
		log.Fatalln("1.", err)
	}
	defer file.Close()

	stream, err := client.UploadAndNotifyProgress(context.Background())
	if err != nil {
		log.Fatalln("2.", err)
	}

	// request
	buf := make([]byte, 5)
	go func() {
		for {
			n, err := file.Read(buf)
			if n == 0 || err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalln("3-1.", err)
			}

			req := &pb.UploadAndNotifyProgressRequest{
				Data: buf[:n],
			}
			sendErr := stream.Send(req)

			if sendErr != nil {
				log.Fatalln("3-2.", sendErr)
			}
			time.Sleep(1 * time.Second)
		}

		err := stream.CloseSend()
		if err != nil {
			log.Fatalln("3-3.", err)
		}
	}()

	// response
	ch := make(chan struct{})
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalln("4-1.", err)
			}

			log.Printf("received message: %v", res.GetMsg())
		}
		close(ch)
	}()
	<-ch
}
