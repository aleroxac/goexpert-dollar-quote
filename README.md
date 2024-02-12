# goexpert-dollar-quote
Primeiro desafio do treinamento GoExpert(FullCycle).



## O desafio
Criar um `server` e um `client` de modo que grave um arquivo `cotacao.txt` com e em uma tabela de um banco `SQLITE` a cotação atual do dólar(`USD`) em relação ao real(`BRL`) a partir de uma `API` pública que fornece tal cotação.



## Como rodar o projeto
``` shell
make serve  # sobe o server
make run    # roda o client
make down   # baixa o server
make clean  # remove o arquivo do db e de cotacao
```



## Funcionalidades da Linguagem Utilizadas
- [x] context
- [x] net/http
- [x] encoding/json
- [x] database/sql
- [x] html/template



## Requisitos
- [x] O client.go precisará receber do server.go apenas o valor atual do câmbio (campo "bid" do JSON). Utilizando o package "context", o client.go terá um timeout máximo de 300ms para receber o resultado do server.go.
- [x] O server.go deverá consumir a API contendo o câmbio de Dólar e Real no endereço: https://economia.awesomeapi.com.br/json/last/USD-BRL e em seguida deverá retornar no formato JSON o resultado para o cliente. 
- [x] Usando o package "context", o server.go deverá registrar no banco de dados SQLite cada cotação recebida, sendo que o timeout máximo para chamar a API de cotação do dólar deverá ser de 200ms e o timeout máximo para conseguir persistir os dados no banco deverá ser de 10ms.
- [x] O client.go deverá realizar uma requisição HTTP no server.go solicitando a cotação do dólar.
- [x] Os 3 contextos deverão retornar erro nos logs caso o tempo de execução seja insuficiente.
- [x] O client.go terá que salvar a cotação atual em um arquivo "cotacao.txt" no formato: Dólar: {valor}
- [x] O endpoint necessário gerado pelo server.go para este desafio será: /cotacao e a porta a ser utilizada pelo servidor HTTP será a 8080.
