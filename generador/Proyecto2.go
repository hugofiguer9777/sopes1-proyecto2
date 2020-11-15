package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Datos struct {
	Name         string `json:"Name"`
	Location     string `json:"Location"`
	Age          int    `json:"Age"`
	Infectedtype string `json:"Infectedtype"`
	State        string `json:"State"`
}

func main() {
	//var entrada string;
	parametros, err2 := ioutil.ReadFile("parametros.txt")
	if err2 != nil {
		fmt.Printf("Error leyendo archivo de parametros: %v", err2)
	}
	contenido := strings.Split(string(parametros), "\n")
	bytesLeidos, err := ioutil.ReadFile(contenido[3])
	//contenido[1]=Hilos
	//contenido[2]=casoa a enviar
	if err != nil {
		fmt.Printf("Error leyendo archivo: %v", err)
		return
	}
	canthilos, errhilos := strconv.Atoi(contenido[1])
	if errhilos != nil {
		fmt.Println("Error con la cantidad de hilos")
		return
	}
	cantcasos, errcasos := strconv.Atoi(contenido[2])
	if errcasos != nil {
		fmt.Println("Error con la cantidad de casos")
		return
	}
	datosporhilo := cantcasos / canthilos
	//contenido := string(bytesLeidos)
	//fmt.Printf("La entrada fue:%s \n",contenido);
	var datos []Datos
	err1 := json.Unmarshal(bytesLeidos, &datos)
	if err1 != nil {
		fmt.Println("error transformando el json: %v", err1)
	}
	contador := 0
	//fmt.Print("valores resultantes:(hilos,casos,datosxhilo) ");
	//fmt.Print(canthilos);
	//fmt.Print(cantcasos);
	//fmt.Println(datosporhilo);
	for i := 0; i < canthilos; i++ {
		//fmt.Print("Haciendo hilo: ");
		//fmt.Println(i);
		var datosenviar []Datos
		for j := 0; j < datosporhilo; j++ {
			if contador >= len(datos) {
				contador = 0
			}
			datosenviar = append(datosenviar, datos[contador])
			contador++
		}
		go enviar(datosenviar, contenido[0])
	}
	time.Sleep(20 * time.Second)
	//fmt.Printf("Todo exitoso\n");
	//fmt.Printf("La entrada fue: %v \n",datos);
	//for i:=0;i<len(datos);i++{

	/*
		fmt.Printf("Valor en la posicion %d :\n",i);
		fmt.Printf("- %s \n",datos[i].Nombre);
		fmt.Printf("- %s \n",datos[i].Departamento);
		fmt.Printf("- %d \n",datos[i].Edad);
		fmt.Printf("- %s \n",datos[i].Forma);
		fmt.Printf("- %s \n",datos[i].Estado);
	*/
	//}
}

func enviar(datos []Datos, balanceador string) {

	for i := len(datos) - 1; i >= 0; i-- {
		enviar, err4 := json.Marshal(datos[i])
		if err4 != nil {
			fmt.Println("error marshal: %v", err4)
		}
		// fmt.Println("-----------IMPRIMIENDO HILO---------------")
		// fmt.Println(bytes.NewBuffer(enviar))

		resp, err3 := http.Post(balanceador, "application/json", bytes.NewBuffer(enviar))
		//http://contour1.p2sopes1.company:30032
		//http://nginx.p2sopes1.company:30036/1
		if err3 != nil {
			fmt.Println("error envio: %v", err3)
		}
		fmt.Println(resp)
		defer resp.Body.Close()
	}
}
