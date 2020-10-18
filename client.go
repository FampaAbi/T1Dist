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
)

func readCsvFile(filePath string) [][]string {
  csvfile, err := os.Open(filePath)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	r := csv.NewReader(csvfile)
  records, err := r.ReadAll()
  if err != nil {
        log.Fatal("Unable to parse file as CSV for ", err)
    }


  return records
}

func mostrarMenu() {
  fmt.Println("Bienvenido Cliente!")
  fmt.Println("Seleccione la acción que desea realizar:")
  fmt.Println("1. Ingresar archivo CSV")
  fmt.Println("2. Solicitar estado del paquete utilizando su número de seguimiento")
  fmt.Println("3. Salir")
}
func getEstado(nro int32, conn *grpc.ClientConn) {
  c := pb.NewLogisticaServiceClient(conn)
  estadito, _ := c.SeguirPaquete(context.Background(), &pb.Seguir{
    Seguimiento: int32(nro),
  })
  if estadito.GetSeguimiento() != 0 {
    fmt.Println("El estado para el paquete con nro de seguimiento ",nro," es ",estadito.GetEstado())
  }
}

func main() {
  var conn *grpc.ClientConn
  conn, err := grpc.Dial("dist61:9000", grpc.WithInsecure())
  if err != nil{
    log.Fatalf("could not connect: %s", err)
  }
  defer conn.Close()
  c := pb.NewLogisticaServiceClient(conn)

  message := pb.Message{
    Body: "Hello from the client!",
  }

  response, err := c.SayHello(context.Background(),&message)
  if err!= nil{
    log.Fatalf("Error when calling SayHello: %s", err)
  }
  var flag_menu = true
  for flag_menu {
    log.Printf("Response from Server: %s", response.Body)
    var opcion int;
    var nro_seg int;
    var tOrdenes int32;
    var flag int
    mostrarMenu()
    fmt.Scanln(&opcion)
    var path string
    if opcion == 1 {
      fmt.Printf("Qué tipo de orden desea realizar? [0: Pyme, 1: Retail]:")
      fmt.Scanln(&flag)
      if flag == 0 {
        path = "pymes.csv"
        } else if flag == 1 {
          path = "retail.csv"
          }else {
            log.Printf("Opción inválida")
          }
          if flag == 0 || flag == 1 {
            records := readCsvFile(path) // lista de listas de ordenes cliente
            //fmt.Println(records)
            //fmt.Println(reflect.ValueOf(records).Kind())
            log.Printf("Escriba segundos que habrá entre cada orden: (entero)")
            fmt.Scanln(&tOrdenes)
            var i =1 ;
            var tipoS int32
            for i < len(records){
              //fmt.Println(records[i])
              //fmt.Println(records[i][0])
              if flag == 0 {
                if records[i][5] == "0" {
                  tipoS = 0 // normal
                  } else {
                    tipoS = 1 //prioritario
                  }
                } else {
                    tipoS = 2 //retail
                }
              valor, _ := strconv.Atoi(records[i][2])
              response, err := c.GenerarOrden(context.Background(), &pb.Orden {
                IDPaquete: records[i][0],
                Tipo: tipoS,
                Nombre: records[i][1],
                Valor: int32(valor),
                Origen: records[i][3],
                Destino: records[i][4],
              })
            if err!= nil{
              log.Fatalf("Error al generar orden: %s", err)
            }
            fmt.Println("El número de seguimiento: ", response.NumeroSeguimiento," corresponde al ID de producto ",records[i][0])
            time.Sleep(time.Duration(tOrdenes)*time.Second)
            i++
          }
        }
      }else if opcion == 2{
        fmt.Println("Ingrese su número de seguimiento")
        fmt.Scanln(&nro_seg)
        getEstado(int32(nro_seg), conn)
      }else if opcion == 3 {
        flag_menu = false
      }else {
        fmt.Println("Ingrese una opción válida")
    }
  }
}
