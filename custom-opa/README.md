# Customized OPA

As seen in https://sched.co/ekEV

## Building With Docker

```
docker build -t custom-opa .
```

## Running

### Example config

__conf.yaml__
```yaml
decision_logs:
  plugin: kafka_logger
plugins:
  kafka_logger:
    topic: opa
    host: localhost:9092
  grpc_api:
    listen: :8123
```

With kafka running locally on 9092 and starting the grpc listener on localhost:8123

> Hint: Check out [https://kafka.apache.org/quickstart](https://kafka.apache.org/quickstart) to get Kafka up and running locally

```
go run . run -s -c conf.yaml -l debug ./policy.rego
```

> Note: The gRPC API handler has a hard coded query so be careful modifying the `policy.rego`.

### gRPC client

There is a test client available to try out plugin, check out [plugins/api/client](./plugins/api/client)

Test with via:

```
cd plugins/api/client
go run main.go localhost:8123 <JWT>
```

This will send a request with the specified JWT. You can modify the `policy.rego` to use a different URL for the JWKS, ideally use one that will verify the JWT provided. Some samples are given below, but they are not going to work indefinitely.

```
go run . localhost:8123 "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IlI0MjdnTG03ek9mTVJwc1hoM1ZnNCJ9.eyJpc3MiOiJodHRwczovL2Rldi12cGRhMWw0ZS51cy5hdXRoMC5jb20vIiwic3ViIjoiWk9XMGZTWTAxeVAwMnRKeEJVcWxxUXRBeVZPQzlSQWJAY2xpZW50cyIsImF1ZCI6Imh0dHBzOi8vb3BhNGZ1bi5jb20vbXlhcHAiLCJpYXQiOjE2MDMyMjY5NzIsImV4cCI6MTYwMzMxMzM3MiwiYXpwIjoiWk9XMGZTWTAxeVAwMnRKeEJVcWxxUXRBeVZPQzlSQWIiLCJndHkiOiJjbGllbnQtY3JlZGVudGlhbHMifQ.OeYtC_LG9TrmtRUiSBZZGtgbrXJXP89fd4lJgTd2cV1cNX0gw6gtPDc5IheDkKYccnggGFNy95tMXSOvv_PMKwowk3u3--YQ22SEvtmSpZl9UIZRp1nK0cI2xXH3EmM_9BOOjg2EUvImEmxLiHtQqR9cquwndrsv2ZsngRXhs0t1QRuySJrlNEeJSViJy4zIU6ElabWM5M0_JEtyglX2jjMrymJVYUb-GZrK83c0AkdMgiFE3XG20z9gPgrGt5Vm3fTPmnv5YRB_CTtX1tcV8UKJc1uU99KsgUzjAjzwKgutHTsP3tgDKjcXfm-q_O7Llv1aS9qKs0eBf9itDKq6vQ"
```