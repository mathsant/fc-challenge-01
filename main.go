package main

import (
	"log"
	"net/http"

	challengeservices "github.com/mathsant/fc-challenge-01/challenge-services"
)

func main() {
	// Cria o servidor HTTP
	http.HandleFunc("/cotacao", challengeservices.CotacaoHandler)
	http.HandleFunc("/client", challengeservices.Client)

	// Inicia o servidor na porta 8080
	log.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
