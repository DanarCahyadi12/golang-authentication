
# Golang authentication

Golang authentication with JWT (Json Web Token) and structured folder.



## Tech Stack

- Golang

- MySQL



## Framework & Library

- Fiber (Http framework)

- GORM (ORM)

- Viper (Configuration)

- Golang Migrate (Database Migration)

- Go Playground Validator (Validation) 
## Configuration

All config is in `config.json` file.
## Run migrations

```bash
 migrate -database "mysql://<your_username>:<your_password>@tcp(<your_host>:<your_port>)/<your_database>?charset=utf8mb4&parseTime=true&loc=Local" -path database/migrations up
```

## API Reference

#### Create account

```http
  POST /signup
```

| Body field | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `name` | `string` | Required  |
| `email`| `string` | Requried |
| `password` | `string` | Required, minimum 8 character |

#### Login / signin

```http
  POST /auth
```

| Body field | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `email`      | `string` | Required |
| `password` | `string` | required |

#### Get token when access token is expired

```http
  GET /auth/token
```

## Run application

```bash
go run main.go
```