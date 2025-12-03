---
description: REST API Standards and Best Practices
---

# REST API Standards and Best Practices

**Endpoint Naming:**
- Use nouns, not verbs (e.g., `/users` not `/getUsers`)
- Use plural nouns for collections (e.g., `/users`, `/orders`)
- Use lowercase with hyphens (kebab-case) for multi-word resources (e.g., `/user-profiles`)
- Keep URIs short and readable, avoid abbreviations
- Do not include file extensions in URIs (e.g., no `.json`)
- Include API version at base path (e.g., `/v1/users`)

**HTTP Methods:**
- `POST /api/v1/my-items` - creation
- `GET /api/v1/my-items/{id}` - retrieval
- `PATCH /api/v1/my-items/{id}` - update
- `DELETE /api/v1/my-items/{id}` - deletion
- `GET /api/v1/my-items` - list with pagination

**Response Structure:**
- Use JSON as standard data format
- Return appropriate HTTP status codes (200, 201, 204, 400, 401, 403, 404, 500)
- Provide consistent error responses with meaningful messages
- Use standard fields like `id`, `created_time`, `updated_time` for resources

**Query Parameters:**
- Use for filtering, sorting, and pagination (e.g., `/users?role=admin&page=2&limit=10`)
- Do not embed filtering in the path
- Use snake_case for parameter names (e.g., `user_id`, `created_after`)

**General Conventions:**
- Minimize nesting depth (avoid more than 2-3 levels)
- Use hierarchical structure only when clear relationship exists (e.g., `/users/{id}/orders`)
- Maintain consistency across all endpoints
- Do not expose internal implementation details in URLs
- Use American English spelling and established abbreviations

**Example:**
```
GET /v1/users
POST /v1/users
GET /v1/users/{user_id}
PATCH /v1/users/{user_id}
DELETE /v1/users/{user_id}
GET /v1/users/{user_id}/orders
```

Response format:
```json
{
  "id": "user_123",
  "name": "John Doe",
  "email": "john@example.com",
  "created_time": "2023-01-01T00:00:00Z",
  "updated_time": "2023-01-02T12:30:00Z"
}
```

Error response:
```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "User not found",
    "details": "User with ID 'invalid_id' does not exist"
  }
}
```