// Copyright 2010 Florian Duraffourg. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package openid

import (
	"xml"
	"io"
	"strings"
)

type XRDSIdentifier struct {
	XMLName xml.Name "Service"
	Type    []string
	URI     string
	LocalID string
}
type XRD struct {
	XMLName xml.Name "XRD"
	Service XRDSIdentifier
}
type XRDS struct {
	XMLName xml.Name "XRDS"
	XRD     XRD
}

// Parse a XRDS document provided through a io.Reader
// Return the OP EndPoint and, if found, the Claimed Identifier
func ParseXRDS(r io.Reader) (string, string) {
	XRDS := new(XRDS)
	err := xml.Unmarshal(r, XRDS)
	if err != nil {
		//fmt.Printf(err.String())
		return "", ""
	}
	XRDSI := XRDS.XRD.Service

	XRDSI.URI = strings.TrimSpace(XRDSI.URI)
	XRDSI.LocalID = strings.TrimSpace(XRDSI.LocalID)

	//fmt.Printf("%v\n", XRDSI)

	if StringTableContains(XRDSI.Type, "http://specs.openid.net/auth/2.0/server") {
		//fmt.Printf("OP Identifier Element found\n")
		return XRDSI.URI, ""
	} else if StringTableContains(XRDSI.Type, "http://specs.openid.net/auth/2.0/signon") {
		//fmt.Printf("Claimed Identifier Element found\n")
		return XRDSI.URI, XRDSI.LocalID
	}
	return "", ""
}


func StringTableContains(t []string, s string) bool {
	for _, v := range t {
		if v == s {
			return true
		}
	}
	return false
}
