package main

import (
	"io"
	"log"
	"net"

	"github.com/Tarea1/Express/logistica"
	"google.golang.org/grpc"
)

func (s *ServerGRPC) Upload(stream messaging.GuploadService_UploadServer) (err error) {
	for {
		_, err = stream.Recv()
		if err != nil {
			if err == io.EOF {
				goto END
			}

			err = error.Wrapf(err, "failed unexpectadely while reading chunks from stream")
			return
		}
	}
	err = stream.SendAndClose(&messaging.UploadStatus{Message: "Upload received with success", Code: messaging.UploadStatusCode_Ok})
	return
}

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen on port 9000: %v", err)
	}

	s := logistica.Server{}

	grpcServer := grpc.NewServer()

	logistica.RegisterLogisticaServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port 9000: %v", err)
	}
}
