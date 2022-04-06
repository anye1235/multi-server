package main

import (
	"context"
	"fmt"
	"io/ioutil"
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

	maxPage = 1000
	cars    []*spiders.QcCar
)

func Start(url string, ch chan []*spiders.QcCar, loopCount int) {
	time.Sleep(time.Second * 1)
	body := downloader.Get(BaseUrl + url)
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Printf("Downloader.Get err: %v", err)
	}

	currentPage := spiders.GetCurrentPage(doc)
	if 0 == currentPage && loopCount < 10 {
		loopCount++
		Start(url, ch, loopCount)
		return
	}
	nextPageUrl, _ := spiders.GetNextPageUrl(doc)

	if currentPage > 0 && currentPage <= maxPage {
		cars := spiders.GetCars(doc)
		log.Printf("input cars numbers: %v", len(cars))
		ch <- cars
		if url := nextPageUrl; url != "" {
			scheduler.AppendUrl(url)
		}

		log.Println(url)
	} else {
		log.Printf("Max page !!! curr: %v, maxPage: %v, url:%s", currentPage, maxPage, url)
		bodyByte, _ := ioutil.ReadAll(body)
		log.Printf("%s", bodyByte)
	}
}

func main() {
	_, err := mongodb.New()
	if err != nil {
		log.Fatal("new mongo err", err)
	}

	citys := spiders.GetCitys()
	log.Printf("total citys: %v", len(citys))
	for _, v := range citys {
		go scheduler.AppendUrl(fmt.Sprintf(StartUrl, v.Pinyin))
	}

	log.Printf("total city urls : %v", len(scheduler.URLs))

	start := time.Now()
	delayTime := time.Second * 36000

	ctx := context.Background()
	ch := make(chan []*spiders.QcCar, 1000)
	totalPop := 0
	totalCreate := 0

L:
	for {
		if url := scheduler.PopUrl(); url != "" {
			totalPop++
			go Start(url, ch, 0)
		}

		log.Printf("total cars len: %v", len(ch))

		select {
		case r := <-ch:
			//cars = append(cars, r...)
			log.Printf("current carlist len: %v", len(cars))
			model.AddCarsOpt(ctx, r...)
			totalCreate++
			go Start(scheduler.PopUrl(), ch, 0)
		case <-time.After(delayTime):
			log.Println("channel Timeout...")
			break L
		}
	}

	//if len(cars) > 0 {
	//	model.AddCars(ctx, cars)
	//}

	log.Printf("Cost Time: %v", time.Since(start).Seconds())
}
