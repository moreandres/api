// Copyright (c) 2021 Andres More

// config

package main

import (
	"io"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	cfg "github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

// https://github.com/dgrijalva/jwt-go

// SetConfig sets configuration
func SetConfig() {
	cfg.SetDefault("LogLevel", "Info")
	cfg.SetDefault("SchemaName", "objects.json")
	cfg.SetDefault("DbUri", "file::memory:?cache=shared")
	cfg.SetDefault("HttpPort", "8080")
	cfg.SetDefault("UseSSL", false)
	cfg.SetDefault("HttpsPort", "8443")
	cfg.SetDefault("CertFile", "api.cer")
	cfg.SetDefault("KeyFile", "api.key")
	cfg.SetDefault("URL", "https://localhost.com")
	cfg.SetDefault("QueryLimit", "512")
	cfg.SetDefault("JwtSecret", "password")
	cfg.SetDefault("MaxAllowed", 20)
	cfg.SetDefault("AccessCidr", "0.0.0.0/0")

	cfg.SetConfigName("api")
	cfg.AddConfigPath(".")
	cfg.AutomaticEnv()
	cfg.SetEnvPrefix("api")

	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{})

	file, err := os.OpenFile("api.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "config",
			"topic": "logging",
			"key":   err.Error(),
		}).Fatal("Could not open logging file")
	}

	log.SetOutput(io.MultiWriter(os.Stdout, file))

	err2 := cfg.ReadInConfig()
	if err2 != nil {
		log.WithFields(log.Fields{
			"event": "config",
			"topic": "read",
			"key":   err2.Error(),
		}).Warn("could not read config")
	}

	setLogLevel()

	cfg.WatchConfig()
	cfg.OnConfigChange(handleConfigChange)
}

func setLogLevel() {

	levels := map[string]log.Level{
		"debug": log.DebugLevel,
		"info":  log.InfoLevel,
		"warn":  log.WarnLevel,
		"error": log.ErrorLevel,
		"fatal": log.FatalLevel,
		"panic": log.PanicLevel,
	}

	level := levels[strings.ToLower(cfg.GetString("LogLevel"))]
	log.SetLevel(level)
}

// handleConfigChange handle configuration changes
func handleConfigChange(e fsnotify.Event) {

	log.WithFields(log.Fields{
		"event": "config",
		"topic": "change",
		"key":   e.Name,
	}).Info("configuration has changed")

	setLogLevel()
}
