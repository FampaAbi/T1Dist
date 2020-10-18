package main

import (
  "golang.org/x/net/context"
  "google.golang.org/grpc"

  "github.com/Tarea1/Express/logistica"

  "encoding/csv"
  "fmt"
  "log"
  "os"
  "bufio"
  "io"
)


func readCsvFile() int{
    f, err := os.Open("retail.csv")
    if err != nil {
        log.Fatal("Unable to read input file " + filePath, err)
    }
    defer f.Close()

    csvReader := csv.NewReader(f)
    for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Question: %s Answer %s\n", record[0], record[1])
	}

    return 1
}
func main(){
  // Open the file
  	csvfile, err := os.Open("retail.csv")
  	if err != nil {
  		log.Fatalln("Couldn't open the csv file", err)
  	}

  	// Parse the file
  	r := csv.NewReader(csvfile)
  	//r := csv.NewReader(bufio.NewReader(csvfile))

  	// Iterate through the records
  	for {
  		// Read each record from csv
  		record, err := r.Read()
  		if err == io.EOF {
  			break
  		}
  		if err != nil {
  			log.Fatal(err)
  		}
  		fmt.Printf("Question: %s Answer %s\n", record[0], record[1])
  	}
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
