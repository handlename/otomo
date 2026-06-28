package trace

import (
	"context"
	"testing"

	"github.com/handlename/otomo/config"
	"github.com/stretchr/testify/assert"
)

func TestInitTracerDisabled(t *testing.T) {
	oldConfig := config.Config
	defer func() { config.Config = oldConfig }()

	config.Config.Otel.Enabled = false
	shutdown, err := InitTracer(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, shutdown)
	err = shutdown(context.Background())
	assert.NoError(t, err)
}

func TestInitTracerEnabledStdout(t *testing.T) {
	oldConfig := config.Config
	defer func() { config.Config = oldConfig }()

	config.Config.Otel.Enabled = true
	config.Config.Otel.Exporter = "stdout"
	config.Config.Otel.ServiceName = "test-service"
	shutdown, err := InitTracer(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, shutdown)
	err = shutdown(context.Background())
	assert.NoError(t, err)
}

func TestInitTracerUnsupportedExporter(t *testing.T) {
	oldConfig := config.Config
	defer func() { config.Config = oldConfig }()

	config.Config.Otel.Enabled = true
	config.Config.Otel.Exporter = "invalid"
	config.Config.Otel.ServiceName = "test-service"
	shutdown, err := InitTracer(context.Background())
	assert.Error(t, err)
	assert.Nil(t, shutdown)
}

func TestInitTracerEnabledOtlp(t *testing.T) {
	oldConfig := config.Config
	defer func() { config.Config = oldConfig }()

	config.Config.Otel.Enabled = true
	config.Config.Otel.Exporter = "otlp"
	config.Config.Otel.ServiceName = "test-service"
	shutdown, err := InitTracer(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, shutdown)
	err = shutdown(context.Background())
	assert.NoError(t, err)
}

