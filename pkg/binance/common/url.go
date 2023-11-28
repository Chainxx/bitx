package common

import "net/url"

const (
	endpoint0 = "api.binance.com"
	endpoint1 = "api1.binance.com"
	endpoint2 = "api2.binance.com"
	endpoint3 = "api3.binance.com"
)

func Url(path, param string) string {
	u := url.URL{}
	u.Scheme = "https"
	u.Host = endpoint0
	u.Path = path
	u.RawQuery = param
	
	return u.String()
}
