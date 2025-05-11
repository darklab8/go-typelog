Example of using with otlp https://opentelemetry.io/docs/languages/go/getting-started/

Utilized this command to simplify debugging, submitting otlp to alloy agent:
- ssh -t -L 4318:127.0.0.1:3000 darklab 'docker run -it --rm --network grafana -p 127.0.0.1:3000:3000 caddy:2.9.1 caddy reverse-proxy --from :3000 --to alloy-traces:4318'

how it looks like in json way
```
$ TYPELOG_JSON=true TYPELOG_SCOPES=true go run .
{"time":"2025-04-13T19:09:28.151092151+02:00","level":"INFO","msg":"started run","scope":"infra.main"}
{"time":"2025-04-13T19:09:32.797456967+02:00","level":"INFO","msg":"rolldice started","span_id":"adaa7d084fba8784","trace_id":"38f8c0869d753c3c484a57c654944b82","scope":"infra.main"}
{"time":"2025-04-13T19:09:32.797536016+02:00","level":"INFO","msg":"Alice is rolling the dice","result":1,"span_id":"adaa7d084fba8784","trace_id":"38f8c0869d753c3c484a57c654944b82","scope":"infra"}
{"time":"2025-04-13T19:09:32.797561367+02:00","level":"ERROR","msg":"Alice is rolling the dice","is_attr":true,"span_id":"adaa7d084fba8784","trace_id":"38f8c0869d753c3c484a57c654944b82","scope":"infra.main"}
```
without json:
```
time=2025-04-13T19:12:15.035+02:00 level=INFO msg="started run" scope=infra.main
time=2025-04-13T19:12:15.545+02:00 level=INFO msg="rolldice started" span_id=0c18debe976dd4d8 trace_id=f9bfc5d805a7960aaabbfe4aa2ca87da scope=infra.main
time=2025-04-13T19:12:15.545+02:00 level=INFO msg="Alice is rolling the dice" result=4 span_id=0c18debe976dd4d8 trace_id=f9bfc5d805a7960aaabbfe4aa2ca87da scope=infra
time=2025-04-13T19:12:15.545+02:00 level=ERROR msg="Alice is rolling the dice" is_attr=true span_id=0c18debe976dd4d8 trace_id=f9bfc5d805a7960aaabbfe4aa2ca87da scope=infra.main
```