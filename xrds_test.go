// Copyright 2010 Florian Duraffourg. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package openid

import (
	"testing"
	"bytes"
)

// ParseXRDS Test

type ParseXRDSTest struct {
	in         []byte
	OPEndPoint string
	ClaimedId  string
}

var ParseXRDSTests = []ParseXRDSTest{
	ParseXRDSTest{[]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?><xrds:XRDS xmlns:xrds=\"xri://$xrds\" xmlns=\"xri://$xrd*($v*2.0)\"><XRD><Service xmlns=\"xri://$xrd*($v*2.0)\">\n<Type>http://specs.openid.net/auth/2.0/signon</Type>\n  <URI>https://www.exampleprovider.com/endpoint/</URI>\n  <LocalID>https://exampleuser.exampleprovider.com/</LocalID>\n		</Service></XRD></xrds:XRDS>"), "https://www.exampleprovider.com/endpoint/", "https://exampleuser.exampleprovider.com/"},
	ParseXRDSTest{[]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<xrds:XRDS xmlns:xrds=\"xri://$xrds\" xmlns=\"xri://$xrd*($v*2.0)\">\n<XRD>\n    <Service>\n        <Type>http://specs.openid.net/auth/2.0/server</Type>\n        <Type>http://openid.net/srv/ax/1.0</Type>\n        <Type>http://openid.net/sreg/1.0</Type>\n        <Type>http://openid.net/extensions/sreg/1.1</Type>\n        <URI priority=\"20\">http://openid.orange.fr/server/</URI>\n    </Service>\n</XRD>\n</xrds:XRDS>"), "http://openid.orange.fr/server/", ""},
}

func TestParseXRDS(t *testing.T) {
	for _, xrds := range ParseXRDSTests {
		var opep, ci = ParseXRDS(bytes.NewBuffer(xrds.in))
		if !bytes.Equal([]byte(opep), []byte(xrds.OPEndPoint)) || !bytes.Equal([]byte(ci), []byte(xrds.ClaimedId)) {
			t.Errorf("ParseXRDS(%s) = (%s, %s) want (%s, %s).", xrds.in, opep, ci, xrds.OPEndPoint, xrds.ClaimedId)
		}
	}
}
