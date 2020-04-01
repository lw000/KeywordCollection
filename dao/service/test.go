package service

import (
	log "github.com/sirupsen/logrus"
)

func TestDao() {
	{
		serv := DomainDaoService{}
		v, err := serv.Query(1)
		if err != nil {
			log.Panic(err)
		}
		log.Printf("%+v", v)
	}

	{
		serv := KeywordsDaoService{}
		v, err := serv.Query(0, 10, 1)
		if err != nil {
			log.Panic(err)
		}
		log.Printf("%+v", v)
	}

	{
		serv := SearchEnginesDaoService{}
		v, err := serv.Query(1)
		if err != nil {
			log.Panic(err)
		}
		log.Printf("%+v", v)
	}

	// {
	// 	cur := time.Now()
	// 	serialNumber := fmt.Sprintf("%d%d%d%d%d%d%d",
	// 		cur.Year(),
	// 		cur.Month(),
	// 		cur.Day(),
	// 		cur.Hour(),
	// 		cur.Minute(),
	// 		cur.Second(),
	// 		tyIdWorker.IdworkerServ().NewId())
	// 	serv := SearchResultDaoService{}
	// 	err := serv.Insert(serialNumber, constant.EngineBaidu, constant.DevicePc, "斗地主", 1, "", time.Now(), nil)
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}
	// }
}
