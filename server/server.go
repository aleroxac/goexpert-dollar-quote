package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	API_PORT      = 8080
	TIMEOUT_API   = 200 * time.Millisecond
	TIMEOUT_DB    = 10 * time.Millisecond
	QUOTE_API_URL = "https://economia.awesomeapi.com.br/json/last"
	SRC_COIN      = "BRL"
	DEST_COIN     = "USD"
)

type Quote struct {
	Usdbrl struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

type QuoteResponse struct {
	Bid float64 `json:"bid" gorm:"primaryKey"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received a request on: /cotacao")

	log.Println("Creating a request")
	ctx_api, cancel := context.WithTimeout(context.Background(), TIMEOUT_API)
	defer cancel()

	request_method := "GET"
	request_endpoint := fmt.Sprintf("%s/%s-%s", QUOTE_API_URL, DEST_COIN, SRC_COIN)
	req, err := http.NewRequestWithContext(ctx_api, request_method, request_endpoint, nil)
	if err != nil {
		log.Fatalf("Fail to create the request: %v", err)
		http.Error(w, fmt.Sprintf("Fail to create the request: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("Making a request")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Fail to make the request: %v", err)
		http.Error(w, fmt.Sprintf("Fail to create the request: %v", err), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	ctx_api_err := ctx_api.Err()
	if ctx_api_err != nil {
		select {
		case <-ctx_api.Done():
			err := ctx_api.Err()

			log.Fatalf("Max timeout reached: %v", err)
			http.Error(w, fmt.Sprintf("Max timeout reached: %v", err), http.StatusRequestTimeout)
			return
		}
	}

	log.Println("Reading the response")
	resp_json, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Fail to read the response: %v", err)
		http.Error(w, fmt.Sprintf("Fail to read the response: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("Decoding the response")
	var quote Quote
	err = json.Unmarshal(resp_json, &quote)
	if err != nil {
		log.Fatalf("Fail to decode the response: %v", err)
		http.Error(w, fmt.Sprintf("Fail to decode the response: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("Converting the quote value")
	w.Header().Set("Content-Type", "application/json")
	quote_float, err := strconv.ParseFloat(quote.Usdbrl.Bid, 64)
	if err != nil {
		log.Fatalf("Fail convert quote value: %v", err)
		http.Error(w, fmt.Sprintf("Fail convert the quote value: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("Checking db file")
	quote_response := QuoteResponse{
		Bid: quote_float,
	}
	cotacao_db_file := fmt.Sprintf("cotacao_%s_%s.db", SRC_COIN, DEST_COIN)
	_, err = os.Stat(cotacao_db_file)
	if err != nil {
		_, err = os.Create(cotacao_db_file)
		if err != nil {
			log.Fatalf("Fail to create %s: %v", cotacao_db_file, err)
			http.Error(w, fmt.Sprintf("Fail to create %s: %v", cotacao_db_file, err), http.StatusInternalServerError)
			return
		}
	}

	log.Println("Opening connection with db")
	db, err := sql.Open("sqlite3", cotacao_db_file)
	if err != nil {
		log.Fatalf("Fail to connecto to db: %v", err)
		http.Error(w, fmt.Sprintf("Fail to create %s: %v", cotacao_db_file, err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	log.Println("Creating a table on db")
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cotacao (id INTEGER PRIMARY KEY AUTOINCREMENT, cotacao FLOAT)`)
	if err != nil {
		log.Fatalf("Fail create table on db: %v", err)
		http.Error(w, fmt.Sprintf("Fail to create table on db: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("Inserting data on db")
	ctx_db, cancel := context.WithTimeout(context.Background(), TIMEOUT_DB)
	defer cancel()

	stmt, err := db.PrepareContext(ctx_db, "INSERT INTO cotacao (cotacao) VALUES (?)")
	if err != nil {
		log.Fatalf("Fail prepare db statement: %v", err)
		http.Error(w, fmt.Sprintf("Fail prepare db statement: %v", err), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(quote_response.Bid)
	if err != nil {
		log.Fatalf("Fail to insert data on db: %v", err)
		http.Error(w, fmt.Sprintf("Fail to insert data on db: %v", err), http.StatusInternalServerError)
		return
	}

	ctx_db_err := ctx_db.Err()
	if ctx_db_err != nil {
		select {
		case <-ctx_db.Done():
			err := ctx_db.Err()

			log.Fatalf("Max timeout reached: %v", err)
			http.Error(w, fmt.Sprintf("Max timeout reached: %v", err), http.StatusRequestTimeout)
			return
		}
	}

	log.Println("Encoding response")
	json_resp, err := json.Marshal(quote_response)
	if err != nil {
		log.Fatalf("Fail to encode response: %v", err)
		http.Error(w, fmt.Sprintf("Fail to encode response: %v", err), http.StatusInternalServerError)
		return
	}
	w.Write(json_resp)
	log.Println("quote:", quote.Usdbrl.Bid)
}

func main() {
	http.HandleFunc("/cotacao", handler)
	log.Printf("Listening on %d", API_PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", API_PORT), nil))
}
