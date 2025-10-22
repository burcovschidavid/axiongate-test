# Shipping API Integration Service

Multi-provider shipping API gateway that transforms generic shipping requests into provider-specific formats.

## Quick Start

### Using Docker

```bash
docker-compose up --build
```

The API will be available at `http://localhost:38089`


## API Usage

### Create Shipment with Specific Provider

```bash
curl -X POST "http://localhost:38089/api/v1/createShipping?provider=A" \
  -H "Content-Type: application/json" \
  -d @payload.json
```

### Broadcast to All Providers

```bash
curl -X POST "http://localhost:38089/api/v1/createShipping" \
  -H "Content-Type: application/json" \
  -d @payload.json
```

### Health Check

```bash
curl http://localhost:38089/health
```

## Request Format

The API accepts a generic shipping model:

```json
{
  "weight": {
    "value": 1000,
    "unit": "Grams"
  },
  "shipper": {
    "contact": {
      "name": "Sender Name",
      "mobileNumber": "0506356566",
      "phoneNumber": "041234567",
      "emailAddress": "sender@test.com",
      "companyName": "Sender Company"
    },
    "address": {
      "line1": "Address Line 1",
      "city": "Dubai",
      "countryCode": "AE",
      "zipCode": "00000"
    }
  },
  "consignee": {
    "contact": {
      "name": "Receiver Name",
      "mobileNumber": "+919441234567",
      "emailAddress": "receiver@test.com"
    },
    "address": {
      "line1": "Receiver Address",
      "city": "Bangalore",
      "countryCode": "IN",
      "zipCode": "1001"
    }
  },
  "dimensions": {
    "length": 10,
    "height": 10,
    "width": 10,
    "unit": "Meter"
  },
  "account": {
    "number": "123"
  },
  "productCode": "International",
  "numberOfPieces": 1,
  "declaredValue": {
    "amount": 250,
    "currency": "AED"
  }
}
```


## Adding New Providers

1. Create new package under `internal/adapters/providers/{provider}/`
2. Implement provider-specific models and mapper
3. Create adapter implementing `ports.ShippingProvider`
4. Register in `cmd/api/main.go`

## Environment Variables

- `PORT` - Server port (default: 8080)
- `DB_HOST` - PostgreSQL host
- `DB_PORT` - PostgreSQL port
- `DB_USER` - Database user
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name
- `PROVIDER_A_URL` - Provider A endpoint
- `PROVIDER_B_URL` - Provider B endpoint

## Database

PostgreSQL with JSONB columns for flexible payload storage.

Run migrations:
```bash
docker-compose up migrate
```


## Testing

The project includes comprehensive test coverage

### Quick API Test
```bash
./test-api.sh
```


