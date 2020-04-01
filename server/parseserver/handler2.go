package parseserv

import (
	"KeywordCollection/constant"
	"KeywordCollection/dao/service"
	"KeywordCollection/dao/table"
	"KeywordCollection/helper"
	"KeywordCollection/searchengine"
	"KeywordCollection/server/dbserver"
	log "github.com/sirupsen/logrus"
	"net/url"
	"strings"
)

// 解析网页数据服务
func saveDomainRankResultHandler(html *htmlParse) error {
	var (
		err     error
		domains []table.TDomain
	)

	if html.Store().Length() > 0 {
		serv := service.DomainDaoService{}
		domains, err = serv.Query(constant.DomainStatusEnable)
		if err != nil {
			log.Error(err)
			return err
		}

		for _, d := range domains {
			var (
				rank int
				job  *dbsrv.UpdateDomainRankTask
			)

			job = &dbsrv.UpdateDomainRankTask{
				DomainId:     d.Id,
				Engine:       html.Word().Engine,
				Device:       html.Word().Device,
				KeywordsId:   html.Word().KeywordId,
				SerialNumber: html.Word().SerialNumber,
			}

			html.Store().Foreach(func(index int, data *searchengine.StoreItem) bool {
				var (
					domainHost    *url.URL
					hyperlinkHost *url.URL
				)
				s1 := helper.MatchUrl(d.Domain)
				s1 = helper.GenderUrl(s1)

				s2 := helper.MatchUrl(data.Hyperlink)
				s2 = helper.GenderUrl(s2)
				domainHost, err = url.Parse(s1)
				if err != nil {
					log.Error(err)
					return true
				}

				hyperlinkHost, err = url.Parse(s2)
				if err != nil {
					log.Error(err)
					return true
				}

				if strings.Compare(hyperlinkHost.Host, domainHost.Host) == 0 {
					rank = index + 1
					job.Content = data.JSON()
					return false
				}

				return true
			})

			if rank > 0 {
				job.Rank = rank
				if err = dbsrv.DBServer().AddJob(job); err != nil {
					log.Error(err)
					return err
				}
			}
		}
	} else {
		job := &dbsrv.UpdateKeywordsStatusTask{
			KeywordId: html.Word().KeywordId,
			Engine:    html.Word().Engine,
			Device:    html.Word().Device,
			Status:    constant.KeywordStatusOk,
		}
		if err = dbsrv.DBServer().AddJob(job); err != nil {
			log.Error(err)
			return err
		}
	}

	return nil
}
