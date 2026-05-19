# SLO and Alerts

## Initial SLOs
- API availability: 99.9%
- p95 latency for read endpoints: < 300ms
- Verification success rate (non-not-found): > 99%

## Alert Conditions
- sustained 5xx error rate > 1%
- p95 latency > 500ms for 10 minutes
- spike in runtime security severity high/critical events
