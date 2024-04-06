1- Packages:
1.1- get apache package `go get -u -v github.com/go-sql-driver/mysql`
1.2- Gorilla Mux `go get github.com/gorilla/mux`

2- REDIS (usage: page count)
2.1- installation and start service: https://redis.io/docs/install/install-redis/install-redis-on-windows/
2.2- Some commands:
    - redis-cli
    - SET "home" 0
    - SET home 30 EX 5
    - ttl home
    - GET home
    - incr home
    - incrby home 3
    - del home
    - HSET monapi home 1 contact 1 download 0 help 1 #Définir Hash
    - HGETALL monapi
    - HGET monapi download
    - HINCRBY monapi dowload 1
    - flushall # Tout vider sur une machine de dév
    - PUBLISH "api:notifications:sport" "live matchs"
    - SUBSCRIBE "api:notif"