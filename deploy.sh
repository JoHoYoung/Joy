docker rm -f joy
docker build -t joy:0.1 .
docker run -it --name joy -p 8080:8080 joy:0.1
