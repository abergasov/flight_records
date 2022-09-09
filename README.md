### flight_records

#### Start service
```bash
make run
# or run on custom port
make port=8081 run
```

endpoint for calculate source and destination in list of flights
```go
POST /api/v1/track
```
example request
```shell
curl --request POST \
  --url http://127.0.0.1:8080/api/v1/track \
  --header 'Content-Type: application/json' \
  --data '[["IND", "EWR"], ["SFO", "ATL"], ["GSO", "IND"], ["ATL", "GSO"]]'
```