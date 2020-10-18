package main

import (
  "github.com/streadway/amqp"
  "encoding/json"
	"encoding/csv"
	"os"
	"strconv"
	"log"
	"fmt"
)


type finanzas struct{
  IDPaquete string
  Tipo string
  Valor string
  Intentos string
  Estado string
}

func failOnError(err error, msg string) {
  if err != nil {
    log.Fatalf("%s: %s", msg, err)
  }
}

func main() {
  // https://www.rabbitmq.com/tutorials/tutorial-one-go.html
  conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
  failOnError(err, "Failed to connect to RabbitMQ")
  defer conn.Close()


  ch, err := conn.Channel()
  failOnError(err, "Failed to open a channel")
  defer ch.Close()
  q, err := ch.QueueDeclare(
		"finances", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

  forever := make(chan bool)

  go func() {
    for d := range msgs {
      log.Printf("Received a message: %s", d.Body)
    }
  }()

  log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
  <-forever


  
  f, err := os.OpenFile("./registrosFIN.csv", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
  if err != nil {
    log.Printf("Unable to open file")
  }
  defer f.Close()
  timestamp := time.Now()
  timestamp2 := timestamp.Format("2006-01-02 00:00")

  var data [][]string //https://golangcode.com/write-data-to-a-csv-file/
                        //https://www.golangprograms.com/sample-program-to-create-csv-and-write-data.html

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
