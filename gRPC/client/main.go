package main

import (
	"context"
	"fmt"
	"gRPC/pb"
	"io"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
)

func main() {

	// サーバーとの接続
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

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
	// callDownload(client)

	/* -----------------------------------------------------------
	* Client Streaming RPC
	-----------------------------------------------------------*/
	// CallUpload(client)

	/* -----------------------------------------------------------
	* Bidirectional Streaming RPC
	-----------------------------------------------------------*/
	CallUploadAndNotifyProgess(client)
}

/*
-----------------------------------------------------------
* Unary RPC
-----------------------------------------------------------
*/
func callListFiles(client pb.FileServiceClient) {
	res, err := client.ListFiles(context.Background(), &pb.ListFilesRequest{})
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
	req := &pb.DownloadRequest{Filename: "name.txt"}
	stream, err := client.Download(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalln(err)
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
	path := "C:\\Users\\MH\\go\\src\\golang_gRPC\\gRPC\\storage\\" + filename

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
	path := "C:\\Users\\MH\\go\\src\\golang_gRPC\\gRPC\\storage\\" + filename

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
