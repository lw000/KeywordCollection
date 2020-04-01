package service

import (
	"KeywordCollection/dao/table"
	"KeywordCollection/global"
	log "github.com/sirupsen/logrus"
)

// 域名数据库服务
type DomainDaoService struct {
}

func (domian *DomainDaoService) Query(status int) ([]table.TDomain, error) {
	return domian.selectFromDb(status)
}

func (domian *DomainDaoService) selectFromDb(status int) ([]table.TDomain, error) {
	rows, err := global.DBReptiledata.DB().Query(`SELECT id, name, domain, status FROM domain WHERE status=?;`, status)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rows.Close()

	var domains []table.TDomain
	for rows.Next() {
		var domain table.TDomain
		err = rows.Scan(&domain.Id, &domain.Name, &domain.Domain, &domain.Status)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		domains = append(domains, domain)
	}
	return domains, nil
}
