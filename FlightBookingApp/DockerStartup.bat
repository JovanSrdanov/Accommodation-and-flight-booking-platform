start http://localhost:4200/swagger/index.html
docker stop flights-app 
docker rmi -f flights-app 
docker compose -f docker-compose.local.yml up
pause