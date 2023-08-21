package model

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {

        database, err := gorm.Open(sqlite.Open("wallet.db"), &gorm.Config{})

        if err != nil {
                panic("Failed to connect to database!")
        }

        err = database.AutoMigrate(&Wallet{})
        if err != nil {
                return
        }
        err = database.AutoMigrate(&Transaction{})
        if err != nil {
                return
        }

        DB = database
}