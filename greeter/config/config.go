package config

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

var (
	once     sync.Once
	smClient *secretsmanager.Client
	region   = "ap-southeast-1"
)

func initSecretManager(ctx context.Context) {
	once.Do(func() {
		cfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(region))
		if err != nil {
			panic(err)
		}
		smClient = secretsmanager.NewFromConfig(cfg)
	})
}

// UnmarshalSecret reads the secret manager and returns the value from secret manager.
func UnmarshalSecret(ctx context.Context, secretID string, v any) error {
	initSecretManager(ctx)
	out, err := smClient.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretID),
		VersionStage: aws.String("AWSCURRENT"),
	})
	if err != nil {
		return err
	}
	buf := []byte(*out.SecretString)
	return json.Unmarshal(buf, &v)
}
