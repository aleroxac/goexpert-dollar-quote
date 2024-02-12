package main

import (
	"context"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	SERVER_ENDPOINT = "http://localhost:8080/cotacao"
	TIMEOUT_API     = 300 * time.Millisecond
	QUOTE_FILE      = "cotacao.txt"
)

type QuoteResponse struct {
	Bid float64
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_API)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", SERVER_ENDPOINT, nil)
	if err != nil {
		log.Fatalf("Fail to create the request: %v", err)
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Fail to make the request: %v", err)
		return
	}
	defer res.Body.Close()

	ctx_err := ctx.Err()
	if ctx_err != nil {
		select {
		case <-ctx.Done():
			err := ctx.Err()
			log.Fatalf("Max timeout reached: %v", err)
			return
		}
	}

	resp_json, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Fail to read the response: %v", err)
		return
	}

	var quote QuoteResponse
	err = json.Unmarshal(resp_json, &quote)
	if err != nil {
		log.Fatalf("Fail to decode the response: %v", err)
		return
	}
	log.Printf("quote: %f", quote.Bid)

	file, err := os.Create(QUOTE_FILE)
	if err != nil {
		log.Fatalf("Fail to create %s: %v", QUOTE_FILE, err)
		return
	}

	tmpl, err := template.New("output").Parse("Dólar: {{.Bid}}")
	if err != nil {
		log.Fatalf("Error creating template: %v", err)
		return
	}

	file_content := QuoteResponse{
		Bid: quote.Bid,
	}

	err = tmpl.Execute(file, file_content)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
		return
	}
}
