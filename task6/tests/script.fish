#!/usr/bin/fish
for i in (seq 20); curl -s -o /dev/null -w "Status: %{http_code}\n" "http://localhost:8080/api/v1/repositories/info?url=github.com/golang/go"; end
