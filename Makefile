## ---------- UTILS
.PHONY: help
help: ## Show this menu
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: clean
clean: ## Clean all temp files
	@sudo rm -rf cotacao_BRL_USD.db cotacao.txt nohup.out



## ----- MAIN
serve: ## put the server up
	@nohup go run server/server.go &

down: ## put the server down
	@ps aux | grep -E "/tmp/go-build.*/exe/server" | grep -v "grep" | cut -d " " -f2 | xargs kill
	@[ -f nohup.out ] && rm nohup.out || true

run: ## run the client
	@go run client/client.go
