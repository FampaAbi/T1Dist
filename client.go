package main

import (
  "log"

  "golang.org/x/net/context"
  "google.golang.org/grpc"

  "github.com/Tarea1/Express/logistica"
)

func main(){
  var conn *grpc.ClientConn
  conn, err := grpc.Dial("dist61:9000", grpc.WithInsecure())
  if err != nil{
    log.Fatalf("could not connect: %s", err)
  }
  defer conn.Close()

  c := logistica.NewLogisticaServiceClient(conn)

  message := logistica.Message{
    Body: "Hello from the client!",
  }

  response, err := c.SayHello(context.Background(),&message)
  if err!= nil{
    log.Fatalf("Error when calling SayHello: %s", err)
  }

  log.Printf("Response from Server: %s", response.Body)
}
