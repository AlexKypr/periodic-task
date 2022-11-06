## How to run assignment

Prerequisites:
1. Have Docker installed

Steps to run server:
1. Unzip file and go to inaccess folder
2. Run cmd `make image` to create docker image
3. Run cmd `make run` to create and run docker container. The default address is 127.0.0.1 and port 8080. If you want to use custom ones then you have to run cmd `make run addr=custom_addr port=custom_port`
4. (optional) To re-run container run cmd `make start`.

Testing Example:
1. make image
2. make run addr=127.0.0.5 port=8082
3. curl --location --request GET 'http://127.0.0.5:8082/ptlist?period=1h&tz=Europe/Athens&t1=20210214T204603Z&t2=20211115T123456Z'
