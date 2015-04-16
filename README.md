# Esgo
Simple elasticsearch client

Esgo is a ```http.Request``` wrapper. As simple as it possibly can be.

## Usage
```go
esgo.Configure("localhost", 9200)

resp, err := esgo.Request("POST", "/index/type/_search", `{
  "query": {
    "match_all": {}
  }
}`)
```
