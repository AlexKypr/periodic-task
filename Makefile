ifdef APP_PORT
PORT := $(APP_PORT)
else
PORT := 8080
endif

.PHONY: image run start

image:
	docker build -t inacc-server:latest .

run:
	docker run --name inacc-server -p $(PORT):8080 inacc-server:latest
#docker run --name inacc-server -p $(PORT):8080 inacc-server:latest
    
start:
	docker start -i inacc-server
