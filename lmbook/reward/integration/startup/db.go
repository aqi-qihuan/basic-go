package startup

import (
	"basic-go/lmbook/reward/repository/dao"
	"context"
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

var db *sql.DB
