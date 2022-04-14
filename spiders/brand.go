package spiders

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"ty/car-prices-master/pkg/mongodb"
	"ty/car-prices-master/utils/httpclient"
)

//接口
type ResultHandle struct {
}

// HandleResult对接不同厂商，返回的格式不一样 这里需要制定返回
func (p *ResultHandle) HandleResult(body []byte, resObj interface{}) (interface{}, error) {
	if err := json.Unmarshal(body, resObj); err != nil {
		return nil, err
	}
	return nil, nil
}

const SPIDERS_BRAND_URL = "https://car.autohome.com.cn/2sc/loadbrand.ashx?area=guangzhou&brand=&ls=&spec=0&minPrice=0&maxPrice=0&minRegisteAge=0&maxRegisteAge=0&MileageId=0&disp=0&stru=0&gb=0&color=0&source=0&listview=0&sell=1&newCar=0&credit=0&sort=0&kw=&ex=c0d0t0p0w0r0u0e0s0a0o0i0b0"

type Brand struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Brandid   int                `bson:"brand_id" json:"brandid"`
	Brandname string             `bson:"brand_name" json:"brandname"`
	Letter    string             `bson:"letter" json:"letter"`
	Pinyin    string             `bson:"pinyin" json:"pinyin"`
	URL       string             `bson:"url" json:"url"`
}

func AddAllBrands() {
	ctx := context.Background()
	var res = make(map[string][]*Brand, 26)
	if _, err := httpclient.DoGet(SPIDERS_BRAND_URL, nil, ctx, &res, new(ResultHandle)); nil != err {
		log.Printf("request error %v", err)
	}

	if nil == res {
		return
	}

	for _, v := range res {
		AddBrand(ctx, v)
	}
}

// 增加品牌
func AddBrand(ctx context.Context, brands []*Brand) {
	for index, brand := range brands {
		if _, err := mongodb.GetMongoClient().Insert(ctx, "", "tbl_brand", brand); err != nil {
			log.Printf("tbl_brand Create index: %v, err : %v", index, err)
		}
	}
}
