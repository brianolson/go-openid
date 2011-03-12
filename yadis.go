// Copyright 2010 Florian Duraffourg. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package openid

import (
	"os"
	"http"
	"xml"
	"fmt"
	"io"
	"bytes"
	"strings"
)


func Yadis(ID string) (io.Reader, os.Error) {
	r, err := YadisRequest(ID, "GET")
	if (err != nil || r == nil) {
		return nil, err
	}

	var contentType = r.Header.Get("Content-Type")

	// If it is an XRDS document, return the Reader
	if strings.HasPrefix(contentType, "application/xrds+xml") {
		return r.Body, nil
	}
	
	// If it is an HTML doc search for meta tags
	if bytes.Equal([]byte(contentType), []byte("text/html")) {
		url, err := searchHTMLMetaXRDS(r.Body)
		if err != nil {
			return nil, err
		}
		return Yadis(url)
	}
	
	// If the response contain an X-XRDS-Location header
	var xrds_location = r.Header.Get("X-Xrds-Location")
	if len(xrds_location) > 0 {
		return Yadis(xrds_location)
	}

	// If nothing is found try to parse it as a XRDS doc
	return nil, nil
}

func YadisRequest (url string, method string) (resp *http.Response, err os.Error) {
	resp = nil

	var request = new(http.Request)
	var client = new(http.Client)
	var Header = make(http.Header)

	request.Method = method
	request.RawURL = url
	
	request.URL , err = http.ParseURL(url)
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

		if response.StatusCode == 301 || response.StatusCode == 302 || response.StatusCode == 303 || response.StatusCode == 307 {
			location := response.Header.Get("Location")
			request.RawURL = location
			request.URL , err = http.ParseURL(location)
			if err != nil {
				return
			}
		} else {
			return response, nil
		}
	}
	return nil, os.ErrorString("Too many redirections")
}

func searchHTMLMetaXRDS(r io.Reader) (string, os.Error) {
	parser := xml.NewParser(r)
	var token xml.Token
	var err os.Error
	for {
		token, err = parser.Token();
		if (token == nil || err != nil) {
			if err == os.EOF {
				break;
			}
			return "", err
		}

		switch token.(type) {
		case xml.StartElement:
			if token.(xml.StartElement).Name.Local == "meta" {
				// Found a meta token. Verify that it is a X-XRDS-Location and return the content
				var content string
				var contentE bool
				var httpEquivOK bool
				contentE = false
				httpEquivOK = false
				for _,v := range token.(xml.StartElement).Attr {
					if v.Name.Local == "http-equiv" && v.Value == "X-XRDS-Location" {
						httpEquivOK = true
					}
					if v.Name.Local == "content" {
						content = v.Value
						contentE = true
					}
				}
				if contentE && httpEquivOK {
					return fmt.Sprint(content), nil
				}
			}
		}
	}
	return "",os.ErrorString("Value not found")
}
