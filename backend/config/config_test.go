package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_DefaultConfig(t *testing.T) {
	cfg, err := LoadConfig()
	require.NoError(t, err)

	expectedCfg := DefaultConfig

	require.Equal(t, expectedCfg, cfg)
}

func Test_LoadConfig(t *testing.T) {
	tests := []struct {
		name       string
		envVars    map[string]string
		wantErr    bool
		wantErrMsg string
		validate   func(*testing.T, *Config)
	}{
		{
			name:    "load_default_config",
			wantErr: false,
			validate: func(t *testing.T, cfg *Config) {
				require.Equal(t, DefaultConfig, cfg)
			},
		},
		{
			name: "custom_database_config",
			envVars: map[string]string{
				"DB_HOST":     "custom-host",
				"DB_PORT":     "5433",
				"DB_USER":     "custom-user",
				"DB_PASSWORD": "custom-pass",
				"DB_DATABASE": "custom-db",
			},
			validate: func(t *testing.T, cfg *Config) {
				require.Equal(t, "custom-host", cfg.Database.Host)
				require.Equal(t, 5433, cfg.Database.Port)
				require.Equal(t, "custom-user", cfg.Database.User)
				require.Equal(t, "custom-pass", cfg.Database.Password)
				require.Equal(t, "custom-db", cfg.Database.Database)
			},
		},
		{
			name: "custom_server_config",
			envVars: map[string]string{
				"PORT":    "8080",
				"SSL":     "true",
				"TIMEOUT": "30",
			},
			validate: func(t *testing.T, cfg *Config) {
				require.Equal(t, uint32(8080), cfg.Server.Port)
				require.True(t, cfg.Server.SSL)
				require.Equal(t, uint32(30), cfg.Server.Timeout)
			},
		},
		{
			name: "custom_aws_config",
			envVars: map[string]string{
				"AWS_DEFAULT_REGION": "us-east-1",
				"AWS_ACCESS_KEY":     "test-access-key",
				"AWS_SECRET_KEY":     "test-secret-key",
			},
			validate: func(t *testing.T, cfg *Config) {
				require.Equal(t, "us-east-1", cfg.Aws.DefaultRegion)
				require.Equal(t, "test-access-key", cfg.Aws.AccessKey)
				require.Equal(t, "test-secret-key", cfg.Aws.SecretKey)
			},
		},
		{
			name: "invalid_port_number",
			envVars: map[string]string{
				"PORT": "invalid",
			},
			wantErr: true,
		},
		{
			name: "empty_database_config",
			envVars: map[string]string{
				"DB_HOST": "",
				"DB_PORT": "0",
			},
			wantErr:    true,
			wantErrMsg: "database host cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset environment before each test
			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}

			cfg, err := LoadConfig()

			if tt.wantErr {
				require.Error(t, err)
				if tt.wantErrMsg != "" {
					require.Equal(t, tt.wantErrMsg, err.Error())
				}
				return
			}

			require.NoError(t, err)
			require.NotNil(t, cfg)

			if tt.validate != nil {
				tt.validate(t, cfg)
			}
		})
	}
}
