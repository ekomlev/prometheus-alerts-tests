package test

import (
	"context"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

const (
	prometheusAddr = "https://prometheus-addres-example.io"
	alertKeyName   = "alertname"
	firingState    = "firing"
	up             = "up"
)

func TestGeneralQuery(t *testing.T) {
	t.Parallel()

	v1api := *getClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, warnings, err := v1api.Query(ctx, up, time.Now())

	assert.Nilf(t, err, "Error querying Prometheus:\n", err)
	assert.Nilf(t, warnings, "Warnings:\n", warnings)

	logger.Log(t, "Result:\n", result)
}

func TestAlerts(t *testing.T) {
	t.Parallel()

	v1api := *getClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	alerts, err := v1api.Alerts(ctx)
	assert.Nilf(t, err, "Error getting Alerts:\n", err)

	logger.Log(t, "Alerts:\n", alerts)

	// Example of getting alert
	kubeletTooManyPodsAlerts := make([]v1.Alert, 0)

	alertValueName := "TooManyPodsAlert"

	for _, alert := range alerts.Alerts {
		if string(alert.Labels[alertKeyName]) == alertValueName {
			kubeletTooManyPodsAlerts = append(kubeletTooManyPodsAlerts, alert)
		}
	}

	for _, alert := range kubeletTooManyPodsAlerts {
		assert.Equalf(t, firingState, string(alert.State), "Incorrect state of alert: ", alert)
	}
}

// Helpers

func getClient(t *testing.T) *v1.API {
	client, err := api.NewClient(api.Config{
		Address: prometheusAddr,
	})

	if err != nil {
		logger.Log(t, "Error creating client:\n", err)
		os.Exit(1)
	}

	v1api := v1.NewAPI(client)

	return &v1api
}
