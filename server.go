package main

import(
  "log"
  "net"

  "google.golang.org/grpc"
  sv "github.com/Tarea1/Express/logistica"
)

func main(){
  lis,err := net.Listen("tcp",":9000")
  if err!= nil {
    log.Fatalf("Failed to listen on port 9000: %v", err)
  }

  s := sv.Server{}
  s.numeroSeguimiento = 0

  grpcServer:= grpc.NewServer()

  sv.RegisterLogisticaServiceServer(grpcServer, &s)

  if err := grpcServer.Serve(lis); err!=nil{
    log.Fatalf("Failed to serve gRPC server over port 9000: %v", err)
  }
}
