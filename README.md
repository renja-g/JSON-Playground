# JSON Playground API

## Project Definition

The JSON Playground API provides a platform to interact with JSON data and test JSON parsing libraries. It generates unique playgrounds that can be utilized for experimentation purposes.

## Playground Creation

To create a playground, send a POST request to `/playgrounds`. Optionally, you can create a playground with pre-populated data by sending a POST request to `/playgrounds?populated`.

- Playground creation includes:
  - Automatic deletion after 5 minutes of inactivity or 30 minutes after creation.
  - Saving each playground as a JSON file in the `playgrounds` directory.
  - Return of a JWT token required for playground access.

## Endpoints

### Create Playground
- **Endpoint**: `/playgrounds`
- **Method**: POST
```json
{
  "id": "string",
  "articles": [],
  "createdAt": "string",
  "usedAt": "string",
}
```

### Create Playground with Data
- **Endpoint**: `/playgrounds?populated`
- **Method**: POST

### Playground Articles
- **Endpoint**: `/{playgroundId}/articles`
- **Methods**: GET, POST

### Specific Article
- **Endpoint**: `/{playgroundId}/articles/{article_id}`
- **Methods**: GET, PUT, PATCH, DELETE

### Article Comments
- **Endpoint**: `/{playgroundId}/articles/{article_id}/comments`
- **Methods**: GET, POST

### Specific Comment
- **Endpoint**: `/{playgroundId}/articles/{article_id}/comments/{comment_id}`
- **Methods**: GET, PUT, PATCH, DELETE

## Schemas

### Playground
```json
{
  "id": "string",
  "articles": [],
  "createdAt": "string",
  "usedAt": "string",
}
```

### Article
```json
{
  "id": "string",
  "title": "string",
  "content": "string",
}
```

### Comment
```json
{
  "id": "string",
  "content": "string",
}
```
