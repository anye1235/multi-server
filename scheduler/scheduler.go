package scheduler

//var URLs = []string{}
//var mutex sync.Mutex
const MaxExistNum = 3 //最大进入次数

var URLs = make(chan string, 1000)
var existUrlMap = make(map[string]int)

func PopUrl() string {
	//mutex.Lock()
	//defer mutex.Unlock()

	//length := len(URLs)
	//if length < 1 {
	//	return ""
	//}
	//url := URLs[0]
	//URLs = URLs[1:]
	//return url

	//delayTime := time.Millisecond * 100
	select {
	case r := <-URLs:
		return r
	default:
		return ""
	}

}

func AppendUrl(url string) {
	URLs <- url
	countExistMapOldNum(url)
}

func countExistMapOldNum(url string) {
	if 0 == len(url) {
		return
	}
	oldNum := existUrlMap[url]
	if oldNum > MaxExistNum {
		return
	}
	existUrlMap[url] = oldNum + 1
}
