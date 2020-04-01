package parseserv

import (
	"KeywordCollection/global"
	"KeywordCollection/models"
	"KeywordCollection/searchengine"
	"KeywordCollection/server/dbserver"
	"github.com/olahol/melody"
	log "github.com/sirupsen/logrus"
)

func saveKeywordResultHandler(html *htmlParse) error {
	var (
		err error
		// rank     = 0
		// topThree string
	)

	// // 获取搜索前三名
	// topThree, err = html.Store().Top(3)
	// if err != nil {
	// 	log.Error(err)
	// }

	if html.Store().Length() > 0 {
		// 计算关键字排名
		html.Store().Foreach(func(index int, data *searchengine.StoreItem) bool {
			// var (
			// 	domainHost    *url.URL
			// 	hyperlinkHost *url.URL
			// )

			// s1 := helper.MatchUrl(html.Word().Domain)
			// s1 = helper.GenerUrl(s1)
			//
			// s2 := helper.MatchUrl(data.Hyperlink)
			// s2 = helper.GenerUrl(s2)
			//
			// domainHost, err = url.Parse(s1)
			// if err != nil {
			// 	log.Error(err)
			// 	return true
			// }
			//
			// hyperlinkHost, err = url.Parse(s2)
			// if err != nil {
			// 	log.Error(err)
			// 	return true
			// }
			//
			// if strings.Compare(hyperlinkHost.Host, domainHost.Host) == 0 {
			// 	rank = index + 1
			// 	return false
			// }

			job := &dbsrv.UpdateResultTask{
				Engine:       html.Word().Engine,
				Device:       html.Word().Device,
				KeywordId:    html.Word().KeywordId,
				Keyword:      html.Word().Keyword,
				SerialNumber: html.Word().SerialNumber,
				Ranks:        index + 1,
				Content:      data.JSON(),
			}

			if err = dbsrv.DBServer().AddJob(job); err != nil {
				log.Error(err)
			}

			return true
		})
	}

	// 通知客户端检索结果
	client, ok := global.QueryClients.Load(html.Word().ClientId)
	if !ok {
		return err
	}

	session := client.(*melody.Session)
	if session.IsClosed() {
		return err
	}

	if html.Store().Length() > 0 {
		html.Store().Foreach(func(index int, data *searchengine.StoreItem) bool {
			var buf []byte
			cmd := &models.WSCMD{MainID: 2, SubID: 1}
			ackQuery := &models.WSAckQuery{Code: 1, Id: html.Word().KeywordId, Rank: 0, Data: data.JSON()}
			buf, err = cmd.EncodeCmd(ackQuery)
			if err != nil {
				log.Error(err)
				return false
			}

			if err = session.Write(buf); err != nil {
				log.Error(err)
				return true
			}

			return true
		})
	} else {
		var buf []byte
		cmd := &models.WSCMD{MainID: 2, SubID: 1}
		ackQuery := &models.WSAckQuery{Code: 1, Id: html.Word().KeywordId, Data: "未检索到数据"}
		buf, err = cmd.EncodeCmd(ackQuery)
		if err != nil {
			log.Error(err)
			return err
		}

		if err = session.Write(buf); err != nil {
			log.Error(err)
		}
	}
	return nil
}
