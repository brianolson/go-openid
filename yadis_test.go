// Copyright 2010 Florian Duraffourg. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Getting test data:
// curl -o test_data/py_id.html --dump-header test_data/py_id.http 'http://localhost:8000/id/bob'
// curl -o test_data/google_yadis.html --dump-header test_data/google_yadis.http 'https://www.google.com/accounts/o8/id'
// curl -o test_data/orange_yadis.html --dump-header test_data/orange_yadis.http "http://orange.fr/"
// curl -o test_data/yahoo_yadis.html --dump-header test_data/yahoo_yadis.http "http://yahoo.com/"
// TODO: facebook? livejournal?


package openid

import (
	"testing"
	"bytes"
)

// searchHTMLMetaXRDS Test

type searchHTMLMetaXRDSTest struct {
	in  []byte
	out string
}

var searchHTMLMetaXRDSTests = []searchHTMLMetaXRDSTest{
	searchHTMLMetaXRDSTest{[]byte("<html><head><meta http-equiv='X-XRDS-Location' content='location'></meta></head></html>"), "location"},
	searchHTMLMetaXRDSTest{[]byte("<html><head><meta http-equiv='X-XRDS-Location' content='location'></head></html>"), "location"},
	searchHTMLMetaXRDSTest{[]byte("<html><head><meta http-equiv=\"x-xrds-location\" content=\"location\"></head></html>"), "location"},
	//searchHTMLMetaXRDSTest{[]byte("<html><head><meta>location</meta></head></html>"), "location"},
}

func TestSearchHTMLMetaXRDS(t *testing.T) {
	for _, l := range searchHTMLMetaXRDSTests {
		content, err := searchHTMLMetaXRDS(bytes.NewBuffer(l.in))
		if err != nil {
			t.Errorf("searchHTMLMetaXRDS error: %s", err.String())
		}
		if !bytes.Equal([]byte(content), []byte(l.out)) {
			t.Errorf("searchHTMLMetaXRDS(%s) = %s want %s.", l.in, content, l.out)
		}
	}
}

// Yadis Test

type YadisTest struct {
	url string
}

var YadisTests = []YadisTest{
	YadisTest{"https://www.google.com/accounts/o8/id"},
	YadisTest{"http://orange.fr/"},
	YadisTest{"http://yahoo.com/"},
}

// Test whether the Yadis function returns no errors and a non nil reader
// Doesn't test the content received
func TestYadis(t *testing.T) {
	for _, yt := range YadisTests {
		var reader, err = Yadis(yt.url)
		if err != nil {
			t.Errorf("Yadis(%s) returned a error: %s", yt.url, err.String())
			continue
		}
		if reader == nil {
			t.Errorf("Yadis(%s) returned a nil reader", yt.url)
		}
	}
}

