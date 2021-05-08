#!/bin/zsh

#domain="http://223.94.61.114:50080"
domain="http://127.0.0.1:8081"

curl -d "message=1111111111111" -X POST ${domain}/1


#curl -d "message=1111111111111" -X POST ${domain}/picmaker\?id\=1 --output txx1.jpg
#curl -d "message=2222222222222" -X POST ${domain}/picmaker\?id\=2 --output txx2.jpg
#curl -d "message=3333333333333" -X POST ${domain}/picmaker\?id\=3 --output txx3.jpg
#curl -d "message=4444444444444" -X POST ${domain}/picmaker\?id\=4 --output txx4.jpg

