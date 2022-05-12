package main

import (
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"github.com/tebeka/selenium/firefox"
	"strings"
)

// Simulator 模拟浏览器
type Simulator struct {
	driverPath string // 浏览器驱动路径
	port       int    // 浏览器驱动端口
	caps       selenium.Capabilities
}

func (s *Simulator) AddChromeCap(cap chrome.Capabilities) {
	s.caps.AddChrome(cap)
}

func (s *Simulator) NewChromeService(opts ...selenium.ServiceOption) (*selenium.Service, error) {
	return selenium.NewChromeDriverService(s.driverPath, s.port, opts...)
}

func (s *Simulator) AddFirefoxCap(cap firefox.Capabilities) {
	s.caps.AddFirefox(cap)
}

func (s *Simulator) NewFirefoxService(opts []selenium.ServiceOption) (*selenium.Service, error) {
	return selenium.NewGeckoDriverService(s.driverPath, s.port, opts...)
}

func (s *Simulator) NewWebDriver() (selenium.WebDriver, error) {
	return selenium.NewRemote(s.caps, fmt.Sprintf("http://localhost:%d/wd/hub", s.port))
}

// NewSimulator 创建模拟器
func NewSimulator(browserName, driverPath string, port int) *Simulator {
	return &Simulator{
		driverPath: driverPath,
		port:       port,
		caps:       selenium.Capabilities{"browserName": strings.ToLower(browserName)},
	}
}
