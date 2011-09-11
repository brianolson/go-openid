// Copyright 2010 Florian Duraffourg. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package openid

import (
	"os"
	"http"
	"url"
	"io"
	"io/ioutil"
	"bytes"
	"log"
	"regexp"
	"strings"
)

func Yadis(ID string) (io.Reader, os.Error) {
	return YadisVerbose(ID, nil)
}

func YadisVerbose(ID string, verbose *log.Logger) (io.Reader, os.Error) {
	r, err := YadisRequest(ID, "GET")
	if err != nil || r == nil {
		return nil, err
	}

	var contentType = r.Header.Get("Content-Type")

	// If it is an XRDS document, return the Reader
	if strings.HasPrefix(contentType, "application/xrds+xml") {
		if verbose != nil {
			verbose.Printf("got xrds from \"%s\"", ID)
		}
		return r.Body, nil
	}

	// If it is an HTML doc search for meta tags
	if bytes.Equal([]byte(contentType), []byte("text/html")) {
		url_, err := searchHTMLMetaXRDS(r.Body)
		if err != nil {
			return nil, err
		}
		if verbose != nil {
			verbose.Printf("fetching xrds found in html \"%s\"", url_)
		}
		return Yadis(url_)
	}

	// If the response contain an X-XRDS-Location header
	var xrds_location = r.Header.Get("X-Xrds-Location")
	if len(xrds_location) > 0 {
		if verbose != nil {
			verbose.Printf("fetching xrds found in http header \"%s\"", xrds_location)
		}
		return Yadis(xrds_location)
	}

	if verbose != nil {
		verbose.Printf("Yadis fails out, nothing found. status=%#v", r.StatusCode)
	}
	// If nothing is found try to parse it as a XRDS doc
	return nil, nil
}

func YadisRequest(url_ string, method string) (resp *http.Response, err os.Error) {
	resp = nil

	var request = new(http.Request)
	var client = new(http.Client)
	var Header = make(http.Header)

	request.Method = method
	request.RawURL = url_

	request.URL, err = url.Parse(url_)
	if err != nil {
		return
	}

	// Common parameters
	request.Proto = "HTTP/1.0"
	request.ProtoMajor = 1
	request.ProtoMinor = 0
	request.ContentLength = 0
	request.Close = true

	Header.Add("Accept", "application/xrds+xml")
	request.Header = Header

	// Follow a maximum of 5 redirections
	for i := 0; i < 5; i++ {
		response, err := client.Do(request)

		if err != nil {
			return nil, err
		}
		if response.StatusCode == 301 || response.StatusCode == 302 || response.StatusCode == 303 || response.StatusCode == 307 {
			location := response.Header.Get("Location")
			request.RawURL = location
			request.URL, err = url.Parse(location)
			if err != nil {
				return
			}
		} else {
			return response, nil
		}
	}
	return nil, os.NewError("Too many redirections")
}

var metaRE *regexp.Regexp
var xrdsRE *regexp.Regexp

func init() {
	// These are ridiculous case insensitive pattern constructions.

	// <[ \t]*meta[^>]*http-equiv=["']x-xrds-location["'][^>]*>
	metaRE = regexp.MustCompile("<[ \t]*[mM][eE][tT][aA][^>]*[hH][tT][tT][pP]-[eE][qQ][uU][iI][vV]=[\"'][xX]-[xX][rR][dD][sS]-[lL][oO][cC][aA][tT][iI][oO][nN][\"'][^>]*>")

	// content=["']([^"']+)["']
	xrdsRE = regexp.MustCompile("[cC][oO][nN][tT][eE][nN][tT]=[\"']([^\"]+)[\"']")
}

func searchHTMLMetaXRDS(r io.Reader) (string, os.Error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	part := metaRE.Find(data)
	if part == nil {
		return "", os.NewError("No -meta- match")
	}
	content := xrdsRE.FindSubmatch(part)
	if content == nil {
		return "", os.NewError("No content in meta tag: " + string(part))
	}
	return string(content[1]), nil
}
