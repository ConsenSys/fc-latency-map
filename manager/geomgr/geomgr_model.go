package geomgr

type GeoIP struct {
	GeopluginRequest   string `json:"geoplugin_request"`
	GeopluginStatus    int    `json:"geoplugin_status"`
	GeopluginDelay     string `json:"geoplugin_delay"`
	GeopluginCredit    string `json:"geoplugin_credit"`
	GeopluginCity      string `json:"geoplugin_city"`
	GeopluginRegion    string `json:"geoplugin_region"`
	GeopluginLatitude  string `json:"geoplugin_latitude"`
	GeopluginLongitude string `json:"geoplugin_longitude"`
	GeopluginTimezone  string `json:"geoplugin_timezone"`
}

type Location struct {
	GeopluginCountryCode string `json:"geoplugin_countryCode"` //nolint:tagliatelle
	GeopluginLatitude    string `json:"geoplugin_latitude"`
	GeopluginLongitude   string `json:"geoplugin_longitude"`
}
