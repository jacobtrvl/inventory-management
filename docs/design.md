## Features & Design

### RESTful APIs
- Add, Get, List, and Delete operations are supported
- Custom metrics endpoint is supported. This does not follow any standards.
  Implemented to demonstrate channels
- Simple pagination is supported
- Filters for List are not supported due to limitations of the DB.
  Filters are easy to write when the underlying DB supports querying. Since this is not the
  case here, I decided to skip supporting filters

### In-Memory DB Design
- A lightweight in-memory DB is implemented to demonstrate concurrency patterns
- The DB is a collection of tables, where each table can hold a map of key-value pairs

### Rate Limiter
- Very simple implementation of the Token Bucket Algorithm

### Error Handling and Status Codes
- It's important to pass correct error codes to the client. This requires the internal packages to also differentiate the error categories. Currently this is not well designed in this system and needs improvement. It's not a good practice to return HTTP status codes from the internal systems.

### Logging Format
- There is a structural difference in logs between slog & Gin. This can be solved by a custom log handler. 