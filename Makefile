ifdef port
app_port = $(port)
else
app_port := 8080
endif

ifdef addr
app_addr = $(addr)
else
app_addr := 127.0.0.1
endif

image:
	docker build -t ptask-server:latest .

run:
	docker run --name ptask-server -p $(app_addr):$(app_port):$(app_port) ptask-server:latest --addr=$(app_addr) --port=$(app_port)
    
start:
	docker start -i ptask-server

.PHONY: image run start