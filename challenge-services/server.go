package challengeservices

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Cotacao struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

func obterCotacao() (string, error) {
	// Cria um contexto com prazo de 200ms
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	// Cria uma requisição HTTP com o contexto
	req, err := http.NewRequest(http.MethodGet, "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return "", fmt.Errorf("erro ao criar a requisição: %v", err)
	}

	// Associa o contexto à requisição
	req = req.WithContext(ctx)

	// Faz a requisição HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro ao fazer a requisição: %v", err)
	}
	defer resp.Body.Close()

	// Lê o corpo da resposta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("erro ao ler o corpo da resposta: %v", err)
	}

	// Decodifica o JSON da resposta
	var cotacao Cotacao
	err = json.Unmarshal(body, &cotacao)
	if err != nil {
		return "", fmt.Errorf("erro ao decodificar o JSON: %v", err)
	}

	return cotacao.USDBRL.Bid, nil
}

func salvarCotacaoNoBanco(cotacao string) error {
	// Abre uma conexão com o banco de dados SQLite
	db, err := sql.Open("sqlite3", "cotacoes.db")
	if err != nil {
		return fmt.Errorf("erro ao abrir o banco de dados: %v", err)
	}
	defer db.Close()

	// Cria a tabela se ela não existir
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS cotacoes (cotacao TEXT)")
	if err != nil {
		return fmt.Errorf("erro ao criar a tabela: %v", err)
	}

	// Insere a cotação na tabela
	_, err = db.Exec("INSERT INTO cotacoes (cotacao) VALUES (?)", cotacao)
	if err != nil {
		return fmt.Errorf("erro ao inserir a cotação no banco de dados: %v", err)
	}

	return nil
}

func CotacaoHandler(w http.ResponseWriter, r *http.Request) {
	// Obtém a cotação
	cotacao, err := obterCotacao()
	if err != nil {
		log.Println("Erro ao obter a cotação:", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}

	// Salva a cotação no banco de dados
	err = salvarCotacaoNoBanco(cotacao)
	if err != nil {
		log.Println("Erro ao salvar a cotação no banco de dados:", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}

	// Monta a resposta em formato JSON
	resposta := map[string]string{
		"cotacao": cotacao,
	}
	jsonResposta, err := json.Marshal(resposta)
	if err != nil {
		log.Println("Erro ao serializar a resposta:", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}

	// Define o cabeçalho da resposta
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResposta)
}
