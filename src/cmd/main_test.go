package main

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLoggerInitialization(t *testing.T) {
	logger := logrus.New()

	assert.NotNil(t, logger, "Logger should not be nil")
}
