clean:
	go run main.go -drop "all"
desafio:
	go run main.go -cep "90010310,04119062,22070900,40080004,69010140" -fileurls "desafio.csv" -collection "products" 

supermercado:
	go run main.go -cep "06020010" -scan "1" -collection "supermercado" -fileurls "supermercado.csv"     