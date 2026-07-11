package s3store

import (
	"testing"

	"github.com/mujhtech/s3ase/config"
	"github.com/stretchr/testify/require"
)

func TestResolveStorageConfigForR2(t *testing.T) {
	cfg := &config.Config{
		ObjectStorage: config.ObjectStorage{Provider: config.ObjectStorageProviderR2},
		R2: config.R2{
			AccountID: "account123", AccessKeyID: "access", SecretKey: "secret", Jurisdiction: "eu",
		},
	}

	resolved, err := resolveStorageConfig(cfg)
	require.NoError(t, err)
	require.Equal(t, config.ObjectStorageProviderR2, resolved.provider)
	require.Equal(t, "auto", resolved.region)
	require.Equal(t, "https://account123.eu.r2.cloudflarestorage.com", resolved.endpoint)
	require.False(t, resolved.usePathStyle)
}

func TestResolveStorageConfigRejectsUnknownProvider(t *testing.T) {
	_, err := resolveStorageConfig(&config.Config{
		ObjectStorage: config.ObjectStorage{Provider: "unknown"},
	})
	require.EqualError(t, err, `unsupported object storage provider "unknown"`)
}

func TestR2AlwaysUsesAutoRegion(t *testing.T) {
	store := &S3Store{provider: config.ObjectStorageProviderR2, region: "auto"}
	require.Equal(t, "auto", store.requestRegion("eu-west-2"))
}

func TestR2BucketCreationOmitsAWSLocationConstraint(t *testing.T) {
	r2Store := &S3Store{provider: config.ObjectStorageProviderR2}
	require.Nil(t, r2Store.createBucketInput("bucket", "auto").CreateBucketConfiguration)

	s3Store := &S3Store{provider: config.ObjectStorageProviderS3}
	input := s3Store.createBucketInput("bucket", "eu-west-2")
	require.NotNil(t, input.CreateBucketConfiguration)
	require.Equal(t, "eu-west-2", string(input.CreateBucketConfiguration.LocationConstraint))
}
