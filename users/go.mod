module github.com/IvanDrf/work-hunter/users

go 1.25.1

require go.yaml.in/yaml/v3 v3.0.4

require (
	github.com/golang-migrate/migrate/v4 v4.19.1
	github.com/google/uuid v1.6.0
)

require github.com/lib/pq v1.10.9

require (
	github.com/IvanDrf/work-hunter/pkg/user-api v0.0.0-20260401110126-ac223fd7dc2f
	github.com/jmoiron/sqlx v1.4.0
	google.golang.org/grpc v1.80.0
	google.golang.org/protobuf v1.36.11
)

require (
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260120221211-b8f7ae30c516 // indirect
)
