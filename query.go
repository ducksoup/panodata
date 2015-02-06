package panodata

import (
	"math"
	"net/url"
	"strconv"
)

// Set is the target photo set for a query.
type Set byte

const (
	SetFull   Set = iota // Public photos
	SetPublic            // All photos
)

// String converts a Set to its string representation in the API.
func (set Set) String() string {
	switch set {
	case SetFull:
		return "full"
	case SetPublic:
		return "public"
	}
	return ""
}

// Size is the photo size Panoramio will return links to.
type Size byte

const (
	SizeOriginal Size = iota
	SizeMedium
	SizeSmall
	SizeThumbnail
	SizeSquare
	SizeMiniSquare
)

// String convert a Size to its string representation in the API.
func (size Size) String() string {
	switch size {
	case SizeOriginal:
		return "original"
	case SizeMedium:
		return "medium"
	case SizeSmall:
		return "small"
	case SizeThumbnail:
		return "thumbnail"
	case SizeSquare:
		return "square"
	case SizeMiniSquare:
		return "mini_square"
	}
	return ""
}

// A Query into Panoramio's database.
type Query struct {
	MinLat, MaxLat float64
	MinLon, MaxLon float64
	Size           Size
	Set            Set
	MapFilter      bool
}

const rEarth = 6371

// QueryRadius returns a Query that looks for images less than radius Km from
// lat, lon. This function uses an approximation that is only valid for small
// values of radius.
func QueryRadius(lat, lon float64, radius float64) Query {
	x := radius / (rEarth * math.Sqrt(2))

	dLat, dLon := math.Abs(x), math.Abs(x/math.Cos(lat))

	return Query{
		MinLat: lat - dLat,
		MaxLat: lat + dLat,
		MinLon: lon - dLon,
		MaxLon: lon + dLon,
	}
}

// URL to the Panoramio API page that performs the query.
func (query Query) URL() string {
	// Create query
	params := url.Values{}
	params.Add("set", query.Set.String())
	params.Add("size", query.Size.String())
	params.Add("minx", strconv.FormatFloat(query.MinLon, 'f', -1, 64))
	params.Add("miny", strconv.FormatFloat(query.MinLat, 'f', -1, 64))
	params.Add("maxx", strconv.FormatFloat(query.MaxLon, 'f', -1, 64))
	params.Add("maxy", strconv.FormatFloat(query.MaxLat, 'f', -1, 64))
	params.Add("mapfilter", strconv.FormatBool(query.MapFilter))

	url := url.URL{
		Scheme:   "http",
		Host:     "www.panoramio.com",
		Path:     "map/get_panoramas.php",
		RawQuery: params.Encode(),
	}

	return url.String()
}
