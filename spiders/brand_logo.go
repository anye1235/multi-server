package spiders

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strconv"
	"ty/car-prices-master/downloader"
)

const Url = "https://m.che168.com/beijing/list/"

func GetBrandLogo() (cars []QcCar) {
	body := downloader.Get(Url)
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Printf("Downloader.Get err: %v", err)
	}

	ctx := context.Background()
	//cityName := GetCityName(doc)
	//i := doc.Find(".replace.block-brand.block-brand-all .replace-main .list-line ").Get(0)
	//fmt.Printf("doc:%v", i)

	doc.Find(".replace.block-brand.block-brand-all .replace-main .list-line ").Each(func(i int, selection *goquery.Selection) {
		//doc.Find("replace block-brand block-brand-all replace-main ul list-line li:not(.line)").Each(func(i int, selection *goquery.Selection) {
		selection.Find("li").Each(func(ii int, selectionSub *goquery.Selection) {
			brandid, _ := selectionSub.Attr("data-brandid")
			//title, _ := selectionSub.Find("img").Attr("title")
			src, _ := selectionSub.Find("img").Attr("data-src")
			if len(brandid) == 0 || len(src) == 0 {
				return
			}

			brandidInt, err := strconv.ParseInt(brandid, 10, 32)
			if nil != err {
				return
			}
			Update(ctx, int(brandidInt), src)
		})

	})

	return cars
}
