---
description: gRPC Standards and Best Practices
---

# gRPC Standards and Best Practices

**Message Naming:**
- UpperCamelCase for messagenames
- For RPC methods, use `{MethodName}Request` and `{MethodName}Response` pattern (e.g., GetUserRequest, GetUserResponse)
- Do not add an unnecessary suffix like `Message`

**Field Naming:**
- Use lower_snake_case for field names (e.g., `user_id`, `error_reason`)
- Use plural forms for repeated fields (e.g. `books`, `song_ids`)
- Avoid including prepositions (e.g., prefer `error_reason` over `reason_for_error`)
- Adhere to American English spelling and established abbreviations (e.g., `id` for `identifier`, `url` for `uniform resource locator`)
- Place numbers after the word component, not after an underscore (e.g., `song_name1` not `song_name_1`)

**General Conventions and Pitfalls:**
- Field names should not use language or framework reserved keywords to avoid code-generation issues
- When defining standard fields, reuse names established by Google’s API design where applicable (e.g., `created_time`, `update_time`, `name` for resource ID)
- For oneof fields, use lower_snake_case as well
- Enum values should be `ALL_UPPER_SNAKE_CASE`, and the zero value usually ends with `_UNSPECIFIED`
- Be mindful of tag numbers: Once assigned, they must never be reused for a different field name, even if the field is deleted (use `reserved` for removed tags/names)
- Message names are not transmitted on the wire except with `Any` fields, so renaming them is not a protocol-breaking change, but clients must be updated accordingly

**Field Deprecation:**
- Use `[deprecated = true]` option to mark fields as deprecated (e.g., `string old_field = 3 [deprecated = true];`)
- MUST NOT delete or reuse deprecated field tag numbers — reserve them via `reserved` after removal
- Add a comment explaining the deprecation reason and migration path
- Deprecated fields SHOULD remain in proto until all clients migrate

**REST Annotations (gRPC-Gateway):**
- Import `"google/api/annotations.proto"` and depend on `"google/api/http.proto"`
- Use `option (google.api.http)` within each RPC method block
- Map HTTP verbs correctly: GET for retrieval, POST for creation, PUT/PATCH for updates, DELETE for deletions
- Use versioned paths with hierarchical structure (e.g., `/v1/users/{user_id}`)
- Bind request fields to path parameters with braces `{...}`
- For GET/DELETE: omit body or use `body: ""`
- For POST/PUT/PATCH: use `body: "*"` or specify specific request field
- Response messages should match expected JSON output structure

**Example:**
```proto
syntax = "proto3";

package books.v1;

import "google/api/annotations.proto";
import "google/protobuf/timestamp.proto";

service BookService {
  rpc GetBook (GetBookRequest) returns (GetBookResponse) {
    option (google.api.http) = {
      get: "/v1/books/{book_id}"
    };
  }
  
  rpc CreateBook (CreateBookRequest) returns (CreateBookResponse) {
    option (google.api.http) = {
      post: "/v1/books"
      body: "*"
    };
  }
  
  rpc UpdateBook (UpdateBookRequest) returns (UpdateBookResponse) {
    option (google.api.http) = {
      patch: "/v1/books/{book_id}"
      body: "book"
    };
  }
  
  rpc DeleteBook (DeleteBookRequest) returns (google.protobuf.Empty) {
    option (google.api.http) = {
      delete: "/v1/books/{book_id}"
    };
  }
}

message GetBookRequest {
  string book_id = 1;
}

message GetBookResponse {
  Book book = 1;
}

message CreateBookRequest {
  string name = 1;
  string author = 2;
  repeated string tags = 3;
}

message CreateBookResponse {
  Book book = 1;
}

message UpdateBookRequest {
  string book_id = 1;
  Book book = 2;
}

message UpdateBookResponse {
  Book book = 1;
}

message DeleteBookRequest {
  string book_id = 1;
}

message Book {
  string id = 1;
  string name = 2;
  string author = 3;
  google.protobuf.Timestamp created_time = 4;
  google.protobuf.Timestamp updated_time = 5;
  repeated string tags = 6;
  
  // Deprecated: use 'author' field instead. Will be removed in v2.
  string author_name = 7 [deprecated = true];
  
  reserved 8, 9;  // Previously: old_field, legacy_data
  reserved "old_field", "legacy_data";
}
```