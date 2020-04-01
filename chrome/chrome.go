package chrome

import (
	"KeywordCollection/constant"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

// ChromeDriver ...
type ChromeDriver struct {
	rw      sync.RWMutex
	service *selenium.Service
	caps    selenium.Capabilities
	port    int
	status  int
}

func (c *ChromeDriver) Status() int {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.status
}

func (c *ChromeDriver) SetStatus(status int) {
	c.rw.Lock()
	defer c.rw.RUnlock()
	c.status = status
}

// ChromeDriverManager ...
type ChromeDriverManager struct {
}

// NewPcService ...
func (c *ChromeDriver) NewPcService(port int) error {
	c.port = port

	opts := []selenium.ServiceOption{}
	// opts := []selenium.ServiceOption{
	//    selenium.StartFrameBuffer(),           // Start an X frame buffer for the browser to run in.
	//    selenium.GeckoDriver(geckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
	// }

	// selenium.SetDebug(true)
	var err error
	c.service, err = selenium.NewChromeDriverService(constant.ChromedriverPath, port, opts...)
	if nil != err {
		log.Error(err)
		return err
	}

	// 链接本地的浏览器 chrome
	c.caps = selenium.Capabilities{
		"browserName": "chrome",
	}

	// 禁止图片加载，加快渲染速度
	// imagCaps := map[string]interface{}{
	// 	"profile.managed_default_content_settings.images": 2,
	// }

	chromeCaps := chrome.Capabilities{
		// Prefs: imagCaps,
		Path: "",
		Args: []string{
			"--headless", // 设置Chrome无头模式，在linux下运行，需要设置这个参数，否则会报错
			"--no-sandbox",
			"--disable-gpu",
			"--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36", // 模拟user-agent，防反爬
		},
	}

	c.caps.AddChrome(chromeCaps)

	return nil
}

// NewMobileService ...
func (c *ChromeDriver) NewMobileService(port int) error {
	c.port = port

	opts := []selenium.ServiceOption{}
	// opts := []selenium.ServiceOption{
	//    selenium.StartFrameBuffer(),           // Start an X frame buffer for the browser to run in.
	//    selenium.GeckoDriver(geckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
	// }

	// selenium.SetDebug(true)

	var err error
	c.service, err = selenium.NewChromeDriverService(constant.ChromedriverPath, port, opts...)
	if nil != err {
		log.Error(err)
		return err
	}

	// 链接本地的浏览器 chrome
	c.caps = selenium.Capabilities{
		"browserName": "chrome",
		"deviceName":  "Apple iPhone 6/7/8 Plus",
	}

	// 禁止图片加载，加快渲染速度
	// imagCaps := map[string]interface{}{
	// 	"profile.managed_default_content_settings.images": 2,
	// }
	chromeCaps := chrome.Capabilities{
		// Prefs: imagCaps,
		Path: "",
		Args: []string{
			"--headless", // 设置Chrome无头模式，在linux下运行，需要设置这个参数，否则会报错
			"--no-sandbox",
			"--disable-gpu",
			"--user-agent=Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 DeviceMobile/15A372 Safari/604.1",
		},
	}
	c.caps.AddChrome(chromeCaps)

	return nil
}

// OpenWebBrowser ...
func (c *ChromeDriver) OpenWebBrowser() (selenium.WebDriver, error) {
	webBrowser, err := selenium.NewRemote(c.caps, fmt.Sprintf("http://localhost:%d/wd/hub", c.port))
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return webBrowser, nil
}

// Stop ...
func (c *ChromeDriver) Stop() {
	if c == nil {
		return
	}
	if c.service == nil {
		return
	}
	err := c.service.Stop()
	if err != nil {
	}
}

// Start ...
func (c *ChromeDriverManager) Start() error {
	return nil
}

// Stop ...
func (c *ChromeDriverManager) Stop() {

}
