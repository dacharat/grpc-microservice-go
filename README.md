# grpc-microservice-go

### cURL Example
```
curl --location --request POST 'http://localhost:8080/instruction' \
--header 'Content-Type: application/json' \
--data-raw '{
    "instructions": [
        {
            "operator": "PUSH",
            "operand": 5
        },
        {
            "operator": "PUSH",
            "operand": 5
        },
        {
            "operator": "MUL"
        }
    ]
}'
```

```
curl --location --request POST 'http://localhost:8080/calculator' \
--header 'Content-Type: application/json' \
--data-raw '{
    "numbers": [
        1,2,3,4,5,6,7,8,9
    ]
}'
```