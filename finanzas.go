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

func terminarFinanzas(dinero finanzas) float32{
  var total float32
  total = gananciaPaquete(dinero)

  f, err := os.OpenFile("./registrosFIN.csv", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
  if err != nil {
    log.Printf("Unable to open file")
  }
  defer f.Close()


  var data [][]string //https://golangcode.com/write-data-to-a-csv-file/
                        //https://www.golangprograms.com/sample-program-to-create-csv-and-write-data.html
  l := fmt.Sprintf("%f", total)
  data = append(data, []string{
    dinero.IDPaquete,
    dinero.Tipo,
    dinero.Valor,
    dinero.Intentos,
    dinero.Valor,
    l,})

  writer := csv.NewWriter(f)
  writer.WriteAll(data)
  return total
}

func gananciaPaquete(dinero finanzas) float32{
  var estado = dinero.Estado
  i, _ := strconv.Atoi(dinero.Intentos)
  var perdida = (i - 1) * 10
  var suma float32
  if dinero.Tipo == "prioritario"{
      if estado == "Recibido"{
        i, _ := strconv.Atoi(dinero.Valor)
        suma = (float32(i) *1.3) - float32(perdida)
        return suma
      }else{
        i, _ := strconv.Atoi(dinero.Valor)
        suma = (float32(i) * 0.3) - float32(perdida)
        return suma
      }
  }else if dinero.Tipo == "normal"{
    if estado == "Recibido"{
      i, _ := strconv.Atoi(dinero.Valor)
      suma = float32(i - perdida)
      return suma
    }else{
      suma= -float32(perdida)
      return suma
    }
  }else{
    i, _ := strconv.Atoi(dinero.Valor)
    suma = float32(i - perdida)
    return suma
  }
}
func main() {
  // https://www.rabbitmq.com/tutorials/tutorial-one-go.html
  conn, err := amqp.Dial("amqp://dist61:dist61@dist61:9000/")
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
  var total float32
  forever := make(chan bool)

  go func() {
    for d := range msgs {
      var dinero finanzas
      json.Unmarshal(d.Body, &dinero)
      var totalTemp = terminarFinanzas(dinero)
      total += totalTemp
      log.Printf("Received a message: %s", d.Body)
    }
  }()

  log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
  <-forever





}
