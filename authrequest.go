// Copyright 2010 Florian Duraffourg. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package openid

import (
	"strings"
	"io"
	"os"
	"url"
)

const (
	_ = iota
	IdentifierXRI
	IdentifierURL
)

func GetRedirectURL(Identifier string, realm string, returnto string) (string, os.Error) {
	var err os.Error
	var Id, IdType = NormalizeIdentifier(Identifier)

	// If the identifier is an XRI, [XRI_Resolution_2.0] will yield an XRDS document that contains the necessary information. It should also be noted that Relying Parties can take advantage of XRI Proxy Resolvers, such as the one provided by XDI.org at http://www.xri.net. This will remove the need for the RPs to perform XRI Resolution locally.
	if IdType == IdentifierXRI {
		// Not implemented yet
		return "", os.NewError("XRI identifier not implemented yed")
	}

	// If it is a URL, the Yadis protocol [Yadis] SHALL be first attempted. If it succeeds, the result is again an XRDS document.
	if IdType == IdentifierURL {
		var reader io.Reader
		reader, err = Yadis(Id)
		if err != nil {
			return "", err
		}
		if reader == nil {
			return "", os.NewError("Yadis returned an empty Reader for the ID: " + Id)
		}

		var endpoint, claimedid = ParseXRDS(reader)
		if len(endpoint) == 0 {
			return "", os.NewError("Unable to parse the XRDS document")
		}

		// At this point we have the endpoint and eventually a claimed id		 
		// Create the authentication request
		return CreateAuthenticationRequest(endpoint, claimedid, realm, returnto), nil
	}

	// If the Yadis protocol fails and no valid XRDS document is retrieved, or no Service Elements are found in the XRDS document, the URL is retrieved and HTML-Based discovery SHALL be attempted.


	return "Not implemented", nil
}

func NormalizeIdentifier(Id string) (Identifier string, IdentifierType int) {
	Identifier = Id
	//1.  If the user's input starts with the "xri://" prefix, it MUST be stripped off, so that XRIs are used in the canonical form.
	if strings.HasPrefix(Identifier, "xri://") {
		Identifier = Identifier[6:]
	}

	// 2. If the first character of the resulting string is an XRI Global Context Symbol ("=", "@", "+", "$", "!") or "(", as defined in Section 2.2.1 of [XRI_Syntax_2.0] (Reed, D. and D. McAlpin, “Extensible Resource Identifier (XRI) Syntax V2.0,” .), then the input SHOULD be treated as an XRI.
	var firstChar = Identifier[0]
	if firstChar == '=' || firstChar == '@' || firstChar == '+' || firstChar == '$' || firstChar == '!' {
		IdentifierType = IdentifierXRI
		return
	}

	// 3. Otherwise, the input SHOULD be treated as an http URL; if it does not include a "http" or "https" scheme, the Identifier MUST be prefixed with the string "http://". If the URL contains a fragment part, it MUST be stripped off together with the fragment delimiter character "#". See Section 11.5.2 (HTTP and HTTPS URL Identifiers) for more information.
	IdentifierType = IdentifierURL
	if !strings.HasPrefix(Identifier, "http://") && !strings.HasPrefix(Identifier, "https://") {
		Identifier = "http://" + Identifier
	}

	// 4. URL Identifiers MUST then be further normalized by both following redirects when retrieving their content and finally applying the rules in Section 6 of [RFC3986] (Berners-Lee, T., “Uniform Resource Identifiers (URI): Generic Syntax,” .) to the final destination URL. This final URL MUST be noted by the Relying Party as the Claimed Identifier and be used when requesting authentication (Requesting Authentication).

	return
}

func CreateAuthenticationRequest(OPEndPoint, ClaimedID, Realm, ReturnTo string) string {
	var p = make(map[string]string)

	p["openid.ns"] = "http://specs.openid.net/auth/2.0"
	p["openid.mode"] = "checkid_setup"

	if len(ClaimedID) == 0 {
		p["openid.claimed_id"] = "http://specs.openid.net/auth/2.0/identifier_select"
		p["openid.identity"] = "http://specs.openid.net/auth/2.0/identifier_select"
	} else {
		p["openid.claimed_id"] = ClaimedID
		p["openid.identity"] = ClaimedID
	}

	p["openid.return_to"] = Realm + ReturnTo
	p["openid.realm"] = Realm

	var url_ string
	url_ = OPEndPoint + "?"

	for k, v := range p {
		url_ += url.QueryEscape(k) + "=" + url.QueryEscape(v) + "&"
	}
	return url_
}
