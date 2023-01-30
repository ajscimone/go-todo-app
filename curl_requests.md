# Curl commands

## Get All Events

curl -i -H "Accept: application/json" -H "Content-Type: application/json" -X GET http://localhost:8080/events

## Post Event

curl -d '{"Title":"Eat Pasta", "Description":"Make spaghetti and eat it"}' -H "Content-Type: application/json" -X POST http://localhost:8080/event
