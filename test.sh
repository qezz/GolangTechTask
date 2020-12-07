#!/bin/bash

# Prep env

red="\e[31m"
green="\e[32m"
reset="\e[0m"
bold="\e[1m"

# Running server
echo -e "${bold}Running server${reset}"
go run cmd/server/main.go -p 8080 &
sleep 2

PID=$!


# Prepare Test data
echo -e "${bold}Preparing Test Data${reset}"


curl -X POST -d '{"question":"whos the man?", "answers": ["you", "me"]}' localhost:8080/buff -D -
curl -X POST -d '{"question":"which country?", "answers": ["UK", "US", "RU"]}' localhost:8080/buff -D -
curl -X POST -d '{"question":"which drink?", "answers": ["coffee", "vodka", "water"]}' localhost:8080/buff -D -


curl -X POST --silent -o /dev/null -d '{"name":"stream1", "buffs": [1]}' localhost:8080/stream -D -
curl -X POST --silent -o /dev/null -d '{"name":"stream2", "buffs": [2]}' localhost:8080/stream -D -
curl -X POST --silent -o /dev/null -d '{"name":"stream3", "buffs": [3]}' localhost:8080/stream -D -
curl -X POST --silent -o /dev/null -d '{"name":"stream4", "buffs": [1,2]}' localhost:8080/stream -D -
curl -X POST --silent -o /dev/null -d '{"name":"stream5", "buffs": [1,3]}' localhost:8080/stream -D -
curl -X POST --silent -o /dev/null -d '{"name":"stream6", "buffs": [2,3]}' localhost:8080/stream -D -
curl -X POST --silent -o /dev/null -d '{"name":"stream7", "buffs": [3,2]}' localhost:8080/stream -D -
curl -X POST --silent -o /dev/null -d '{"name":"stream8", "buffs": [3,1]}' localhost:8080/stream -D -
curl -X POST --silent -o /dev/null -d '{"name":"stream9", "buffs": [1,2]}' localhost:8080/stream -D -
curl -X POST --silent -o /dev/null -d '{"name":"stream10", "buffs": [1,3]}' localhost:8080/stream -D -
curl -X POST --silent -o /dev/null -d '{"name":"stream11", "buffs": [1]}' localhost:8080/stream -D -
curl -X POST --silent -o /dev/null -d '{"name":"stream12", "buffs": [2,1,3]}' localhost:8080/stream -D -
curl -X POST --silent -o /dev/null -d '{"name":"stream13", "buffs": [3]}' localhost:8080/stream -D -
curl -X POST --silent -o /dev/null -d '{"name":"stream14", "buffs": [1,3,2]}' localhost:8080/stream -D -
curl -X POST --silent -o /dev/null -d '{"name":"stream15", "buffs": [2]}' localhost:8080/stream -D -
curl -X POST --silent -o /dev/null -d '{"name":"stream16", "buffs": [3]}' localhost:8080/stream -D -
curl -X POST --silent -o /dev/null -d '{"name":"stream17", "buffs": [1,3,2]}' localhost:8080/stream -D -
curl -X POST --silent -o /dev/null -d '{"name":"stream18", "buffs": [2,1]}' localhost:8080/stream -D -
curl -X POST --silent -o /dev/null -d '{"name":"stream19", "buffs": [3]}' localhost:8080/stream -D -
curl -X POST --silent -o /dev/null -d '{"name":"stream20", "buffs": [1]}' localhost:8080/stream -D -


# Test Get Stream

if [ $(curl -s -o /dev/null -w "%{http_code}" -X GET  localhost:8080/stream/12) -eq 200 ]
then
    echo -e "${bold}===> Get Stream Test is ${green}Ok${reset}"
else
    echo -e "${bold}===> Get Stream Test is ${red}Failed${reset}"
fi


# Test Error
if [ "x$(curl --silent -X GET  localhost:8080/stream/100 )" == "xfailed to get stream, stream not found" ]
then
    echo -e "${bold}===> Get Stream Error Test is ${green}Ok${reset}"
else
    echo -e "${bold}===> Get Stream Error Test is ${red}Failed${reset}"
fi


# Test Pagination List

if [ $(curl --silent -X GET  localhost:8080/stream/page=2?pagesize=10 | grep id | wc -l) -eq 10 ]
then
    echo -e "${bold}===> Get Pagination List Test is ${green}Ok${reset}"
else
    echo -e "${bold}===> Get Pagination List Test is ${red}Failed${reset}"
fi



# Exit
echo -e "${bold}Exiting...${reset}"
kill -9 $PID
ps aux | grep "main -p 8080" | grep -v grep |awk '{print $2}' | xargs -I{} kill -9 {} > /dev/null 2>&1


