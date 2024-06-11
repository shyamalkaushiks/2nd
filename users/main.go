package main

import (
	"fmt"
	"os"
	"users/config"
	"users/logger"
	model "users/models"
	"users/services"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	var err error
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false})
	err = logger.SetupLogger(logger.Log)
	if err != nil {
		log.Error().Err(err).Msg("error while setup logger")
		return
	}
	//load config

	err = config.LoadConfig()
	if err != nil {
		logger.Log.Error().Err(err).Msg("error while setup config")
	}
	logger.Log.Info().Msg("enter for main")
	//=========================== connection to db ===========================//
	// check database parameter validation

	if config.Config.DB_USER == "" || config.Config.DB_PASSWORD == "" || config.Config.DB_SERVER == "" || config.Config.DB_PORT == "" || config.Config.DB_DATABASE == "" {
		logger.Log.Error().Err(err).Msg("database env parameters are not found")
		return
	}
	var dsn string
	switch config.Config.DATABASE {
	case "postgres":
		logger.Log.Info().Msg("database connection : postgres")
		dsn = "host=" + config.Config.DB_SERVER + " port=" + config.Config.DB_PORT + " user=" + config.Config.DB_USER + " dbname=" + config.Config.DB_DATABASE + " password=" + config.Config.DB_PASSWORD + " sslmode=disable"
	case "mysql":
		logger.Log.Info().Msg("database connection : mysql")
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", config.Config.DB_USER, config.Config.DB_PASSWORD, config.Config.DB_SERVER, config.Config.DB_PORT, config.Config.DB_DATABASE)
	default:
		logger.Log.Error().Err(err).Msg("Invalid database slection")
		log.Error().Err(err).Msg("Invalid database slection")
		return
	}
	log.Info().Msg(dsn)
	db, err := gorm.Open(config.Config.DATABASE, dsn)
	if err != nil {
		fmt.Println("Error in db :", err.Error())
		logger.Log.Error().Err(err).Msg("Error in database connection")
		return
	}
	defer db.Close()

	db.LogMode(true)
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(10)

	logger.Log.Info().Msg("database connected successfully...")
	model.DBConn = db

	router := gin.Default()
	//router.L("F:/Dean-ai/users/templates/*")
	router.Use(func(c *gin.Context) {
		// router.LoadHTMLGlob("templates/*")
		router.LoadHTMLGlob("templates/*")
		c.Set("DB", db)
		c.Next()
	})
	router.Use(CORSMiddleware())

	port := config.Config.SERVICE_PORT
	if port == "" {
		logger.Log.Error().Err(err).Msg("service port not found")
		return
	}
	router.LoadHTMLGlob("templates/*")
	// router.LoadHTMLFiles("templates/*")
	routes := &services.HandlerService{}
	routes.Bootstrap(router)

	log.Info().Msg("Starting server on :" + port)
	logger.Log.Info().Msg("Starting server on :" + port)
	router.Run(":" + port)
}
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3006")
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Token")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}
