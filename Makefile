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
	docker build -t inacc-server:latest .

run:
	docker run --name inacc-server -p $(app_addr):$(app_port):8080 inacc-server:latest --port=$(app_port)
    
start:
	docker start -i inacc-server

.PHONY: image run start