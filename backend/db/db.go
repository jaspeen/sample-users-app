package db

import (
	"github.com/jaspeen/sample-users-app/config"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

func InitDB(dbConnectString string, logMode bool) (*gorm.DB, error) {
	dataSourceName := "postgres://" + dbConnectString
	db, err := gorm.Open("postgres", dataSourceName)

	if err != nil {
		return nil, err
	}

	db.LogMode(logMode)

	db.AutoMigrate(&User{}, &RefreshToken{})

	// create admin user if not exist
	var adminUser User
	dbRes := db.Take(&adminUser, "email = ?", config.C.AdminEmail)
	if err == nil && dbRes.RowsAffected == 0 {
		passHash, err := bcrypt.GenerateFromPassword([]byte(config.C.AdminPass), 14)
		if err != nil {
			return nil, err
		}
		name := "Admin"
		err = db.Save(&User{
			FirstName: &name,
			LastName:  &name,
			Email:     config.C.AdminEmail,
			Admin:     true,
			Password:  string(passHash),
		}).Error
		return db, err
	}
	return db, nil
}
