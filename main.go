package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"ty/car-prices-master/pkg/mongodb"

	"github.com/PuerkitoBio/goquery"
	"ty/car-prices-master/downloader"
	"ty/car-prices-master/model"
	"ty/car-prices-master/scheduler"
	"ty/car-prices-master/spiders"
)

var (
	//https://car.autohome.com.cn/2sc/hefei/list/
	//StartUrl = "/2sc/%s/a0_0msdgscncgpi1ltocsp1exb4/"
	//StartUrl =  "/2sc/110100/index.html"
	StartUrl = "/2sc/%s/list/"
	BaseUrl  = "https://car.autohome.com.cn"

	maxPage int = 99
	cars    []*spiders.QcCar
)

func Start(url string, ch chan []*spiders.QcCar) {
	body := downloader.Get(BaseUrl + url)
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Printf("Downloader.Get err: %v", err)
	}

	currentPage := spiders.GetCurrentPage(doc)
	nextPageUrl, _ := spiders.GetNextPageUrl(doc)

	if currentPage > 0 && currentPage <= maxPage {
		cars := spiders.GetCars(doc)
		log.Println(cars)
		ch <- cars
		if url := nextPageUrl; url != "" {
			scheduler.AppendUrl(url)
		}

		log.Println(url)
	} else {
		log.Println("Max page !!!")
	}
}

func main() {
	_, err := mongodb.New()
	if err != nil {
		log.Fatal("new mongo err", err)
	}

	citys := spiders.GetCitys()
	for _, v := range citys {
		scheduler.AppendUrl(fmt.Sprintf(StartUrl, v.Pinyin))
	}

	start := time.Now()
	delayTime := time.Second * 6000

	ctx := context.Background()
	ch := make(chan []*spiders.QcCar)

L:
	for {
		if url := scheduler.PopUrl(); url != "" {
			go Start(url, ch)
		}

		select {
		case r := <-ch:
			cars = append(cars, r...)
			go Start(scheduler.PopUrl(), ch)
		case <-time.After(delayTime):
			log.Println("Timeout...")
			break L
		}
	}

	if len(cars) > 0 {
		model.AddCars(ctx, cars)
	}

	log.Printf("Time: %s", time.Since(start)-delayTime)
}
