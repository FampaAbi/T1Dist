# sd-t1

Integrantes:
	Camilo farah, rol 201773534-1
	Christian Sepúlveda, rol 201791003-8

Para ejecutar:
	Para iniciar el servidor RabbitMQ, ejecutar el siguiente comando tanto en ssh root@dist54 y en ssh root@dist56:
		>> /sbin/service rabbitmq-server start

	Luego, cada maquina virtual (VM) tiene asignada una de las cuatro entidades. Por lo tanto, se debe ejecutar el siguiente comando para cada uno de ellas:
		- Para ssh root@dist53:  
		>> make clientes
		
		- Para ssh root@dist54:
		>> make logistica

		- Para ssh root@dist55:
		>> make camiones

		- Para ssh root@dist56:
		>> make finanza

Consideraciones:
	- Si desea utilizar sus propios datos para los pedidos de paquetes de clientes, modifique los archivos csv ("pymes.csv" y/o "retail.csv") como corresponda:
		- La primera línea del archivo se ignorará, pues equivaldrá a los headers. Luego:
			- para retail: id,producto,valor,tienda,destino
			- para pymes: id,producto,valor,tienda,destino,prioritario

	