package model

import (
	"context"
	"fmt"
	"log"
	"ty/car-prices-master/pkg/mongodb"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"ty/car-prices-master/spiders"
)

var (
	DB *gorm.DB

	host     string = "127.0.0.1"
	username string = "root"
	password string = "123456"
	dbName   string = "spiders"
)

func init() {
	var err error
	DB, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", username, password, dbName))
	if err != nil {
		log.Fatalf(" gorm.Open.err: %v", err)
	}

	DB.SingularTable(true)
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "sp_" + defaultTableName
	}
}

func AddCars(ctx context.Context, cars []*spiders.QcCar) {
	for index, car := range cars {
		if _, err := mongodb.GetMongoClient().Insert(ctx, "", "tbl_car_price", car); err != nil {
			log.Printf("db.Create index: %v, err : %v", index, err)
		}
	}
}

func AddCarsOpt(ctx context.Context, cars ...*spiders.QcCar) {
	for index, car := range cars {
		if _, err := mongodb.GetMongoClient().Insert(ctx, "", "tbl_car_price_2", car); err != nil {
			log.Printf("db.Create index: %v, err : %v", index, err)
		}
	}
}
