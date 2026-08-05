package config

type S3Config struct {
	Endpoint     string `env:"S3_ENDPOINT"`
	BucketName   string `env:"S3_BUCKET_NAME"`
	RootUser     string `env:"S3_ROOT_USER"`
	RootPassword string `env:"S3_ROOT_PASSWORD"`
	UseSSL       bool   `env:"S3_USE_SSL"`
}
