#!/bin/bash

ENDPOINT="http://localhost:8080/api/v1/short"

generate_random_url() {
  echo "https://example$(($RANDOM % 10000)).com"
}

NUM_REQUESTS=5000

for ((i=1; i<=NUM_REQUESTS; i++))
do
  RANDOM_URL=$(generate_random_url)
  DATA=$(printf '{"original": "%s"}' "$RANDOM_URL")
  
  RESPONSE=$(curl -s -w "%{http_code}" -o /dev/null -X POST -H "Content-Type: application/json" -d "$DATA" "$ENDPOINT")
  
  if [ "$RESPONSE" -eq 200 ]; then
    echo "Request $i: Success"
  else
    echo "Request $i: Failed with status code $RESPONSE"
  fi
done
