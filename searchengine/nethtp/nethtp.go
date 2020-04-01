package nethtp

import (
	"KeywordCollection/constant"
	"fmt"
	tyutils "github.com/lw000/gocommon/utils"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/parnurzeal/gorequest"
	log "github.com/sirupsen/logrus"
)

// GetRequest ...
func GetRequest(device string, cookies map[string]string) *gorequest.SuperAgent {
	request := gorequest.New()
	switch device {
	case constant.DevicePc:
		request.AppendHeader("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
		request.AppendHeader("Accept-Encoding", "gzip, deflate, br")
		request.AppendHeader("Accept-Language", "zh-CN,zh;q=0.9")
		request.AppendHeader("Connection", "keep-alive")
		request.AppendHeader("Cache-Control", "max-age=0")
		request.AppendHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.131 Safari/537.36")
	case constant.DeviceMobile:
		request.AppendHeader("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
		request.AppendHeader("Accept-Encoding", "gzip, deflate, br")
		request.AppendHeader("Accept-Language", "zh-CN,zh;q=0.9")
		request.AppendHeader("Connection", "keep-alive")
		request.AppendHeader("Cache-Control", "max-age=0")
		request.AppendHeader("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 DeviceMobile/15A372 Safari/604.1")
	default:
		return nil
	}

	for k, v := range cookies {
		request.AddCookie(&http.Cookie{Name: k, Value: v})
	}

	return request
}

func doGetRequest(device string, strurl string, cookies map[string]string) (string, error) {
	v, _, err := tyutils.TExecTime("doGetRequest", func() (interface{}, error) {
		request := GetRequest(device, cookies)
		request.Timeout(time.Second * 5)

		resp, body, errs := request.Get(strurl).End()
		// resp, body, errs := request.Get(requestUrl).Retry(2, time.Second*5, http.StatusBadRequest, http.StatusInternalServerError).End(func(response gorequest.Response, body string, errs []error) {
		// 	log.Println(body)
		// })

		if len(errs) != 0 {
			log.Error(errs)
			return "", errs[0]
		}
		defer resp.Body.Close()

		if body == "" {
			log.Error("body is empty")
			return "", errs[0]
		}

		if resp == nil {
			log.Error("resp is nil")
			return "", errs[0]
		}

		return body, nil
	})
	if err != nil {
		return "", err
	}

	return v.(string), nil
}

// FirstPageRequest ...
func FirstPageRequest(device string, host string, path string, u *url.Values, cookies map[string]string) (string, error) {
	if u != nil {
		requestURL := fmt.Sprintf("%s%s?%s", host, path, u.Encode())
		log.Info(requestURL)
		return doGetRequest(device, requestURL, cookies)
	}

	requestURL := fmt.Sprintf("%s%s?%s", host, path, "")
	return doGetRequest(device, requestURL, cookies)
}

// NextPageRequest ...
func NextPageRequest(device string, host string, u *url.URL, cookies map[string]string) (string, error) {
	if host == "" {
		requestURL := u.String()
		log.Info(requestURL)
		return doGetRequest(device, requestURL, cookies)
	}

	requestURL := fmt.Sprintf("%s/%s", host, u.String())
	return doGetRequest(device, requestURL, cookies)
}

// SaveHTML ...
func SaveHTML(dir string, wd string, device string, page int, body string) {
	// TODO: 不保存网页
	// return

	go func() {
		ok, err := tyutils.PathExists(dir)
		if err != nil {
			log.Error(err)
			return
		}
		if !ok {
			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				log.Error(err)
				return
			}
		}

		f, err := os.Create(fmt.Sprintf("%s/%s_%s_%d.html", dir, wd, device, page))
		if err != nil {
			log.Error(err)
			return
		}

		n, err := f.Write([]byte(body))
		if n > 0 {

		}
	}()
}
