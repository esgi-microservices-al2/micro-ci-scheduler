package model

import (
	"github.com/System-Glitch/goyave/v2/database"
	"github.com/jinzhu/gorm"
)

func init() {
	database.RegisterModel(&Job{})
}

type Job struct {
	gorm.Model
	CronExpression string `gorm:"size:255"`
	IdProject      int    // to see with project module
	Name           string `gorm:"size:255"`
}
