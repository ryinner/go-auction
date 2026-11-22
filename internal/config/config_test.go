package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string            // имя подтеста, видно в выводе -v
		env     map[string]string // ВСЕ три переменные задаём явно
		wantErr string            // подстрока в ошибке; "" = ошибки быть не должно
	}{
		{
			name: "valid",
			env: map[string]string{
				"DB_URL":     "postgres://u:p@localhost:5432/auction",
				"JWT_SECRET": strings.Repeat("a", 32),
				"PORT":       "8080",
			},
			wantErr: "",
		},
		{
			name: "missing DB_URL",
			env: map[string]string{
				"DB_URL":     "",
				"JWT_SECRET": strings.Repeat("a", 32),
				"PORT":       "8080",
			},
			wantErr: "DB_URL is empty",
		},
		{
			name: "missing JWT_SECRET",
			env: map[string]string{
				"DB_URL":     "postgres://u:p@localhost:5432/auction",
				"JWT_SECRET": "",
				"PORT":       "8080",
			},
			wantErr: "JWT_SECRET is empty",
		},
		{
			name: "not secure JWT_SECRET",
			env: map[string]string{
				"DB_URL":     "postgres://u:p@localhost:5432/auction",
				"JWT_SECRET": strings.Repeat("a", 31),
				"PORT":       "8080",
			},
			wantErr: "JWT_SECRET is not secure",
		},
		{
			name: "missing PORT",
			env: map[string]string{
				"DB_URL":     "postgres://u:p@localhost:5432/auction",
				"JWT_SECRET": strings.Repeat("a", 32),
				"PORT":       "",
			},
			wantErr: "PORT is empty",
		},
		{
			name: "out of range (lower) PORT",
			env: map[string]string{
				"DB_URL":     "postgres://u:p@localhost:5432/auction",
				"JWT_SECRET": strings.Repeat("a", 32),
				"PORT":       "0",
			},
			wantErr: "PORT must be between 1 and 65535",
		},
		{
			name: "out of range (bigger) PORT",
			env: map[string]string{
				"DB_URL":     "postgres://u:p@localhost:5432/auction",
				"JWT_SECRET": strings.Repeat("a", 32),
				"PORT":       "65536",
			},
			wantErr: "PORT must be between 1 and 65535",
		},
		{
			name: "not integer PORT",
			env: map[string]string{
				"DB_URL":     "postgres://u:p@localhost:5432/auction",
				"JWT_SECRET": strings.Repeat("a", 32),
				"PORT":       "abc",
			},
			wantErr: "PORT must be integer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			_, err := Load()

			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.ErrorContains(t, err, tt.wantErr)
		})
	}
}
