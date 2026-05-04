#!/usr/bin/fish
echo "Test run: RPS=20 Duration=5s"
echo "Results:"
echo "GET http://localhost:8080/api/v1/subscriptions/info" | \
vegeta attack -rate=20 -duration=5s | \
vegeta report -type=json | jq '.status_codes'
