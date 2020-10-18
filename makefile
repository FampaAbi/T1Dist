logistica:
	/sbin/service rabbitmq-server start
	go run logistica.go
finanzas:
	/sbin/service rabbitmq-server start
	go run finanzas.go
cliente:
	go run client.go
camion:
	go run camiones.go
