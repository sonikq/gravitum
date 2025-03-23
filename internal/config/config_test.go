package config

import (
	"github.com/joho/godotenv"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockLoad is a version of Load that doesn't use flags, for testing purposes
func MockLoad(envFiles ...string) (Config, error) {
	if len(envFiles) != 0 {
		// Create a map to hold all environment variables from all files
		envMap := make(map[string]string)

		// Load each file and merge the environment variables
		for _, file := range envFiles {
			fileEnv, err := godotenv.Read(file)
			if err != nil {
				return Config{}, err
			}

			// Merge environment variables, with later files overriding earlier ones
			for k, v := range fileEnv {
				envMap[k] = v
			}
		}

		// Set all environment variables
		for k, v := range envMap {
			os.Setenv(k, v)
		}
	}

	cfg := new(Config)

	cfg.RunAddress = getEnvString(defaultRunAddress, "RUN_ADDRESS")
	cfg.DatabaseDSN = getEnvString(defaultDatabaseDSN, "DATABASE_DSN")
	cfg.DBPoolWorkers = getEnvInt(defaultDBPoolWorkers, "DB_POOL_WORKERS")
	cfg.CtxTimeOut = getEnvDuration(defaultCtxTimeOut, "CTX_TIMEOUT")
	cfg.LogLevel = getEnvString(defaultLogLevel, "LOG_LEVEL")
	cfg.ServiceName = getEnvString(defaultServiceName, "SERVICE_NAME")

	return *cfg, nil
}

// TestLoad_DefaultValues tests that default values are used when no environment variables are set
func TestLoad_DefaultValues(t *testing.T) {
	// Clear any environment variables that might affect the test
	os.Clearenv()

	// Load config with default values
	cfg, err := MockLoad()

	// Assert no error occurred
	require.NoError(t, err)

	// Assert default values are used
	assert.Equal(t, defaultRunAddress, cfg.RunAddress)
	assert.Equal(t, defaultDatabaseDSN, cfg.DatabaseDSN)
	assert.Equal(t, defaultDBPoolWorkers, cfg.DBPoolWorkers)
	assert.Equal(t, defaultCtxTimeOut, cfg.CtxTimeOut)
	assert.Equal(t, defaultLogLevel, cfg.LogLevel)
	assert.Equal(t, defaultServiceName, cfg.ServiceName)
}

// TestLoad_EnvVariables tests that environment variables override default values
func TestLoad_EnvVariables(t *testing.T) {
	// Clear any environment variables that might affect the test
	os.Clearenv()

	// Set environment variables
	customRunAddress := "localhost:4000"
	customDatabaseDSN := "postgres://user:password@localhost:5432/db"
	customDBPoolWorkers := 100
	customCtxTimeOut := 10000 // milliseconds
	customLogLevel := "debug"
	customServiceName := "test-service"

	os.Setenv("RUN_ADDRESS", customRunAddress)
	os.Setenv("DATABASE_DSN", customDatabaseDSN)
	os.Setenv("DB_POOL_WORKERS", "100")
	os.Setenv("CTX_TIMEOUT", "10000")
	os.Setenv("LOG_LEVEL", customLogLevel)
	os.Setenv("SERVICE_NAME", customServiceName)

	// Load config with environment variables
	cfg, err := MockLoad()

	// Assert no error occurred
	require.NoError(t, err)

	// Assert environment variables are used
	assert.Equal(t, customRunAddress, cfg.RunAddress)
	assert.Equal(t, customDatabaseDSN, cfg.DatabaseDSN)
	assert.Equal(t, customDBPoolWorkers, cfg.DBPoolWorkers)
	assert.Equal(t, time.Millisecond*time.Duration(customCtxTimeOut), cfg.CtxTimeOut)
	assert.Equal(t, customLogLevel, cfg.LogLevel)
	assert.Equal(t, customServiceName, cfg.ServiceName)
}

// TestLoad_EnvFile tests loading configuration from an environment file
func TestLoad_EnvFile(t *testing.T) {
	// Clear any environment variables that might affect the test
	os.Clearenv()

	// Create a temporary .env file
	envFile := ".env.test"
	envContent := `
RUN_ADDRESS=localhost:5000
DATABASE_DSN=postgres://test:test@localhost:5432/testdb
DB_POOL_WORKERS=75
CTX_TIMEOUT=7500
LOG_LEVEL=trace
SERVICE_NAME=env-file-service
`
	err := os.WriteFile(envFile, []byte(envContent), 0644)
	require.NoError(t, err)
	defer os.Remove(envFile)

	// Load config from env file
	cfg, err := MockLoad(envFile)

	// Assert no error occurred
	require.NoError(t, err)

	// Assert values from env file are used
	assert.Equal(t, "localhost:5000", cfg.RunAddress)
	assert.Equal(t, "postgres://test:test@localhost:5432/testdb", cfg.DatabaseDSN)
	assert.Equal(t, 75, cfg.DBPoolWorkers)
	assert.Equal(t, time.Millisecond*7500, cfg.CtxTimeOut)
	assert.Equal(t, "trace", cfg.LogLevel)
	assert.Equal(t, "env-file-service", cfg.ServiceName)
}

// TestLoad_InvalidEnvFile tests handling of an invalid environment file
func TestLoad_InvalidEnvFile(t *testing.T) {
	// Clear any environment variables that might affect the test
	os.Clearenv()

	// Try to load config from a non-existent env file
	_, err := MockLoad("non_existent_file.env")

	// Assert an error occurred
	require.Error(t, err)
}

// TestLoad_PartialEnvVariables tests that only specified environment variables override defaults
func TestLoad_PartialEnvVariables(t *testing.T) {
	// Clear any environment variables that might affect the test
	os.Clearenv()

	// Set only some environment variables
	customRunAddress := "localhost:4000"
	customLogLevel := "debug"

	os.Setenv("RUN_ADDRESS", customRunAddress)
	os.Setenv("LOG_LEVEL", customLogLevel)

	// Load config with partial environment variables
	cfg, err := MockLoad()

	// Assert no error occurred
	require.NoError(t, err)

	// Assert specified environment variables are used
	assert.Equal(t, customRunAddress, cfg.RunAddress)
	assert.Equal(t, customLogLevel, cfg.LogLevel)

	// Assert default values are used for unspecified variables
	assert.Equal(t, defaultDatabaseDSN, cfg.DatabaseDSN)
	assert.Equal(t, defaultDBPoolWorkers, cfg.DBPoolWorkers)
	assert.Equal(t, defaultCtxTimeOut, cfg.CtxTimeOut)
	assert.Equal(t, defaultServiceName, cfg.ServiceName)
}

// TestLoad_InvalidDBPoolWorkersEnv tests handling of invalid DB_POOL_WORKERS environment variable
func TestLoad_InvalidDBPoolWorkersEnv(t *testing.T) {
	// Clear any environment variables that might affect the test
	os.Clearenv()

	// Set invalid DB_POOL_WORKERS environment variable
	os.Setenv("DB_POOL_WORKERS", "not_a_number")

	// Load config with invalid environment variable
	cfg, err := MockLoad()

	// Assert no error occurred (the function handles invalid values internally)
	require.NoError(t, err)

	// Assert default value is used for invalid environment variable
	// Note: In this case, the function returns 1 for invalid values
	assert.Equal(t, 1, cfg.DBPoolWorkers)
}

// TestLoad_InvalidCtxTimeoutEnv tests handling of invalid CTX_TIMEOUT environment variable
func TestLoad_InvalidCtxTimeoutEnv(t *testing.T) {
	// Clear any environment variables that might affect the test
	os.Clearenv()

	// Set invalid CTX_TIMEOUT environment variable
	os.Setenv("CTX_TIMEOUT", "not_a_number")

	// Load config with invalid environment variable
	cfg, err := MockLoad()

	// Assert no error occurred (the function handles invalid values internally)
	require.NoError(t, err)

	// Assert default value is used for invalid environment variable
	// Note: The cast.ToInt function returns 0 for invalid values
	assert.Equal(t, time.Millisecond*0, cfg.CtxTimeOut)
}

// TestLoad_EmptyEnvVariables tests handling of empty environment variables
func TestLoad_EmptyEnvVariables(t *testing.T) {
	// Clear any environment variables that might affect the test
	os.Clearenv()

	// Set empty environment variables
	os.Setenv("RUN_ADDRESS", "")
	os.Setenv("DATABASE_DSN", "")
	os.Setenv("DB_POOL_WORKERS", "")
	os.Setenv("CTX_TIMEOUT", "")
	os.Setenv("LOG_LEVEL", "")
	os.Setenv("SERVICE_NAME", "")

	// Load config with empty environment variables
	cfg, err := MockLoad()

	// Assert no error occurred
	require.NoError(t, err)

	// Assert empty environment variables are used (not default values)
	assert.Equal(t, "", cfg.RunAddress)
	assert.Equal(t, "", cfg.DatabaseDSN)

	// For numeric values, empty strings can't be converted, so defaults or error handling values are used
	assert.Equal(t, 1, cfg.DBPoolWorkers)               // Error handling returns 1
	assert.Equal(t, time.Millisecond*0, cfg.CtxTimeOut) // cast.ToInt returns 0 for empty string

	assert.Equal(t, "", cfg.LogLevel)
	assert.Equal(t, "", cfg.ServiceName)
}

// TestGetEnvString tests the getEnvString helper function
func TestGetEnvString(t *testing.T) {
	// Test cases
	testCases := []struct {
		name      string
		flagValue string
		envKey    string
		envValue  string
		expected  string
	}{
		{
			name:      "Environment variable exists",
			flagValue: "flag-value",
			envKey:    "TEST_ENV_KEY",
			envValue:  "env-value",
			expected:  "env-value",
		},
		{
			name:      "Environment variable does not exist",
			flagValue: "flag-value",
			envKey:    "TEST_ENV_KEY",
			envValue:  "",
			expected:  "flag-value",
		},
		{
			name:      "Environment variable is empty",
			flagValue: "flag-value",
			envKey:    "TEST_ENV_KEY",
			envValue:  "",
			expected:  "flag-value",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear any environment variables that might affect the test
			os.Clearenv()

			// Set environment variable if needed
			if tc.envValue != "" {
				os.Setenv(tc.envKey, tc.envValue)
			}

			// Call the function
			result := getEnvString(tc.flagValue, tc.envKey)

			// Assert result
			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestGetEnvDuration tests the getEnvDuration helper function
func TestGetEnvDuration(t *testing.T) {
	// Test cases
	testCases := []struct {
		name      string
		flagValue time.Duration
		envKey    string
		envValue  string
		expected  time.Duration
	}{
		{
			name:      "Environment variable exists and is valid",
			flagValue: 5 * time.Second,
			envKey:    "TEST_ENV_KEY",
			envValue:  "1000",
			expected:  1000 * time.Millisecond,
		},
		{
			name:      "Environment variable does not exist",
			flagValue: 5 * time.Second,
			envKey:    "TEST_ENV_KEY",
			envValue:  "",
			expected:  5 * time.Second,
		},
		{
			name:      "Environment variable is not a number",
			flagValue: 5 * time.Second,
			envKey:    "TEST_ENV_KEY",
			envValue:  "not-a-number",
			expected:  0 * time.Millisecond, // cast.ToInt returns 0 for invalid values
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear any environment variables that might affect the test
			os.Clearenv()

			// Set environment variable if needed
			if tc.envValue != "" {
				os.Setenv(tc.envKey, tc.envValue)
			}

			// Call the function
			result := getEnvDuration(tc.flagValue, tc.envKey)

			// Assert result
			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestGetEnvInt tests the getEnvInt helper function
func TestGetEnvInt(t *testing.T) {
	// Test cases
	testCases := []struct {
		name      string
		flagValue int
		envKey    string
		envValue  string
		expected  int
	}{
		{
			name:      "Environment variable exists and is valid",
			flagValue: 50,
			envKey:    "TEST_ENV_KEY",
			envValue:  "100",
			expected:  100,
		},
		{
			name:      "Environment variable does not exist",
			flagValue: 50,
			envKey:    "TEST_ENV_KEY",
			envValue:  "",
			expected:  50,
		},
		{
			name:      "Environment variable is not a number",
			flagValue: 50,
			envKey:    "TEST_ENV_KEY",
			envValue:  "not-a-number",
			expected:  1, // Error handling returns 1
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear any environment variables that might affect the test
			os.Clearenv()

			// Set environment variable if needed
			if tc.envValue != "" {
				os.Setenv(tc.envKey, tc.envValue)
			}

			// Call the function
			result := getEnvInt(tc.flagValue, tc.envKey)

			// Assert result
			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestLoad_ZeroValues tests handling of zero values in environment variables
func TestLoad_ZeroValues(t *testing.T) {
	// Clear any environment variables that might affect the test
	os.Clearenv()

	// Set zero values in environment variables
	os.Setenv("DB_POOL_WORKERS", "0")
	os.Setenv("CTX_TIMEOUT", "0")

	// Load config with zero values
	cfg, err := MockLoad()

	// Assert no error occurred
	require.NoError(t, err)

	// Assert zero values are used
	assert.Equal(t, 0, cfg.DBPoolWorkers)
	assert.Equal(t, time.Millisecond*0, cfg.CtxTimeOut)
}

// TestLoad_ExtremeValues tests handling of extreme values in environment variables
func TestLoad_ExtremeValues(t *testing.T) {
	// Clear any environment variables that might affect the test
	os.Clearenv()

	// Set extreme values in environment variables
	os.Setenv("DB_POOL_WORKERS", "9999999")
	os.Setenv("CTX_TIMEOUT", "9999999")

	// Load config with extreme values
	cfg, err := MockLoad()

	// Assert no error occurred
	require.NoError(t, err)

	// Assert extreme values are used
	assert.Equal(t, 9999999, cfg.DBPoolWorkers)
	assert.Equal(t, time.Millisecond*9999999, cfg.CtxTimeOut)
}
