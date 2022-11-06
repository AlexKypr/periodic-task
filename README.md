## How to run assignment

Prerequisites:
1. Have Docker installed

Steps:
1. Unzip file and go to inaccess folder
2. Run cmd `make image` to create docker image
3. Run cmd `make run` to create and run docker container. The default address is 127.0.0.1 and port 8080. If you want to use custom ones then you have to run cmd `make run addr=custom_addr port=custom_port`
4. (optional) To re-run container run cmd `make start`.