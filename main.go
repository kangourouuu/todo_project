// @title todo_project Management API
// @version 1.0
// @description API for managing todo_project
// @host localhost:9002
// @BasePath /
package main

import (

	"os"
	"path/filepath"


	api "todo_project/api"
	"todo_project/common/log"
	"todo_project/internal"
	"todo_project/internal/sqlclient"
	auth "todo_project/middleware"
	"todo_project/common/limiter"
	"todo_project/internal/redis"
	server "todo_project/server/http"

	"github.com/caarlos0/env/v10"
	//"github.com/go-chi/chi/v5"
	//"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"


	//swaggerFiles
	//swaggerFiles
	//"gorm.io/driver/postgres" 
	//_ "todo_project/docs"

)

type Config struct {
	Dir     string `env:"CONFIG_DIR" envDefault:"configs/config.json"`
	Port    string
	LogType string
	LogFile string
	DB      string
	Redis   string
}

var config Config
var redisClient redis.IRedis

func init() {
	// Parse environment variables
	if err := env.Parse(&config); err != nil {
		log.Fatal("Failed to parse environment variables: ", err)
	}

	// Initialize viper for configuration
	viper.SetConfigFile(config.Dir)
	viper.SetConfigType("json")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Failed to read config file: ", err)
	}

	// Set config values from viper
	config = Config{
		Dir:     config.Dir,
		Port:    viper.GetString("main.port"),
		LogType: viper.GetString("main.log_type"),
		LogFile: viper.GetString("main.log_file"),
		DB:      viper.GetString("main.db"),
		Redis:   viper.GetString("main.redis"),
	}

	// Initialize logger
	if config.LogType == "FILE" {
		if err := os.MkdirAll(filepath.Dir(config.LogFile), 0755); err != nil {
			log.Fatal("Failed to create log directory: ", err)
		}
		f, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal("Failed to open log file: ", err)
		}
		logrus.SetOutput(f)
	}

	// Initialize repository if DB is enabled
	if config.DB == "enabled" {
		gormsql := sqlclient.GormSqlConfig{
			Driver:   "postgresql",
			Host:     viper.GetString("db.host"),
			Port:     viper.GetInt("db.port"),
			Database: viper.GetString("db.database"),
			Username: viper.GetString("db.username"),
			Password: viper.GetString("db.password"),
		}
		initRepo(gormsql)
	}
	if config.Redis == "enabled" {
	var err error
	redisClient, err = redis.NewRedis(redis.RedisConfig{
		Addr:         viper.GetString("redis.address"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.db"),
		Poolsize:     15,
		PoolTimeOut:  15,
		IdleTimeOut:  15,
		ReadTimeOut:  15,
		WriteTimeOut: 15,
	})

	if err != nil {
		panic(err)
	}

	limiter.RateLimit = limiter.NewRateLimiter(viper.GetString(`redis.address`), viper.GetString(`redis.password`))
	}
}

func initRepo(gormSqlConfig sqlclient.GormSqlConfig) {
	internal.GormSqlClient = sqlclient.NewGormSqlClient(gormSqlConfig)
}

func main() {
	if config.DB != "enabled" {
		logrus.Fatal("Database is disabled, cannot start server")
	}

	logrus.Infof("DB name from config: %s", viper.GetString("db.database"))

	engine := server.NewEngine()

	engine.Use(auth.AuthMiddleWare())

	apiV2 := engine.Group("/api/v2")
	api.SetupRoutes(apiV2, redisClient)

	appServer := server.New(config.Port, engine)
	if err := appServer.Run(); err != nil {
		logrus.Fatalf("Server could not be started: %v", err)
	} 
}


