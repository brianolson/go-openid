// Copyright 2010 Florian Duraffourg. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package openid

import (
	"log"
	"os"
	"http"
	"regexp"
	"bytes"
	"url"
)

// Verify that the url given match a successfull authentication
// Return:
// * true if authenticated, false otherwise
// * The Claimed identifier if authenticated
// * Eventually an error
func Verify(url_ string) (grant bool, identifier string, err os.Error) {
	grant = false
	identifier = ""
	err = nil

	//var urlm map[string]string
	//urlm, err = url2map(url)
	var values url.Values
	values, err = url.ParseQuery(url_)
	if err != nil {
		return false, "", err
	}

	// The value of "openid.return_to" matches the URL of the current request (Section 11.1)
	// To be implemented in a global way

	// Discovered information matches the information in the assertion (Section 11.2)

	// An assertion has not yet been accepted from this OP with the same value for "openid.response_nonce" (Section 11.3)

	// The signature on the assertion is valid and all fields that are required to be signed are signed (Section 11.4)

	return VerifyValues(values)
	//if err != nil {
	//	return grant, identifier, err
	//}

	//identifier = urlm["openid.claimed_id"]

	//return grant, identifier, err
}

var REVerifyDirectIsValid = "is_valid:true"
var REVerifyDirectNs = regexp.MustCompile("ns:([a-zA-Z0-9:/.]*)")

// Like Verify on a parsed URL
func VerifyValues(values url.Values) (grant bool, identifier string, err os.Error) {
	err = nil

	var postArgs url.Values
	postArgs = url.Values(map[string][]string{})
	//postArgs = new(http.Values)
	postArgs.Set("openid.mode", "check_authentication")

	// Create the url
	URLEndPoint := values.Get("openid.op_endpoint")
	if URLEndPoint == "" {
		log.Printf("no openid.op_endpoint")
		return false, "", os.NewError("no openid.op_endpoint")
	}
	for k, v := range values {
		if k == "openid.op_endpoint" {
			continue // skip it
		}
		postArgs[k] = v
	}
	postContent := postArgs.Encode()

	// Post the request
	var client = new(http.Client)
	postReader := bytes.NewBuffer([]byte(postContent))
	response, err := client.Post(URLEndPoint, "application/x-www-form-urlencoded", postReader)
	if err != nil {
		log.Printf("VerifyValues failed at post")
		return false, "", err
	}

	// Parse the response
	// Convert the reader
	// We limit the size of the response to 1024 bytes but it should be large enough for most cases
	buffer := make([]byte, 1024)
	_, err = response.Body.Read(buffer)
	if err != nil {
		log.Printf("VerifyValues failed reading response")
		return false, "", err
	}

	// Check for ns
	rematch := REVerifyDirectNs.FindSubmatch(buffer)
	if rematch == nil {
		return false, "", os.NewError("VerifyValues: ns value not found on the response of the OP")
	}
	nsValue := string(rematch[1])
	if !bytes.Equal([]byte(nsValue), []byte("http://specs.openid.net/auth/2.0")) {
		return false, "", os.NewError("VerifyValues: ns value not correct: " + nsValue)
	}

	// Check for is_valid
	match, err := regexp.Match(REVerifyDirectIsValid, buffer)
	if err != nil {
		return false, "", err
	}

	identifier = values.Get("openid.claimed_id")

	return match, identifier, nil
}

// Transform an url string into a map of parameters/value
func url2map(url_ string) (map[string]string, os.Error) {
	pmap := make(map[string]string)
	var start, end, eq, length int
	var param, value string
	var err os.Error

	length = len(url_)
	start = 0
	for start < length && url_[start] != '?' {
		start++
	}
	if start >= length {
		start = -1
	}
	end = start
	for end < length {
		start = end + 1
		eq = start
		for eq < length && url_[eq] != '=' {
			eq++
		}
		end = eq + 1
		for end < length && url_[end] != '&' {
			end++
		}

		param, err = url.QueryUnescape(url_[start:eq])
		if err != nil {
			return nil, err
		}
		value, err = url.QueryUnescape(url_[eq+1 : end])
		if err != nil {
			return nil, err
		}

		pmap[param] = value
	}
	return pmap, nil
}
