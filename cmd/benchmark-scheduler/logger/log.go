// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package logger

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const defaultSampleValue = 100

type Logger interface {
	Info(string, ...interface{})
	Debug(string, ...interface{})
	Error(message string, err error)
	WithValues(...interface{}) Logger
}

type zapLogger struct {
	*zap.Logger
}

func New() Logger {
	cfg := zap.Config{
		Encoding:         "console",
		Development:      false,
		Sampling:         &zap.SamplingConfig{Initial: defaultSampleValue, Thereafter: defaultSampleValue},
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	z, err := cfg.Build(zap.AddCallerSkip(1), zap.AddCaller())
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	return &zapLogger{z}
}

func (z *zapLogger) Info(message string, keyAndValues ...interface{}) {
	if len(keyAndValues) == 0 {
		z.Logger.Info(message)
		return
	}
	z.Logger.Info(message, z.handleFields(keyAndValues)...)
}

func (z *zapLogger) Debug(message string, keyAndValues ...interface{}) {
	if len(keyAndValues) == 0 {
		z.Logger.Debug(message)
		return
	}
	z.Logger.Debug(message, z.handleFields(keyAndValues)...)
}

func (z *zapLogger) Error(message string, err error) {
	z.Logger.Error(message, zap.Error(err))
}

func (z *zapLogger) WithValues(keyAndValues ...interface{}) Logger {
	if len(keyAndValues) == 0 {
		return z
	}
	with := z.Logger.With(z.handleFields(keyAndValues)...)
	return &zapLogger{with}
}

func (z *zapLogger) handleFields(keyAndValues []interface{}) []zap.Field {
	fields := make([]zap.Field, 0, len(keyAndValues)/2) //nolint:gomnd //Reason:source is a pairs, but we need only half of them
	for i := 0; i < len(keyAndValues); i += 2 {
		if _, ok := keyAndValues[i].(zap.Field); ok {
			z.Logger.DPanic("strongly-typed Zap Field passed to logr", zap.Any("zap field", keyAndValues[i]))
			break
		}
		// make sure this isn't a mismatched key
		if i == len(keyAndValues)-1 {
			z.Logger.DPanic("odd number of arguments passed as key-value pairs for logging", zap.Any("ignored key", keyAndValues[i]))
			break
		}
		// process a key-value pair,
		// ensuring that the key is a string
		key, val := keyAndValues[i], keyAndValues[i+1]
		keyStr, isString := key.(string)
		if !isString {
			// if the key isn't a string, DPanic and stop logging
			z.Logger.DPanic("non-string key argument passed to logging, ignoring all later arguments", zap.Any("invalid key", key))
			break
		}
		fields = append(fields, zap.Any(keyStr, val))
	}
	return fields
}
