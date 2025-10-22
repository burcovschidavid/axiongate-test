#!/bin/bash

API_URL="http://localhost:38089"

if ! command -v jq &> /dev/null; then
    echo "jq not found, install it for pretty output: brew install jq"
    exit 1
fi

echo "========================================="
echo "Shipping API Test Suite"
echo "========================================="
echo ""

echo "1. Testing Provider A"
curl -s -X POST "${API_URL}/api/v1/createShipping?provider=A" \
  -H "Content-Type: application/json" \
  -d @sample-payload.json | jq '.'
echo ""
sleep 1

echo "2. Testing Provider B"
curl -s -X POST "${API_URL}/api/v1/createShipping?provider=B" \
  -H "Content-Type: application/json" \
  -d @sample-payload.json | jq '.'
echo ""
sleep 1

echo "3. Testing broadcast (both providers)"
curl -s -X POST "${API_URL}/api/v1/createShipping" \
  -H "Content-Type: application/json" \
  -d @sample-payload.json | jq '.'
echo ""
sleep 1

echo "4. Testing invalid provider"
curl -s -X POST "${API_URL}/api/v1/createShipping?provider=X" \
  -H "Content-Type: application/json" \
  -d @sample-payload.json | jq '.'
echo ""
sleep 1

echo "5. Testing health endpoint"
curl -s "${API_URL}/health"
echo ""
echo ""

echo "6. Checking database records"
docker-compose exec -T postgres psql -U postgres -d shipping -c "SELECT id, provider, success, created_at FROM shipment_records ORDER BY created_at DESC LIMIT 10;"
echo ""

echo "Done"
