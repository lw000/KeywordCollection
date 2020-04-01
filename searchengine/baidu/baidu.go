package baidu

import (
	"KeywordCollection/searchengine/nethtp"
	tyutils "github.com/lw000/gocommon/utils"
	"time"

	log "github.com/sirupsen/logrus"
)

func getRealURL(device string, href string) string {
	v, _, err := tyutils.TExecTime("getRealUrlAddress", func() (interface{}, error) {
		request := nethtp.GetRequest(device, nil)
		request.Timeout(time.Second * 3)
		resp, body, errs := request.Get(href).End()
		if len(errs) != 0 {
			log.Error(errs)
			return "", errs[0]
		}
		defer resp.Body.Close()

		if body != "" {
		}

		href = resp.Request.URL.String()

		return href, nil
	})
	if err != nil {
		return ""
	}
	return v.(string)
}
