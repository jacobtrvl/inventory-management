## Inventory Management Microservice
For design & feature considerations, see [docs/design.md](docs/design.md).

### Environment Variables
```bash
export ADDR=":8080"
```
By default ADDR is set to :8080

### Running Service
```bash
# Build and run
make run
```

The service will start on `http://127.0.0.1:8080` or on ADDR set in environment variable


### Running unit tests

```bash
make test
```

### API Endpoints & Payload



#### Add Product
POST /products

- id is optional. System will generate a UID if id is not supplied


Request body format:

```json
{
  "id": "string (optional)",
  "name": "string (optional)",
  "price": "number (optional)",
  "stock": "integer (optional)"
}
```

Example:

```bash
curl --location 'http://127.0.0.1:8080/products' \
--header 'Content-Type: application/json' \
--data '{
    "id": "12",
  "name": "Laptop",
  "price": 999.99,
  "stock": 10
}'
```

#### List products (with pagination)
GET /products
```bash
curl --location 'http://127.0.0.1:8080/products'
```
##### Pagination
```bash
curl --location 'http://127.0.0.1:8080/products?page=2&limit=1'
```
#### Get based on ID
GET /products/<id>

```
curl --location 'http://127.0.0.1:8080/products/12'
```

#### Update product
PUT /products/<id>

Request body format:
```json
{
  "name": "string (optional)",
  "price": "number (optional)",
  "stock": "integer (optional)"
}
```
Example:
```bash
curl --location --request PUT 'http://127.0.0.1:8080/products/1' \
--header 'Content-Type: application/json' \
--data '{
    "price": 12

}'
```

#### Delete Product by ID
DELETE /products/<id>

```bash
curl --location --request DELETE 'http://127.0.0.1:8080/products/12'
```

#### Get metrics
GET /metrics
```
curl --location 'http://127.0.0.1:8080/metrics'
```

### Testing rate limiter

```
make ratelimiter
```

### Docker build & run

```bash
docker build . -t <image_name>
docker run --rm -p 8080:8080 <image_name>
```