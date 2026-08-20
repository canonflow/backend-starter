package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/canonflow/backend-starter/internal/config"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type DBProperty struct {
	Driver                string
	Host                  string
	Port                  string
	Username              string
	Password              string
	Database              string
	IdleConnection        int
	MaxConnection         int
	MaxLifetimeConnection int
}

func GetDriverConnection(property DBProperty) string {
	switch strings.ToLower(property.Driver) {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			property.Username,
			property.Password,
			property.Host,
			property.Port,
			property.Database,
		)
	default:
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			property.Username,
			property.Password,
			property.Host,
			property.Port,
			property.Database,
		)
	}
}

func getDatabaseDriver(driver string) func(string) gorm.Dialector {
	switch strings.ToLower(driver) {
	case "mysql":
		return mysql.Open
	default:
		return mysql.Open
	}
}

func NewDatabase(driver string, log *zap.Logger) *gorm.DB {
	username := config.Get[string](config.DBUsername)
	password := config.Get[string](config.DBPassword)
	host := config.Get[string](config.DBHost)
	port := config.Get[string](config.DBPort)

	database := config.Get[string](config.DBName)
	idleConnection := config.Get[int](config.DBIdle)
	maxConnection := config.Get[int](config.DBMax)
	maxLifetimeConnection := config.Get[int](config.DBLifeTime)

	dsn := GetDriverConnection(DBProperty{
		Driver:                driver,
		Username:              username,
		Password:              password,
		Database:              database,
		Host:                  host,
		Port:                  port,
		IdleConnection:        idleConnection,
		MaxConnection:         maxConnection,
		MaxLifetimeConnection: maxLifetimeConnection,
	})
	databaseDriver := getDatabaseDriver(driver)

	db, err := gorm.Open(databaseDriver(dsn), &gorm.Config{
		Logger: gormLogger.New(&zapWriter{
			Logger: log,
		}, gormLogger.Config{
			SlowThreshold:             time.Second * 5,
			Colorful:                  false,
			IgnoreRecordNotFoundError: false,
			ParameterizedQueries:      true,
			LogLevel:                  gormLogger.Info,
		}),
	})
	if err != nil {
		log.Fatal("Failed to connect with database",
			zap.String("error", err.Error()),
		)
	}

	connection, err := db.DB()
	if err != nil {
		log.Fatal("Failed to connect with database",
			zap.String("error", err.Error()),
		)
	}

	connection.SetMaxIdleConns(idleConnection)
	connection.SetMaxOpenConns(maxConnection)
	connection.SetConnMaxLifetime(time.Second * time.Duration(maxLifetimeConnection))

	return db
}

type zapWriter struct {
	Logger *zap.Logger
}

func (l *zapWriter) Printf(message string, args ...interface{}) {
	l.Logger.Debug(fmt.Sprintf("[GORM TRACE] "+message, args))
}
