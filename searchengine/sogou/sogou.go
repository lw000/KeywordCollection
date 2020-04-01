package sogou

import (
	"KeywordCollection/searchengine/nethtp"
	tyutils "github.com/lw000/gocommon/utils"
	"time"
)

func getRealURL(device string, href string, fn func(body string) string) string {
	v, _, err := tyutils.TExecTime("getRealUrlAddress", func() (interface{}, error) {
		request := nethtp.GetRequest(device, nil)
		request.Timeout(time.Second * 3)
		resp, body, errs := request.Get(href).End()
		if len(errs) != 0 {
			return "", errs[0]
		}
		defer resp.Body.Close()

		if body != "" {
		}

		if fn != nil {
			href = fn(body)
		}
		return href, nil
	})
	if err != nil {
		return ""
	}
	return v.(string)
}
