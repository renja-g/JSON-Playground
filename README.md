# JSON Playground API

## Project Definition

The JSON Playground provides a mock REST API to play around with JSON data.

## Playground Creation

To create a playground, send a POST request to `/playgrounds`.
- Automatic deletion after 30 minutes.
- Create an own SQLite database for the playground.
- Return the playground ID and a JWT that will be needed to interact with the playground.

## Endpoints

### Open Endpoints

| Method | Endpoints                                    |
| ------ | -------------------------------------------- |
| GET    | `/articles`                                  |
| GET    | `/articles/{articleId}`                      |
| GET    | `/articles/{articleId}/comments`             |
| GET    | `/articles/{articleId}/comments/{commentID}` |
| POST   | `/playgrounds`                               |

### Playground Endpoints
> All endpoints require a valid JWT.

| Method         | Endpoints                                                               |
| -------------- | ----------------------------------------------------------------------- |
| GET/POST       | `/playgrounds/{playgroundId}/articles`                                  |
| GET/PUT/DELETE | `/playgrounds/{playgroundId}/articles/{articleId}`                      |
| GET/POST       | `/playgrounds/{playgroundId}/articles/{articleId}/comments`             |
| GET/PUT/DELETE | `/playgrounds/{playgroundId}/articles/{articleId}/comments/{commentID}` |
| GET            | `/playgrounds/{playgroundId}/comments/{commentID}`                      |