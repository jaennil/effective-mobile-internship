sudo docker build -t 2 --build-arg GREETING=custom_greeting .
sudo docker run -p 8080:80 2
