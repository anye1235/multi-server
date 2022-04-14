package spiders

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"ty/car-prices-master/pkg/mongodb"
)

type City struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	Name     string             `bson:"name" json:"name"`
	Pinyin   string             `bson:"pinyin" json:"pinyin"`
	CiTyId   int                `bson:"id" json:"id"`
	Citysite int                `bson:"citysite" json:"citysite"` // 不知道啥用
	Pid      int                `bson:"pid" json:"pid"`           // 不知道啥用
}

// 中国
const PID_0 = 0

// 省
const PID_1 = 1

// 市（一级）
const PID_2 = 2

// 县
const PID_3 = 3

func AddAllCity() {
	var citys []City
	json.Unmarshal([]byte(cityArray), &citys)

	ctx := context.Background()
	for _, city := range citys {
		if city.CiTyId%10000 == 0 {
			city.Pid = PID_1
		} else if city.CiTyId%100 == 0 {
			city.Pid = PID_2
		} else {
			city.Pid = PID_3
		}

		AddCity(ctx, &city)
	}

}

// 增加品牌
func AddCity(ctx context.Context, city *City) {
	if _, err := mongodb.GetMongoClient().Insert(ctx, "", "tbl_city", city); err != nil {
		log.Printf("tbl_brand Create index: %v, err : %v", city, err)
	}
}
