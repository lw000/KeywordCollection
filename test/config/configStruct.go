package config

type Config struct {
	Count       int64 `json:"count"`
	Millisecond int64 `json:"millisecond"`
	Query       struct {
		Data []struct {
			Device string `json:"device"`
			Domain string `json:"domain"`
			Engine string `json:"engine"`
			ID     string `json:"id"`
			Page   string `json:"page"`
			Wd     string `json:"wd"`
		} `json:"data"`
		Urladdr string `json:"urladdr"`
	} `json:"query"`
	Ws struct {
		Host   string `json:"host"`
		Path   string `json:"path"`
		Scheme string `json:"scheme"`
	} `json:"ws"`
}
