package main

import (
  "log"
 "fmt"
  "golang.org/x/net/context"
  pb "github.com/Tarea1/Express/logistica"
  "net"
  "google.golang.org/grpc"
  "time"
  "strconv"
  "os"
  "encoding/csv"
  "github.com/streadway/amqp"
  "encoding/json"
)

type Papi struct{ //struct que maneja la info general de la logistica (EL PAPI)
    numeroSeguimiento int
    colaNormal []Paquete
    colaPrioritaria []Paquete
    colaRetail []Paquete
    pedidos []Paquete
}

type Paquete struct {
  IDPaquete string
  seguimiento int
  tipo string
  valor int32
  intentos int
  estado string
  origen string
  destino string
  IDCamion int32
  Fecha_entrega string
}

func(s *Papi) SayHello(ctx context.Context, message *pb.Message) (*pb.Message,error){
  log.Printf("Received message body from client: %s", message.Body)
  return &pb.Message{Body: "Hello From the Server!"}, nil
}

func (s *Papi) SeguirPaquete(ctx context.Context, message *pb.Seguir) (*pb.Estado,error) {
  i := encontrarP(s.pedidos, int(message.GetSeguimiento()))
  if i != -1 {
    estadito := &pb.Estado{
        Estado:     s.pedidos[i].estado,
        Seguimiento: int32(s.pedidos[i].seguimiento),
        Intentos: int32(s.pedidos[i].intentos),
      }
    return estadito, nil
  }
  return &pb.Estado{}, nil
}

func (s *Papi) GenerarOrden(ctx context.Context, orden *pb.Orden) (*pb.NroSeguimiento, error) {
  //log.Printf("Received ID from client: %s", message.GetValor())
  s.numeroSeguimiento++
  fmt.Println("Se recibió la orden con ID: ", orden.GetIDPaquete())
  crearRegistro(s.numeroSeguimiento,orden)
  paquete, tipo := CrearPaquete(orden, s.numeroSeguimiento)
  if tipo == 0 { //normal
    s.colaNormal = append(s.colaNormal, paquete)

  } else if tipo == 1 { //prioritario
    s.colaPrioritaria = append(s.colaPrioritaria, paquete)
  } else {
    s.colaRetail = append(s.colaRetail, paquete)
  }
  s.pedidos = append(s.pedidos, paquete)
  //fmt.Println(s.colaNormal)
  //fmt.Println(s.colaPrioritaria)
  //fmt.Println(s.colaPyme)
  return &pb.NroSeguimiento{NumeroSeguimiento: int32(s.numeroSeguimiento)}, nil //nro_seguimiento
}

func CrearPaquete(orden *pb.Orden, numeroSeguimiento int) (Paquete, int32){
  var tipoO = orden.GetTipo()
  var tipo string
  //fmt.Println(tipoO)
  if tipoO == 0{
      tipo = "normal"
  }else if tipoO == 1{
      tipo = "prioritario"
  }else{
      tipo = "retail"
  }
  var paquete = Paquete {
    IDPaquete: orden.GetIDPaquete(),
    seguimiento: int(numeroSeguimiento),
    tipo: tipo,
    valor: orden.GetValor(),
    intentos: 0,
    estado: "En bodega",
    origen: orden.GetOrigen(),
    destino: orden.GetDestino(),
  }
  return paquete, tipoO
}

func crearRegistro(nro_seguimiento int, orden *pb.Orden) {
	f, err := os.OpenFile("./registrosLOG.csv", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Printf("Unable to open file")
	}
	defer f.Close()

	timestamp := time.Now()
	timestamp2 := timestamp.Format("2006-01-02 00:00")

	var data [][]string //https://golangcode.com/write-data-to-a-csv-file/
                        //https://www.golangprograms.com/sample-program-to-create-csv-and-write-data.html
  var tipoO = orden.GetTipo()
  var tipo string
  //fmt.Println(tipoO)
  if tipoO == 0{
      tipo = "normal"
      //fmt.Println("AAAaA")
  }else if tipoO == 1{
      tipo = "prioritario"
  }else{
      tipo = "retail"
  }

	data = append(data, []string{timestamp2,
    orden.GetIDPaquete(),
		tipo,
		orden.GetNombre(),
		strconv.Itoa(int(orden.GetValor())),
		orden.GetOrigen(),
		orden.GetDestino(),
    strconv.Itoa(nro_seguimiento)})

	writer := csv.NewWriter(f)
	writer.WriteAll(data)
}

func (s *Papi) UpdateEstado(ctx context.Context, estado_msg *pb.Estado) (*pb.Message, error) {
  i := encontrarP(s.pedidos, int(estado_msg.GetSeguimiento()))
  //fmt.Println("nr seguimiento:",estado_msg.GetSeguimiento())
  if i != -1 {
    s.pedidos[i].estado = estado_msg.GetEstado()
    s.pedidos[i].intentos = int(estado_msg.GetIntentos())

    // RABBIT --------
    conn, err := amqp.Dial("amqp://dist61:dist61@dist61:9000/")
	   failOnError(err, "Failed to connect to RabbitMQ")
	    defer conn.Close()

	     ch, err := conn.Channel()
	     failOnError(err, "Failed to open a channel")
       defer ch.Close()

	      q, err := ch.QueueDeclare(
		    "hello", // name
		     false,   // durable
		     false,   // delete when unused
		     false,   // exclusive
		     false,   // no-wait
		    nil,     // arguments
	       )
	failOnError(err, "Failed to declare a queue")

	finanzas := Paquete{
		IDPaquete: s.pedidos[i].IDPaquete,
		tipo:      s.pedidos[i].tipo,
		valor:     s.pedidos[i].valor,
		intentos:  s.pedidos[i].intentos,
		estado:    s.pedidos[i].estado}

	b, _ := json.Marshal(finanzas)

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(b),
		})
	failOnError(err, "Failed to publish a message")
    //---------------
    fmt.Println("Entrega finalizada para el paquete con número de entrega:",estado_msg.GetSeguimiento(),"con estado ",estado_msg.GetEstado())
    return &pb.Message{Body: "Orden Finalizada "}, nil
  }
  return &pb.Message {Body: "Ocurrió un error"}, nil
}
func failOnError(err error, msg string) {
  if err != nil {
    log.Fatalf("%s: %s", msg, err)
  }
}
func (s *Papi) GetPackage(ctx context.Context, camion *pb.Camion)(*pb.PaqueteMSG/*,int*/, error){

  if camion.GetTipo() == "normal" {
    var suma = len(s.colaPrioritaria) + len(s.colaNormal)
    if suma == 0{
      return &pb.PaqueteMSG{},nil
    }else if suma == 1{ // hay 1 entre los 2, se debe chequear cual
        if len(s.colaPrioritaria) == 1{
          i := encontrarP(s.pedidos, s.colaPrioritaria[0].seguimiento)
          pCamion, temp := crearPmsje (s.colaPrioritaria,0)
          s.colaPrioritaria = temp

          if i != -1 {
            s.pedidos[i].estado = "En camino"
            s.pedidos[i].IDCamion = int32(camion.GetIDCamion())
            return pCamion,nil
          }
        }else{
          i := encontrarP(s.pedidos, s.colaNormal[0].seguimiento)
          pCamion, temp := crearPmsje (s.colaNormal,0)
          s.colaNormal = temp

          if i != -1 {
            s.pedidos[i].estado = "En camino"
            s.pedidos[i].IDCamion = int32(camion.GetIDCamion())
            return pCamion,nil
          }
        }
    }else{
      if len(s.colaPrioritaria) > 0{
        i := encontrarP(s.pedidos, s.colaPrioritaria[0].seguimiento)
        pCamion, temp := crearPmsje (s.colaPrioritaria,1)
        s.colaPrioritaria = temp

        if i != -1 {
          s.pedidos[i].estado = "En camino"
          s.pedidos[i].IDCamion = int32(camion.GetIDCamion())
          return pCamion,nil
        }
      }else{
        i := encontrarP(s.pedidos, s.colaNormal[0].seguimiento)
        pCamion, temp := crearPmsje (s.colaNormal,1)
        s.colaNormal = temp

        if i != -1 {

          s.pedidos[i].estado = "En camino"
          s.pedidos[i].IDCamion = int32(camion.GetIDCamion())
          return pCamion,nil
        }
      }
    }
  }else{// caso retail
    var suma = len(s.colaPrioritaria) + len(s.colaRetail)
    if suma == 0{
      return &pb.PaqueteMSG{},nil
    }else if suma == 1{ // hay 1 entre los 2, se debe chequear cual
        if len(s.colaRetail) == 1{
          i := encontrarP(s.pedidos, s.colaRetail[0].seguimiento)
          pCamion, temp := crearPmsje (s.colaRetail,0)
          s.colaRetail = temp

          if i != -1 {
            s.pedidos[i].estado = "En camino"
            s.pedidos[i].IDCamion = int32(camion.GetIDCamion())



            return pCamion,nil
          }
        }else{
          i := encontrarP(s.pedidos, s.colaPrioritaria[0].seguimiento)
          pCamion, temp := crearPmsje (s.colaPrioritaria,0)
          s.colaPrioritaria = temp

          if i != -1 {
            s.pedidos[i].estado = "En camino"
            s.pedidos[i].IDCamion = int32(camion.GetIDCamion())



            return pCamion,nil
          }
        }
    }else{
      if len(s.colaRetail) > 0{
        i := encontrarP(s.pedidos, s.colaRetail[0].seguimiento)
        pCamion, temp := crearPmsje (s.colaRetail,1)
        s.colaRetail = temp

        if i != -1 {
          s.pedidos[i].estado = "En camino"
          s.pedidos[i].IDCamion = int32(camion.GetIDCamion())



          return pCamion,nil
        }
      }else{
        i := encontrarP(s.pedidos, s.colaPrioritaria[0].seguimiento)
        pCamion, temp := crearPmsje (s.colaPrioritaria,1)
        s.colaPrioritaria = temp

        if i != -1 {

          s.pedidos[i].estado = "En camino"
          s.pedidos[i].IDCamion = int32(camion.GetIDCamion())



          return pCamion,nil
        }
      }
    }
  }
  return &pb.PaqueteMSG{},nil
}

func crearPmsje (s []Paquete, flag int32)(*pb.PaqueteMSG,[]Paquete){
  if len(s) > 1 {
    pack := s[0]
    temp  := s[1:]
    pCamion := &pb.PaqueteMSG{
        IDPaquete:   pack.IDPaquete,
        Tipo:        pack.tipo,
        Origen:      pack.origen,
        Destino:     pack.destino,
        Valor:       int32(pack.valor),
        Flagxd:      flag,
        Seguimiento: int32(pack.seguimiento),
      }
      log.Printf("Estado actualizado del paquete con numero de seguimiento %v a En camino", pack.seguimiento)

      return pCamion,temp
    }else if len(s) == 1{
      pack := s[0]
      pCamion := &pb.PaqueteMSG{
          IDPaquete:          pack.IDPaquete,
          Tipo:        pack.tipo,
          Origen:      pack.origen,
          Destino:     pack.destino,
          Valor:       int32(pack.valor),
          Flagxd:      flag,
          Seguimiento: int32(pack.seguimiento),
        }
        log.Printf("Estado actualizado del paquete con numero de seguimiento %v a En camino", pack.seguimiento)

        return pCamion,[]Paquete{}
    }
    return &pb.PaqueteMSG{},[]Paquete{}
}

func encontrarP(pedidos []Paquete, nro_seg int) (int) {
  var i = 0
  for i < len(pedidos){
    if pedidos[i].seguimiento == nro_seg {
      return i
    }else{
      i++
    }
  }
  return -1
}

func main(){
  lis,err := net.Listen("tcp",":9000")
  if err!= nil {
    log.Fatalf("Failed to listen on port 9000: %v", err)
  }

  s := Papi{}
  s.numeroSeguimiento = 0

  grpcServer:= grpc.NewServer()

  pb.RegisterLogisticaServiceServer(grpcServer, &s)

  if err := grpcServer.Serve(lis); err!=nil{
    log.Fatalf("Failed to serve gRPC server over port 9000: %v", err)
  }
}
