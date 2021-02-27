// Copyright (c) 2021 Andres More

// Main

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"io/ioutil"

	"github.com/qri-io/jsonschema"
	log "github.com/sirupsen/logrus"
	cfg "github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"

	revision "github.com/appleboy/gin-revision-middleware"
	limit "github.com/aviddiviner/gin-limit"
	access "github.com/bu/gin-access-limit"
	security "github.com/danielkov/gin-helmet"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	stats "github.com/semihalev/gin-stats"

	"gorm.io/datatypes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Object to handle entities
type Object struct {
	gorm.Model
	Attributes datatypes.JSON
}

// DB database
var DB *gorm.DB

func setupDatabase() *gorm.DB {

	var db *gorm.DB
	var err error

	db, err = gorm.Open(sqlite.Open(cfg.GetString("DbUri")), &gorm.Config{}) // Logger: logger.Default.LogMode(logger.Info)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "setup",
			"topic": "database",
			"key":   err.Error(),
		}).Fatal("Could not open database")
	}

	var dbi *sql.DB

	dbi, err = db.DB()
	if err != nil {
		log.WithFields(log.Fields{
			"event": "setup",
			"topic": "interface",
			"key":   err.Error(),
		}).Fatal("Could not get database interface")
	}

	err = dbi.Ping()
	if err != nil {
		log.WithFields(log.Fields{
			"event": "setup",
			"topic": "ping",
			"key":   err.Error(),
		}).Fatal("Could not ping interface")
	}

	err = db.AutoMigrate(&Object{})
	if err != nil {
		log.WithFields(log.Fields{
			"event": "setup",
			"topic": "migrate",
			"key":   err.Error(),
		}).Fatal("Could not migrate database")
	}

	return db
}

// Schema attributes
var Schema jsonschema.Schema

// setupResource set up resource
func setupResource() string {

	schemaFile := cfg.GetString("SchemaName")
	resourceName := strings.TrimSuffix(schemaFile, filepath.Ext(schemaFile))

	log.WithFields(log.Fields{
		"event": "resource",
		"topic": "name",
		"key":   resourceName,
	}).Info("resource name is " + resourceName)

	schemaContent, err := ioutil.ReadFile(schemaFile)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "resource",
			"topic": "schema",
			"key":   err.Error(),
		}).Fatal("could not read schema")
	}

	Schema := &jsonschema.Schema{}
	if err := json.Unmarshal(schemaContent, Schema); err != nil {
		log.WithFields(log.Fields{
			"event": "resource",
			"topic": "unmarshal",
			"key":   err.Error(),
		}).Fatal("could not unmarshal schema")
	}

	return resourceName
}

// setupRouter route resource
func setupRouter() *gin.Engine {

	SetConfig()

	resource := setupResource()

	DB = setupDatabase()

	router := gin.New()
	file, err := os.OpenFile("api.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "logging",
			"topic": "file",
			"key":   err.Error(),
		}).Fatal("could not open log")
	}
	gin.DefaultWriter = io.MultiWriter(file) // os.Stdout

	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	router.Use(gin.Recovery())
	router.Use(revision.Middleware())
	router.Use(limit.MaxAllowed(cfg.GetInt("MaxAllowed")))
	router.Use(security.Default())
	router.Use(access.CIDR(cfg.GetString("AccessCidr")))
	router.Use(stats.RequestStats())
	router.Use(cors.Default())
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(requestid.New())

	pprof.Register(router)

	v1 := router.Group("/v1")
	{
		v1.GET("/health", HandleHealth)
		v1.GET("/"+resource, HandleGet)
		v1.GET("/"+resource+"/"+":id", HandleGetItem)
		v1.POST("/"+resource, HandlePost)
		v1.PATCH("/"+resource+"/:id", HandlePatch)
		v1.DELETE("/"+resource+"/:id", HandleDelete)
	}
	return router
}

func main() {

	router := setupRouter()

	if cfg.GetBool("UseSsl") {
		router.RunTLS(":"+cfg.GetString("HttpsPort"), cfg.GetString("CertFile"), cfg.GetString("KeyFile"))
	} else {
		router.Run(":" + cfg.GetString("HttpPort"))
	}
}
