module github.com/pmorelli92/how-i-test-microservices-2023

go 1.19

require (
	github.com/Netflix/go-env v0.0.0-20220526054621-78278af1949d
	github.com/go-chi/chi/v5 v5.0.8
	github.com/golang-migrate/migrate/v4 v4.15.2
	github.com/google/uuid v1.3.0
	github.com/jackc/pgx/v5 v5.2.0
)

replace github.com/golang-migrate/migrate/v4 => github.com/treuherz/migrate/v4 v4.15.3-0.20230126133952-c198cf34413e

require (
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jackc/pgerrcode v0.0.0-20201024163028-a0d42d470451 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/puddle/v2 v2.1.2 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	golang.org/x/crypto v0.0.0-20220829220503-c86fa9a7ed90 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/text v0.4.0 // indirect
)
