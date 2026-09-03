# 005 - Structured logs

## Status
Proposed

## Context
Plain-text log lines are easy to write but hard to query once logs leave a developer's terminal: fields like level, timestamp or operation ID have to be parsed back out with regexes, and formats drift between call sites. We needed a log format that observability tooling (log aggregators, search, alerting) can index and filter on directly, without a parsing layer in between.

## Decision
We use `logrus` configured with a JSON formatter (`infra/monitor/logger.go`) so every log line is emitted as a JSON object, with an optional pretty-print mode for local development.

An example of application logs:

```
{"level":"info","msg":"Starting Doc server on http://localhost:1323","time":"2026-09-03T15:55:42.635Z"}
{"level":"info","msg":"Starting server on http://localhost:8484","time":"2026-09-03T15:55:42.635Z"}
{"level":"info","method":"GET","ms":3,"msg":"http request","operationId":"01a067fb-c70d-7c66-9511-80b41581c665","status":200,"time":"2026-09-03T15:55:48.881Z","uri":"/v1/catalog?page=1"}
{"level":"info","method":"GET","ms":1,"msg":"http request","operationId":"01a067fb-d864-7653-ae1d-91af07f05ba0","status":200,"time":"2026-09-03T15:55:53.318Z","uri":"/v1/catalog?page=2"}
{"level":"info","method":"GET","ms":1,"msg":"http request","operationId":"01a067fb-f271-71ba-83b0-df2ff764ec55","status":200,"time":"2026-09-03T15:55:59.986Z","uri":"/v1/catalog/PROD001"}
{"level":"info","method":"OPTIONS","ms":0,"msg":"http request","operationId":"01a067fc-9da3-7774-a4dd-4c62d5af46d0","status":204,"time":"2026-09-03T15:56:43.811Z","uri":"/v1/categories"}
{"codes":["CAT101","CAT102"],"level":"info","msg":"categories created","operationId":"01a067fc-9da3-7b90-9202-8cbd994f7f6c","time":"2026-09-03T15:56:43.814Z"}
{"level":"info","method":"POST","ms":2,"msg":"http request","operationId":"01a067fc-9da3-7b90-9202-8cbd994f7f6c","status":201,"time":"2026-09-03T15:56:43.814Z","uri":"/v1/categories"}
```

## Consequences
Logs can be ingested and indexed by observability tools as structured data, letting fields such as level, timestamp and [004-operation-id](004-operation-id) be filtered and searched on directly instead of parsed from free text. The trade-off is that raw log output is less readable in a plain terminal without pretty-print or a JSON-aware viewer, and every log call must go through the structured logger rather than ad-hoc `fmt.Println`-style output.
