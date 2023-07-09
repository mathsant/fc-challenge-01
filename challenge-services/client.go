package challengeservices

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func Client(w http.ResponseWriter, r *http.Request) {
	// Cria um contexto com prazo de 300ms
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	// Cria uma requisição HTTP com o contexto
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatal("Erro ao criar a requisição:", err)
	}

	// Associa o contexto à requisição
	req = req.WithContext(ctx)

	// Faz a requisição HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Erro ao fazer a requisição:", err)
	}
	defer resp.Body.Close()

	// Lê o corpo da resposta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Erro ao ler o corpo da resposta:", err)
	}

	// Salva a cotação em um arquivo
	err = ioutil.WriteFile("cotacao.txt", body, 0644)
	if err != nil {
		log.Fatal("Erro ao salvar a cotação em arquivo:", err)
	}

	fmt.Println("Cotação salva em cotacao.txt")
}
