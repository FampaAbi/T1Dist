(se ve mejor del browser)

Sebastián Rojas 201773598-8
Fabio Pazos 201773503-1


Considerar -> 	dist61 : logistica
		dist62 : finanzas
		dist63 : clientes
		dist64 : camiones

Para ejecutar:
	-Ejecutar make run en cada maquina ejecutará el código correspondiente.
	-Tanto en la máquina asociada a logística como a finanzas ademas del correspondiente go run file.go se ejecuta
	el siguiente comando /sbin/service rabbitmq-server start, el cual inicializa el server de rabbit.
	

Consideraciones:
	- El timestamp utilizado para la creación de los diferentes registros presentó problemas, ya que el formateo generaba que la hora permaneciera fija
	  en 00:00, a pesar de los esfuerzos y los múltiples tutoriales revisados.
	- Se asumió que la fecha de entrega asociada a los paquetes era 0 en caso de que estos no fueran entregados.
	- En el csv de ejemplo (pymes) se trabajó el booleano de la siguiente forma:
		1: Prioritario
		0: Normal
	- El archivo proto se entrega compilado para facilitar la revisión de la tarea.
	- A tener en cuenta que se decidió debido al contexto del problema no optar por resetear la creación de los diferentes csv de registros asociados a
	  camiones, finanzas y logística, ya que en estos se guardaba información necesaria para la construcción de reportes "históricos" como el balance
	  global de la "empresa".
	
	

	
