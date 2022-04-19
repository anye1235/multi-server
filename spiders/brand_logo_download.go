package spiders

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func DownLoadLogo() {
	list := GetAll(context.Background())
	if nil == list {
		return
	}
	fileStr, _ := os.Getwd()
	for _, brand := range list {
		if len(brand.Pinyin) != 0 && len(brand.LogoUrl) == 0 {
			continue
		}
		SaveBrandLogo(brand.LogoUrl, brand.Pinyin, fileStr)
	}

}

func SaveBrandLogo(url, name, fileStr string) (n int64, err error) {
	fmt.Println(name)
	name = fileStr + `\logo\` + name + ".png"
	out, err := os.Create(name)
	defer out.Close()
	resp, err := http.Get(url)
	defer resp.Body.Close()
	pix, err := ioutil.ReadAll(resp.Body)
	n, err = io.Copy(out, bytes.NewReader(pix))
	return
}
