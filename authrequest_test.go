// Copyright 2010 Florian Duraffourg. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package openid

import (
	"testing"
	"bytes"
)

// NormalizeIdentifier Test

type NormalizeIdentifierTest struct {
	in, out string
	t       int
}

var NormalizeIdentifierTests = []NormalizeIdentifierTest{
	//NormalizeIdentifierTest{"example.com", "http://example.com/", IdentifierURL},
	//NormalizeIdentifierTest{"http://example.com", "http://example.com/", IdentifierURL},
	NormalizeIdentifierTest{"https://example.com/", "https://example.com/", IdentifierURL},
	NormalizeIdentifierTest{"http://example.com/user", "http://example.com/user", IdentifierURL},
	NormalizeIdentifierTest{"http://example.com/user/", "http://example.com/user/", IdentifierURL},
	NormalizeIdentifierTest{"http://example.com/", "http://example.com/", IdentifierURL},
	NormalizeIdentifierTest{"=example", "=example", IdentifierXRI},
	NormalizeIdentifierTest{"xri://=example", "=example", IdentifierXRI},
}

func TestNormalizeIdentifier(testing *testing.T) {
	for _, nit := range NormalizeIdentifierTests {
		v, t := NormalizeIdentifier(nit.in)
		if !bytes.Equal([]byte(v), []byte(nit.out)) || t != nit.t {
			testing.Errorf("NormalizeIdentifier(%s) = (%s, %d) want (%s, %d).", nit.in, v, t, nit.out, nit.t)
		}
	}
}

// GetRedirectURL Test

var Identifiers = []string{
	"https://www.google.com/accounts/o8/id",
	"orange.fr",
	"yahoo.com",
}

// Just check that there is no errors returned by GetRedirectURL
func TestGetRedirectURL(t *testing.T) {
	for _, url := range Identifiers {
		_, err := GetRedirectURL(url, "http://example.com", "/loginCheck")
		if err != nil {
			t.Errorf("GetRedirectURL() returned the error: %s", err.String())
		}
	}
}
