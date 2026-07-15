module github.com/IvanDrf/work-hunter/users

go 1.25.1

require go.yaml.in/yaml/v3 v3.0.4

require (
	github.com/golang-migrate/migrate/v4 v4.19.1
	github.com/google/uuid v1.6.0
)

require github.com/lib/pq v1.10.9

require (
	github.com/DATA-DOG/go-sqlmock v1.5.2
	github.com/IvanDrf/work-hunter/pkg/common v0.0.0-20260611120708-a46d5979968d
	github.com/IvanDrf/work-hunter/pkg/user-api v0.0.0-20260715074652-0590c2144997
	github.com/jmoiron/sqlx v1.4.0
	github.com/stretchr/testify v1.10.0
	google.golang.org/grpc v1.80.0
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260120221211-b8f7ae30c516 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
