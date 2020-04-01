package helper

import (
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

func MatchUrl(s string) string {
	if s == "" {
		return ""
	}

	reg, err := regexp.Compile(`(([a-zA-Z0-9\._-]+\.[a-zA-Z]{2,6})|([0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}))(:[0-9]{1,4})*(/[a-zA-Z0-9\&%_\./-~-]*)?`)
	if err != nil {
		log.Error(err)
		return ""
	}
	s = reg.FindString(s)
	return s
}

func GenderUrl(u string) string {
	u = strings.TrimSpace(u)
	// i := strings.Index(u, "http://")
	// if i >= 0 {
	// 	return u
	// }

	if strings.HasPrefix(u, "http://") {
		return u
	}

	// i = strings.Index(u, "https://")
	// if i >= 0 {
	// 	return u
	// }

	if strings.HasPrefix(u, "https://") {
		return u
	}

	u = "http://" + u
	return u
}
