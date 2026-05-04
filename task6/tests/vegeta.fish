#!/usr/bin/fish
echo "GET http://localhost:8080/api/v1/subscriptions/info" | vegeta attack -rate=20 -duration=5s | vegeta report
