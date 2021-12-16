In order to load balance : 
sudo docker-compose up -d --scale web-dynamic=3
sudo docker-compose logs -f web-dynamic
