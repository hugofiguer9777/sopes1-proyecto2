package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"log"

	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
)

type caso struct {
	Name         string `bson:"name"`
	Location     string `bson:"location"`
	Age          int    `bson:"age"`
	Infectedtype string `bson:"infectedtype"`
	State        string `bson:"state"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// var ch *amqp.Channel
// var q amqp.Queue

func main() {

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", indexRoute).Methods("GET")
	router.HandleFunc("/nuevoCaso", nuevoCaso).Methods("POST")

	fmt.Println("Server on port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))

}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "API REST en Golang")
	fmt.Println("GET: /")
}

func nuevoCaso(w http.ResponseWriter, r *http.Request) {
	var casoNuevo caso

	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprint(w, "Datos no validos")
	}

	// Unmarshal
	err = json.Unmarshal(reqBody, &casoNuevo)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	body, err := json.Marshal(casoNuevo) //fmt.Sprintf("%+v", casoNuevo)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	conn, err := amqp.Dial("amqp://guest:guest@Rabbitmq:5672/")
	failOnError(err, "Fallo al conectar con RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Falo al abrir el canal")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"covid-cases", // nombre queue
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Fallo al declarar la cola")

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(string(body)),
		})
	log.Printf("Enviando: %s", body)
	failOnError(err, "Fallo al enviar el mensaje")

	w.Header().Set("content-type", "application/json")
	fmt.Fprint(w, "Valor recibido")
}
