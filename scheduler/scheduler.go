package scheduler

//var URLs = []string{}
//var mutex sync.Mutex

var URLs = make(chan string, 1000)

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
}
