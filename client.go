package main

import (
	"log"
	"os"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/Tarea1/Express/logistica"
)

func (c *ClienteGRPC) UploadFile(ctx context.Context, f string) (stats Stats, err error) {
	file, err = os.Open(f)
	stream, err := c.client.Upload(ctx)
	stats.StartedAt = time.Now()
	buf = make([]byte, c.chunkSize)
	for writing {
		n, err = file.Read(buf)
		stream.Send(&messaging.Chunk{Content: buf[:n]})
	}
	stats.FinishedAt = time.Now()
	status, err = stream.CloseAndRecv()
}

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("dist61:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %s", err)
	}
	defer conn.Close()

	c := logistica.NewLogisticaServiceClient(conn)

	message := logistica.Message{
		Body: "Hello from the client!",
	}

	response, err := c.SayHello(context.Background(), &message)
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}

	log.Printf("Response from Server: %s", response.Body)
}
