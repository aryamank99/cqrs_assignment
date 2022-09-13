package database

import (
	"fmt"
	"sync"
	"tenantMicroserviceAssignment/cmd"
	"tenantMicroserviceAssignment/models"
)
import "gorm.io/gorm"
import "gorm.io/driver/mysql"

// lock mutex
var lock = &sync.Mutex{}

type DbConn struct {
	commandsConn *gorm.DB
	storeConn    *gorm.DB
}

var instance *DbConn

func Connect(storeDbConfig cmd.StoreDatabaseConfig, commandsDbConfig cmd.CommandDatabaseConfig) *DbConn {
	lock.Lock()
	defer lock.Unlock()
	if instance == nil {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", storeDbConfig.DatabaseUsername, storeDbConfig.DatabasePassword, storeDbConfig.DatabaseConnectionUrl, storeDbConfig.DatabaseConnectionPort, storeDbConfig.DatabaseName)
		storeDb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect store database")
		}
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", commandsDbConfig.DatabaseUsername, commandsDbConfig.DatabasePassword, commandsDbConfig.DatabaseConnectionUrl, commandsDbConfig.DatabaseConnectionPort, commandsDbConfig.DatabaseName)
		commandsDb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		instance = &DbConn{
			commandsConn: commandsDb,
			storeConn:    storeDb,
		}
		performMigrations(instance.commandsConn, instance.storeConn)
	}

	return instance
}

func GetStoreDbConn() *gorm.DB {
	if instance.storeConn == nil {
		panic("db connection not established")
	}
	return instance.storeConn
}

func GetCommandsDbConn() *gorm.DB {
	if instance.commandsConn == nil {
		panic("db connection not established")
	}
	return instance.commandsConn
}

func performMigrations(commandsDb *gorm.DB, storeDb *gorm.DB) {
	// Migrate the schema
	err := storeDb.AutoMigrate(
		&models.Tenant{},
		&models.LoginCredentials{},
	)
	if err != nil {
		panic("unable to perform store migrations...")
	}
	err = commandsDb.AutoMigrate(
		&models.CreateTenantRequest{},
		&models.UpdateTenantCredentialsRequest{},
		&models.UpdateTenantLicenseRequest{},
		&models.UpdateTenantStatusRequest{},
	)
	if err != nil {
		panic("unable to perform commands migrations...")
	}

}
