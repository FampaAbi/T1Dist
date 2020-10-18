package main

import (
	//"bufio"
	"encoding/csv"
	"fmt"
  "time"
	"log"
	"os"

  "golang.org/x/net/context"
  "google.golang.org/grpc"
  pb "github.com/Tarea1/Express/logistica"
  "strconv"
  "math/rand"
)

type camiones struct{
  tipoC string
  IDcamion int
  paqueteA PaqueteCamion
  paqueteB PaqueteCamion
  total int
}

type PaqueteCamion struct{
  IDPaquete string
  Seguimiento int
  Tipo string
  Valor int
  Origen string
  Destino string
  Intentos int
  //Estado string
  Fecha_entrega string
}



func traerPaquetes(camion *camiones, tEspera int, conn *grpc.ClientConn) {
  c := pb.NewLogisticaServiceClient(conn)
	p, _ := c.GetPackage(context.Background(), &pb.Camion{
		IDCamion:     int32(camion.IDcamion),
		Tipo:          camion.tipoC,})

  if p.GetFlagxd() == 3 {
		//Paquete_Camion := PaqueteCamion{}
    fmt.Println("xd")
	}else{


    if p.GetIDPaquete() != "" {
      fmt.Println("La ID del paquete recibido por el camion ",camion.IDcamion," es:", p.GetIDPaquete())
      //crearRegistroC(Paquete_Camion,camion.IDcamion)
      Paquete_Camion := PaqueteCamion{
    		IDPaquete:    p.GetIDPaquete(),
    		Tipo:         p.GetTipo(),
    		Valor:        int(p.GetValor()),
    		Origen:       p.GetOrigen(),
    		Destino:      p.GetDestino(),
    		Seguimiento:  int(p.GetSeguimiento()),
    		Intentos:     0,
    		Fecha_entrega: "",
      }
    	camion.total++
      camion.paqueteA = Paquete_Camion
    }

    if p.GetFlagxd() == 1{ // hay otro paquete para este camion
      p2, _ := c.GetPackage(context.Background(), &pb.Camion{
    		IDCamion:     int32(camion.IDcamion),
    		Tipo:          camion.tipoC,})


      if p2.GetIDPaquete() != ""  {
        fmt.Println("La ID del paquete recibido por el camion ",camion.IDcamion," es:", p2.GetIDPaquete())
        Paquete_Camion2 := PaqueteCamion{
          IDPaquete:    p2.GetIDPaquete(),
          Tipo:         p2.GetTipo(),
          Valor:        int(p2.GetValor()),
          Origen:       p2.GetOrigen(),
          Destino:      p2.GetDestino(),
          Seguimiento:  int(p2.GetSeguimiento()),
          Intentos:     0,
          Fecha_entrega: "",

        }
        camion.paqueteB = Paquete_Camion2
        camion.total++
        //crearRegistroC(Paquete_Camion2,camion.IDcamion)
      }

    }else{
      time.Sleep(time.Duration(tEspera)*time.Second)
      p2, _ := c.GetPackage(context.Background(), &pb.Camion{
    		IDCamion:            int32(camion.IDcamion),
    		Tipo:          camion.tipoC,})
      if p2.GetIDPaquete() != "" {
        fmt.Println("La ID del paquete recibido por el camion ",camion.IDcamion," es:", p2.GetIDPaquete())
        Paquete_Camion2 := PaqueteCamion{
      		IDPaquete:    p2.GetIDPaquete(),
      		Tipo:         p2.GetTipo(),
      		Valor:        int(p2.GetValor()),
      		Origen:       p2.GetOrigen(),
      		Destino:      p2.GetDestino(),
      		Seguimiento:  int(p2.GetSeguimiento()),
      		Intentos:     0,
      		Fecha_entrega: "",
        }
        camion.paqueteB = Paquete_Camion2
        camion.total++
        /*if p2.GetIDPaquete() != "" {
          fmt.Println("La ID del paquete recibido por el camion ",camion.IDcamion," es:", p2.GetIDPaquete())
          //crearRegistroC(Paquete_Camion2,camion.IDcamion)
    	 }*/
     }
   }


  }

}

func crearRegistroC( pCam PaqueteCamion, idC int) {
  f, err := os.OpenFile("./registrosCAM"+strconv.Itoa(idC)+".csv", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Printf("Unable to open file")
	}
	defer f.Close()

	timestamp := time.Now()
	timestamp2 := timestamp.Format("2006-01-02 00:00")

	var data [][]string //https://golangcode.com/write-data-to-a-csv-file/
                        //https://www.golangprograms.com/sample-program-to-create-csv-and-write-data.html


	data = append(data, []string{
    pCam.IDPaquete,
		pCam.Tipo,
		strconv.Itoa(pCam.Valor),
		pCam.Origen,
		pCam.Destino,
    strconv.Itoa(pCam.Intentos),
    timestamp2,})

	writer := csv.NewWriter(f)
	writer.WriteAll(data)
}

func prob80() int {
  rand.Seed(time.Now().UnixNano())
	prob := rand.Intn(100)
  
	if prob < 81 {
		return 1
	} else {
		return 0
	}
}

func updateEstado(estado string, nr_seg int, intentos int, conn *grpc.ClientConn) {
  c := pb.NewLogisticaServiceClient(conn)
	msg, _ := c.UpdateEstado(context.Background(), &pb.Estado{
		Estado: estado,
		Seguimiento: int32(nr_seg),
    Intentos: int32(intentos),
  })
  fmt.Println(msg)
}

func despacharPaquetes(camion *camiones, tEntrega int, conn *grpc.ClientConn){
  var flag int
  if (camion.paqueteA).Valor > (camion.paqueteB).Valor { // decidir prioridad entrega paquete
    flag = 1
  }else{
    flag = 0
  }
  var tipoPA = (camion.paqueteA).Tipo
  var tipoPB = (camion.paqueteB).Tipo
  //fmt.Println("a:",camion.paqueteA)
  //fmt.Println("b:",camion.paqueteB)
  //fmt.Println("CAMION:",camion)
  for camion.total > 0{
    //fmt.Println("Camion.total:",camion.total)
    entregado := prob80()

    time.Sleep(time.Duration(tEntrega)*time.Second)
    if flag == 1{ //se entrega A primero
      (camion.paqueteA).Intentos++
      if entregado == 1{ //se entrega
        crearRegistroC(camion.paqueteA,camion.IDcamion)
            //ACTUALIZAR ESTADO ENTREGADO
        //fmt.Println("contenido paquete A:",camion.paqueteA)
        updateEstado("Recibido" , (camion.paqueteA).Seguimiento, (camion.paqueteA).Intentos , conn)
        camion.paqueteA = PaqueteCamion{}
        //fmt.Println("aLIMPIOOOO:",camion.paqueteA)
        camion.total--// considerar tope intentos
      }else{ // no se entrego pero quizs igual finalizo la orden
        if tipoPA == "retail"{
          if (camion.paqueteA).Intentos == 3{
            crearRegistroC(camion.paqueteA,camion.IDcamion)
            //ACTUAIZA Estado  NO ENTREGADO
            updateEstado("No recibido" , (camion.paqueteA).Seguimiento, (camion.paqueteA).Intentos , conn)
            camion.paqueteA = PaqueteCamion{}
            //fmt.Println("aLIMPIOOOO:",camion.paqueteA)
            camion.total--// considerar tope intento
          }
        }else{ // es pyme
            if (camion.paqueteA).Valor <=10{
              crearRegistroC(camion.paqueteA,camion.IDcamion)
              //ACTUAIZA Estado  NO ENTREGADO
              updateEstado("No recibido" , (camion.paqueteA).Seguimiento, (camion.paqueteA).Intentos , conn)
              camion.paqueteA = PaqueteCamion{}
              //fmt.Println("aLIMPIOOOO:",camion.paqueteA)
              camion.total--// considerar tope intento
            }else{
              if (camion.paqueteA).Intentos == 2{
                crearRegistroC(camion.paqueteA,camion.IDcamion)
                //ACTUAIZA Estado  NO ENTREGADO
                updateEstado("No recibido" , (camion.paqueteA).Seguimiento, (camion.paqueteA).Intentos , conn)
                camion.paqueteA = PaqueteCamion{}
                //fmt.Println("aLIMPIOOOO:",camion.paqueteA)
                camion.total--// considerar tope intento
              }
            }
        }
      }
      if (camion.paqueteB).Tipo != "" {
        flag = 0
      }
    }else{ // PAQUETE B
      (camion.paqueteB).Intentos++
      if entregado == 1{ //se entrega
        crearRegistroC(camion.paqueteB,camion.IDcamion)
            //ACTUALIZAR ESTADO ENTREGADO
        updateEstado("Recibido" , (camion.paqueteB).Seguimiento, (camion.paqueteB).Intentos , conn)
        camion.paqueteB = PaqueteCamion{}
        camion.total--// considerar tope intentos
      }else{ // no se entrego pero quizs igual finalizo la orden
        if tipoPB == "retail"{
          if (camion.paqueteB).Intentos == 3{
            crearRegistroC(camion.paqueteB,camion.IDcamion)
            //ACTUAIZA Estado  NO ENTREGADO
            updateEstado("No recibido" , (camion.paqueteB).Seguimiento, (camion.paqueteB).Intentos , conn)
            camion.paqueteB = PaqueteCamion{}
            camion.total--// considerar tope intento
          }
        }else{ // es pyme
            if (camion.paqueteB).Valor <=10{
              crearRegistroC(camion.paqueteB,camion.IDcamion)
              //ACTUAIZA Estado  NO ENTREGADO
              updateEstado("No recibido" , (camion.paqueteB).Seguimiento, (camion.paqueteB).Intentos , conn)
              camion.paqueteB = PaqueteCamion{}
              camion.total--// considerar tope intento
            }else{
              if (camion.paqueteB).Intentos == 2{
                crearRegistroC(camion.paqueteB,camion.IDcamion)
                //ACTUAIZA Estado  NO ENTREGADO
                updateEstado("No recibido" , (camion.paqueteB).Seguimiento, (camion.paqueteB).Intentos , conn)
                camion.paqueteB = PaqueteCamion{}
                camion.total--// considerar tope intento
              }
            }
        }
      }
      if (camion.paqueteA).Tipo != "" {
        flag = 1
      }
    }

  }//for
}
func main() {
  var conn *grpc.ClientConn
  conn, err := grpc.Dial(":9000", grpc.WithInsecure())
  if err != nil{
    log.Fatalf("could not connect: %s", err)
  }
  defer conn.Close()
  c := pb.NewLogisticaServiceClient(conn)

  var camionA = camiones {
    tipoC : "retail",
    IDcamion : 1,
    paqueteA : PaqueteCamion{},
    paqueteB : PaqueteCamion{},
    total : 0,
  }
  var camionB = camiones {
    tipoC : "retail",
    IDcamion : 2,
    paqueteA : PaqueteCamion{},
    paqueteB : PaqueteCamion{},
    total : 0,
  }
  var camionC = camiones {
    tipoC : "normal",
    IDcamion : 3,
    paqueteA : PaqueteCamion{},
    paqueteB : PaqueteCamion{},
    total : 0,
  }

  message := pb.Message{
    Body: "Hello from camiones",
  }

  response, err := c.SayHello(context.Background(),&message)
  if err!= nil{
    log.Fatalf("Error when calling SayHello: %s", err)
  }

  log.Printf("Response from Logistica: %s", response.Body)

  var tEspera int; // 0: retail 1: pyme
  log.Printf("Escriba segundos que esperará un camión el segundo paquete: (entero)")
  fmt.Scanln(&tEspera)
  var tEntrega int; // 0: retail 1: pyme
  log.Printf("Escriba segundos que se demorará la entrega de un paquete: (entero)")
  fmt.Scanln(&tEntrega)

  for{
    if camionA.total == 0{
      traerPaquetes(&camionA,tEspera,conn)
    }
    if camionB.total == 0{
      traerPaquetes(&camionB,tEspera,conn)
    }
    if camionC.total == 0{
      traerPaquetes(&camionC,tEspera,conn)
    }
    //--------------------
    if camionA.total != 0 {
      despacharPaquetes(&camionA,tEntrega,conn)
    }
    if camionB.total != 0 {
      despacharPaquetes(&camionB,tEntrega,conn)
    }
    if camionC.total != 0 {
      despacharPaquetes(&camionC,tEntrega,conn)
    }
    //Camiones cargados
  }
}
