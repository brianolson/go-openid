package openid

import (
	"fmt"
	"io"
	"regexp"
	"os"
	"bytes"
	"http"
	)


type OpenID struct {
	Identifier string
	Params map[string] string
	RPUrl string
	Hostname string
	Request string
	Realm string
	ReturnTo string
}

func (o *OpenID) normalizeIdentifier() {
	match,_ := regexp.MatchString("https?://", o.Identifier)
	if ! match {
		o.Identifier = fmt.Sprintf("http://%s",o.Identifier)
	}

	return

}


func mapToUrlEnc (params map[string] string) string {
	url := ""
	for k,v := range (params) {
		url = fmt.Sprintf("%s&%s=%s",url,k,v)
	}
	//return http.URLEscape(url[1:])Aa
	return url[1:]
}

func urlEncToMap (url string) map[string] string {
	// We don't know how elements are in the URL so we create a list first and push elements on it
	pmap := make(map[string] string)
	url,_ = http.URLUnescape(url)
	var start, end, eq, length int
	length = len(url)
	start = 0
	for start < length && url[start] != '?' { start ++ }
	end = start
	for end < length {
		start = end + 1
		eq = start
		for eq < length && url[eq] != '=' { eq++ }
		end = eq + 1
		for end < length && url[end] != '&' { end++ }
	
		fmt.Printf("Trouve: %s : %s\n", url[start:eq], url[eq+1:end])
		pmap[url[start:eq]] = url[eq+1:end]
	}
	return pmap
}

func (o *OpenID) GetUrl() string {
	o.normalizeIdentifier()

	URI := Yadis(o.Identifier)
	if URI == "" {
		return ""
	}
	params := map[string] string {
		"openid.ns": "http://specs.openid.net/auth/2.0",
		"openid.mode" : "checkid_setup",
		"openid.return_to": fmt.Sprintf("%s%s", o.Realm, o.ReturnTo),
		"openid.realm": o.Realm,
		"openid.claimed_id" : "http://specs.openid.net/auth/2.0/identifier_select",
		"openid.identity" : "http://specs.openid.net/auth/2.0/identifier_select",

	}
	return fmt.Sprintf("%s?%s",URI, mapToUrlEnc(params))
}

func (o *OpenID) Verify() (grant bool, err os.Error) {
	grant = false
	err = nil
	

	// The value of "openid.return_to" matches the URL of the current request
	// if ! MExists(o.Params, "openid.return_to") {
	// 	err = os.ErrorString("The value of 'openid.return_to' is not defined")
	// 	return
	// }
	// if (fmt.Sprintf("%s%s", o.Hostname, o.Request) != o.Params["openid.return_to"]) {
	// 	err = os.ErrorString("The value of 'openid.return_to' does not match the URL of the current request")
	// 	return
	// }

	// Discovered information matches the information in the assertion

	// An assertion has not yet been accepted from this OP with the same value for "openid.response_nonce"

	// The signature on the assertion is valid and all fields that are required to be signed are signed
	grant, err = o.VerifyDirect()
	
	return
}

func (o *OpenID) ParseRPUrl(url string) {
	o.Params = urlEncToMap(url)
}

func (o *OpenID) VerifyDirect() (grant bool, err os.Error) {
	grant = false
	err = nil

	o.Params["openid.mode"] = "check_authentication"

	headers := map[string] string {
		"Content-Type" : "application/x-www-form-urlencoded",
	}
	fmt.Printf("Verification: %s\n",o.Params["openid.op_endpoint"])
	r,error := post(o.Params["openid.op_endpoint"],
		headers,
		bytes.NewBuffer([]byte(mapToUrlEnc(o.Params))))
	if error != nil {
		fmt.Printf("erreur: %s\n", error.String())
		err = error
		return
	}
	fmt.Printf("Post done\n")
	if (r != nil) {
		buffer := make([]byte, 1024)
		fmt.Printf("Buffer created\n")
		io.ReadFull(r.Body, buffer)
		fmt.Printf("Body extracted: %s\n", buffer)
		grant, err = regexp.Match("is_valid:true", buffer)
		fmt.Printf("Response: %v\n", grant)
	}else {
		err = os.ErrorString("No response from POST verification")
		return
	}

	return
}


func MExists(datas map[string] string, index string) bool {
	_, present := datas[index]
	return present
}
