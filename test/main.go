// test project main.go
package main

import (
	"KeywordCollection/models"
	"KeywordCollection/test/config"
	"KeywordCollection/test/ws"
	"encoding/json"
	tyhttp "github.com/lw000/gocommon/network/http"
	"log"
	"net/url"
	"time"
)

// var (
// 	urladdr string
// 	engine  string
// 	device  string
// 	wd      string
// 	page    int
// )
//
// func init() {
// 	flag.StringVar(&engine, "urladdr", "", "-urladdr")
// 	flag.StringVar(&engine, "engine", "", "-engine baidu/sogou/360/shenma")
// 	flag.StringVar(&device, "device", "pc", "-device pc/mobile")
// 	flag.StringVar(&wd, "wd", "", "-wd 搜索关键词")
// 	flag.IntVar(&page, "page", 1, "-page 页数")
// }

// HTTPReponseData ...
type HTTPReponseData struct {
	C int               `json:"c"`
	M string            `json:"m"`
	D map[string]string `json:"d"`
}

var (
	clientID string
)

func httpGet(url string) {
	_, data, er := tyhttp.DoHttpGet(url, nil, time.Second*time.Duration(60))
	if er != nil {
		log.Println(er)
		return
	}

	// log.Println(string(data))

	var d HTTPReponseData
	er = json.Unmarshal(data, &d)
	if er != nil {
		log.Println(er)
		return
	}
	log.Println(d)
}

func query(urladdr, clientID, id, engine, device, wd, domain, page string) {
	u := url.Values{}
	u.Add("clientID", clientID)
	u.Add("engine", engine)
	u.Add("device", device)
	u.Add("domain", domain)
	u.Add("wd", wd)
	u.Add("id", id)
	u.Add("page", page)
	httpGet(urladdr + "?" + u.Encode())
}

func main() {
	// if !flag.Parsed() {
	// 	flag.Parse()
	// }
	//
	// if wd == "" {
	// 	flag.PrintDefaults()
	// 	return
	// }

	cfg := config.NewConfig()
	er := cfg.Load("./conf/conf.json")
	if er != nil {
		log.Println(er)
		return
	}

	if cfg.Query.Urladdr == "" {
		log.Println(er)
		return
	}

	client := &ws.FastWsClient{}
	client.HandleConnected(func() {
		log.Printf("connected")
	})

	client.HandleDisConnected(func() {
		log.Println("disconnected")
	})

	client.HandleMessage(func(data []byte) {
		var cmd models.WSCMD
		er = cmd.Decode(data)
		if er != nil {
			log.Println(er)
			return
		}

		switch cmd.MainID {
		case 1:
			switch cmd.SubID {
			case 1:
				clientID = cmd.Data
				log.Println("clientID", clientID)

				for _, w := range cfg.Query.Data {
					go query(cfg.Query.Urladdr, clientID, w.ID, w.Engine, w.Device, w.Wd, w.Domain, w.Page)
					time.Sleep(time.Microsecond * time.Duration(cfg.Millisecond))
				}
			case 2:

			}

		case 2:
			switch cmd.SubID {
			case 1:
				log.Println(cmd.Data)
			case 2:

			}
		default:
			log.Println("未知消息")
		}

	})

	er = client.Open(cfg.Ws.Scheme, cfg.Ws.Host, cfg.Ws.Path)
	if er != nil {
		log.Println(er)
		return
	}

	go client.Run()

	select {}
}
