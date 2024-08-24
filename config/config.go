package config

type Config struct {
    BucketName      string `env:"BUCKET_NAME,required"`
    S3Endpoint      string `env:"S3_ENDPOINT,required"`
    AccessKeyID     string `env:"ACCESS_KEY_ID,required"`
    SecretAccessKey string `env:"SECRET_ACCESS_KEY,required"`
    Port            string `env:"PORT,default=3000"`
}

