package logger

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	originalEnv := os.Getenv("LOG_LEVEL")
	defer os.Setenv("LOG_LEVEL", originalEnv)

	oldDefault := slog.Default()
	defer slog.SetDefault(oldDefault)

	tests := []struct {
		name          string
		envLogLevel   string
		expectedLevel slog.Level
	}{
		{
			name:          "debug level",
			envLogLevel:   "DEBUG",
			expectedLevel: slog.LevelDebug,
		},
		{
			name:          "info level",
			envLogLevel:   "INFO",
			expectedLevel: slog.LevelInfo,
		},
		{
			name:          "warn level",
			envLogLevel:   "WARN",
			expectedLevel: slog.LevelWarn,
		},
		{
			name:          "error level",
			envLogLevel:   "ERROR",
			expectedLevel: slog.LevelError,
		},
		{
			name:          "invalid level 1",
			envLogLevel:   "INVALID",
			expectedLevel: slog.LevelInfo,
		},
		{
			name:          "empty level 2",
			envLogLevel:   "",
			expectedLevel: slog.LevelInfo,
		},
		{
			name:          "lowercase debug",
			envLogLevel:   "debug",
			expectedLevel: slog.LevelDebug,
		},
		{
			name:          "mixed case",
			envLogLevel:   "InFo",
			expectedLevel: slog.LevelInfo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("LOG_LEVEL", tt.envLogLevel)
            _ = Init()

            currentLogger := slog.Default()
            ctx := context.Background()

            assert.True(t, currentLogger.Handler().Enabled(ctx, tt.expectedLevel), "Level %s should be active ", tt.envLogLevel)

            var levelBelow slog.Level
            switch tt.expectedLevel {
            case slog.LevelInfo:
                levelBelow = slog.LevelDebug
            case slog.LevelWarn:
                levelBelow = slog.LevelInfo
            case slog.LevelError:
                levelBelow = slog.LevelWarn
            default:
                levelBelow = -5
            }

            if levelBelow != -5 {
                assert.False(t, currentLogger.Handler().Enabled(ctx, levelBelow))
            }
		})
	}
}
