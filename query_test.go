package panodata

import (
	"math"
	"net/url"
	"testing"
)

func approxLatLonDist(lat1, lon1 float64, lat2, lon2 float64) float64 {
	dLat, dLon := lat2-lat1, lon2-lon1
	x, y := dLon*math.Cos(lat1), dLat
	return rEarth * math.Sqrt(x*x+y*y)
}

func TestQueryRadius(t *testing.T) {
	var lat, lon float64 = 45, 0
	var radius float64 = 10

	query := QueryRadius(lat, lon, radius)

	if d := approxLatLonDist(lat, lon, query.MinLat, query.MinLon); d > radius*math.Sqrt(2) {
		t.Error("dist (min) =", d, "want", radius*math.Sqrt(2))
	}

	if d := approxLatLonDist(lat, lon, query.MaxLat, query.MaxLon); d > radius*math.Sqrt(2) {
		t.Error("dist (max) =", d, "want", radius*math.Sqrt(2))
	}
}

func TestQueryURLDefaults(t *testing.T) {
	// Test defaults
	queryDefault := Query{}

	urlDefault, err := url.Parse(queryDefault.URL())
	if err != nil {
		t.Fatal("Malformed URL: ", err)
	}
	if urlDefault.Scheme != "http" {
		t.Error("Scheme =", urlDefault.Scheme, "want http")
	}
	if urlDefault.Host != "www.panoramio.com" {
		t.Error("Host =", urlDefault.Host, "want www.panoramio.com")
	}
	if urlDefault.Path != "/map/get_panoramas.php" {
		t.Error("Path =", urlDefault.Path, "want /map/get_panoramas.php")
	}

	params := urlDefault.Query()
	if set := params.Get("set"); set != "full" {
		t.Error("set =", set, "want full")
	}
	if size := params.Get("size"); size != "original" {
		t.Error("size =", size, "want original")
	}
	if minx := params.Get("minx"); minx != "0" {
		t.Error("minx =", minx, "want 0")
	}
	if miny := params.Get("miny"); miny != "0" {
		t.Error("miny =", miny, "want 0")
	}
	if maxx := params.Get("maxx"); maxx != "0" {
		t.Error("maxx =", maxx, "want maxx")
	}
	if maxy := params.Get("maxy"); maxy != "0" {
		t.Error("maxy =", maxy, "want maxy")
	}
	if mapFilter := params.Get("mapfilter"); mapFilter != "false" {
		t.Error("mapfilter =", mapFilter, "want false")
	}
}
