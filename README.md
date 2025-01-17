### A simple middleware built with Go


##### File structure
```
api/
├── cmd/
│   └── main.go
├── internal/
│   ├── model/
│   │   └── user.go         # Domain models
│   ├── dto/
│   │   └── auth_dto.go     # Request/Response DTOs
│   ├── repository/
│   │   └── user_repo.go    # Database operations
│   ├── service/
│   │   └── auth_service.go # Business logic
│   ├── handler/
│   │   └── auth_handler.go # HTTP handlers
│   └── middleware/
│       └── auth.go         # JWT middleware
└── config/
    └── config.go           # App configuration

``` 

##### How to run application
````
go run cmd/main.go
```