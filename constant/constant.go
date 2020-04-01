package constant

const (
	// EngineBaidu ...
	EngineBaidu = "baidu"
	// EngineSogou ...
	EngineSogou = "sogou"
	// EngineShenma ...
	EngineShenma = "shenma"
	// EngineSo360 ...
	EngineSo360 = "so360"
)

const (
	// DevicePc ...
	DevicePc = "pc"
	// DeviceMobile ...
	DeviceMobile = "mobile"
)

const (
	// ChromedriverPath ...
	ChromedriverPath = "./chromedriver/chromedriver"
)

const (
	// DefaultDir ...
	DefaultDir = "./html"
	// BaiduDir ...
	BaiduDir = "./html/baidu"
	// SougouDir ...
	SougouDir = "./html/sogou"
	// So360Dir ...
	So360Dir = "./html/so360"
	// ShenmaDir ...
	ShenmaDir = "./html/shenma"
)

const (
	KeywordStatusDisable  int = 0 /*iota*/
	KeywordStatusDefault  int = 1 // 等待检索
	KeywordStatusChecking int = 2 // 检索中
	KeywordStatusFail     int = 3 // 检索失败
	KeywordStatusOk       int = 4 // 检索完成
)

const (
	EnginesStatusDisable int = 0 /*iota*/
	EnginesStatusEnable  int = 1 // 启用
)

const (
	DomainStatusDisable int = 0 /*iota*/
	DomainStatusEnable  int = 1 // 启用
)

func GetStoreDir(engine string) string {
	switch engine {
	case EngineBaidu:
		return BaiduDir
	case EngineSo360:
		return So360Dir
	case EngineSogou:
		return SougouDir
	case EngineShenma:
		return ShenmaDir
	default:
	}
	return DefaultDir
}
