package downloader

import (
	"bytes"
	"encoding/binary"
	"github.com/axgle/mahonia"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"

	"ty/car-prices-master/fake"
)

func Get(url string) io.Reader {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("http.NewRequest err", "%v", err)
	}

	req.Header.Add("User-Agent", fake.GetUserAgent())
	req.Header.Add("Referer", "https://car.autohome.com.cn")

	ip := RandIp()
	req.Header.Add("X-Forwarded-For", ip)
	req.Header.Add("Client-Ip", ip)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("client.Do err", "%v", err)
	}

	mah := mahonia.NewDecoder("gbk")
	return mah.NewReader(resp.Body)
}

//ip到数字
func ip2Long(ip string) uint32 {
	var long uint32
	binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &long)
	return long
}

//数字到IP
func backtoIP4(ipInt int64) string {
	// need to do two bit shifting and “0xff” masking
	b0 := strconv.FormatInt((ipInt>>24)&0xff, 10)
	b1 := strconv.FormatInt((ipInt>>16)&0xff, 10)
	b2 := strconv.FormatInt((ipInt>>8)&0xff, 10)
	b3 := strconv.FormatInt((ipInt & 0xff), 10)
	return b0 + "." + b1 + "." + b2 + "." + b3
}

func RandIp() string {
	// or if you prefer the super fast way
	ip1 := ip2Long("221.177.0.0")
	ip2 := ip2Long("221.177.7.255")
	//ip1 := ip2Long("192.168.0.0")
	//ip2 := ip2Long("192.168.0.255")
	x := ip2 - ip1
	return backtoIP4(int64(ip1) + int64(x))
}
