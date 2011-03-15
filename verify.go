// Copyright 2010 Florian Duraffourg. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package openid

import (
	"os"
	"http"
	"regexp"
	"bytes"
)

// Verify that the url given match a successfull authentication
// Return:
// * true if authenticated, false otherwise
// * The Claimed identifier if authenticated
// * Eventually an error
func Verify(url string) (grant bool, identifier string, err os.Error) {
	grant = false
	identifier = ""
	err = nil

	var urlm map[string]string
	urlm, err = url2map(url)
	if err != nil {
		return false, "", err
	}

	// The value of "openid.return_to" matches the URL of the current request (Section 11.1)
	// To be implemented in a global way

	// Discovered information matches the information in the assertion (Section 11.2)

	// An assertion has not yet been accepted from this OP with the same value for "openid.response_nonce" (Section 11.3)

	// The signature on the assertion is valid and all fields that are required to be signed are signed (Section 11.4)

	grant, err = verifyDirect(urlm)
	if err != nil {
		return
	}

	return
}

var REVerifyDirectIsValid = "is_valid=true"
var REVerifyDirectNs = regexp.MustCompile("ns=([^&]*)")

func verifyDirect(urlm map[string]string) (grant bool, err os.Error) {
	grant = false
	err = nil

	urlm["openid.mode"] = "check_authentication"

	// Create the url
	URLEndPoint := urlm["openid.op_endpoint"]
	var postContent string
	for k, v := range urlm {
		postContent += http.URLEscape(k) + "=" + http.URLEscape(v) + "&"
	}

	// Post the request
	var client = new(http.Client)
	postReader := bytes.NewBuffer([]byte(postContent))
	response, err := client.Post(URLEndPoint, "application/x-www-form-urlencoded", postReader)
	if err != nil {
		return false, err
	}

	// Parse the response
	// Convert the reader -- Warning, response.ContentLength might be -1!!
	buffer := make([]byte, response.ContentLength)
	_, err = response.Body.Read(buffer)
	if err != nil {
		return false, err
	}

	// Check for ns
	rematchs := REVerifyDirectNs.FindSubmatch(buffer)
	if len(rematchs) < 1 {
		return false, os.ErrorString("verifyDirect: ns value not found on the response of the OP")
	}
	nsValue, err := http.URLUnescape(string(rematchs[1]))
	if err != nil {
		return false, err
	}
	if !bytes.Equal([]byte(nsValue), []byte("http://specs.openid.net/auth/2.0")) {
		return false, os.ErrorString("verifyDirect: ns value not correct: " + nsValue)
	}

	// Check for is_valid
	match, err := regexp.Match(REVerifyDirectIsValid, buffer)
	if err != nil {
		return false, err
	}

	return match, nil
}

// Transform an url string into a map of parameters/value
func url2map(url string) (map[string]string, os.Error) {
	pmap := make(map[string]string)
	var start, end, eq, length int
	var param, value string
	var err os.Error

	length = len(url)
	start = 0
	for start < length && url[start] != '?' {
		start++
	}
	if start >= length {
		start = -1
	}
	end = start
	for end < length {
		start = end + 1
		eq = start
		for eq < length && url[eq] != '=' {
			eq++
		}
		end = eq + 1
		for end < length && url[end] != '&' {
			end++
		}

		param, err = http.URLUnescape(url[start:eq])
		if err != nil {
			return nil, err
		}
		value, err = http.URLUnescape(url[eq+1 : end])
		if err != nil {
			return nil, err
		}

		pmap[param] = value
	}
	return pmap, nil
}
