package panodata

import "testing"

func TestAddRange(t *testing.T) {
	url := "http://www.test.com/"
	if added := addRange(url, 0, 100); added != "http://www.test.com/&from=0&to=100" {
		t.Error("added = ", added, " want http://www.test.com/&from=0&to=100")
	}
}
