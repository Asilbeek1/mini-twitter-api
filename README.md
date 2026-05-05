# mini-twitter-api

A small Twitter-like API built in Go with PostgreSQL, Redis, and JWT authentication.

## What it includes

- User registration and login
- JWT access + refresh tokens
- Post creation, update, delete, and feed listing
- PostgreSQL for storage
- Redis for token management and rate limiting
- Docker Compose config for easy local launch

## Quick start

1. This project uses environment variables for configuration.
You must create your own .env file before running the application
   ```bash
   #Database Credentials
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name

#Config path
CONFIG_PATH=config/local.yaml


#JWT Keys
ACCESS_KEY=your_access_key_here
REFRESH_KEY=your_refresh_key_here
   ```

2. Start the app and dependencies:
   ```bash
   sudo docker compose up --build
   ```

3. The API will be available at:
   ```text
   http://localhost:8080
   ```

## API endpoints

### Public

- `POST /api/v1/register` — create a new user
- `POST /api/v1/login` — get access + refresh tokens

### Authenticated

- `POST /api/v1/logout`
- `POST /api/v1/refresh`
- `GET /api/v1/users/{id}`
- `PUT /api/v1/users/{id}`
- `DELETE /api/v1/users/{id}`
- `GET /api/v1/posts/feed`
- `POST /api/v1/posts`
- `GET /api/v1/posts/{id}`
- `PUT /api/v1/posts/{id}`
- `DELETE /api/v1/posts/{id}`
- `GET /api/v1/users/{id}/posts`

## Example JSON payloads

Register:
```json
{
  "username": "testuser",
  "email": "test@example.com",
  "first_name": "Test",
  "password": "password123"
}
```

Login:
```json
{
  "email": "test@example.com",
  "password": "password123"
}
```

Refresh token:
```json
{
  "refresh_token": "your_refresh_token_here"
}
```

Update profile:
```json
{
  "first_name": "UpdatedName",
  "second_name": "UpdatedSecond",
  "email": "updated@example.com"
}
```

Create post:
```json
{
  "title": "My First Post",
  "description": "This is a description of my post."
}
```

Update post:
```json
{
  "title": "Updated Title",
  "description": "Updated description."
}
```

## Notes

- The app reads config from `config/local.yaml`.
- Docker Compose already includes PostgreSQL and Redis.
- For authenticated requests, add header:
  ```text
  Authorization: Bearer <access_token>
  ```

## Stop

```bash
sudo docker compose down
```
