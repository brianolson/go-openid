package openid

import (
	"fmt"
	"io"
	"regexp"
	"os"
	"bytes"

	)


type OpenID struct {
	Identifier string
	//Params map[string] string
	RPUrl string
	Hostname string
	Request string
}

func (o *OpenID) normalizeIdentifier() {
	return

}

func Yadis(url string) string{

	r, err := get (url, nil)
	if (err != nil) { return "" }

	var buffer = make([]byte,1024)
	io.ReadFull(r.Body, buffer)
	
	URIRegex := regexp.MustCompile("<URI>.*</URI>")
	uris := URIRegex.MatchStrings(string(buffer))
	if len(uris) < 1 {
		return ""
	}
	
	uri := uris[0][5:len(uris[0])-6]
	fmt.Printf("%s\n",uri)
	return uri
}

func (o *OpenID) GetUrl() string {
	o.normalizeIdentifier()

	URI := Yadis(o.Identifier)
	params := fmt.Sprintf("?openid.ns=http://specs.openid.net/auth/2.0&openid.claimed_id=http://specs.openid.net/auth/2.0/identifier_select&openid.identity=http://specs.openid.net/auth/2.0/identifier_select&openid.return_to=http://157.159.46.13:8083/go/loginCheck&openid.realm=http://157.159.46.13:8083/&openid.mode=checkid_setup")
	return fmt.Sprintf("%s%s",URI,params)
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


func (o *OpenID) VerifyDirect() (grant bool, err os.Error) {
	grant = false
	err = nil

	//o.Params["openid.mode"] = "check_authentication"

	// Create the new post
	// Copy everything exept for openid.mode
	params := "openid.mode=check_authentication"
	url := ""

	var start, end, eq, length int
	length = len(o.RPUrl)
	start = 0
	for start < length && o.RPUrl[start] != '?' { start ++ }
	end = start
	for end < length {
		start = end + 1
		eq = start
		for eq < length && o.RPUrl[eq] != '=' { eq++ }
		end = eq + 1
		for end < length && o.RPUrl[end] != '&' { end++ }
	
		fmt.Printf("Trouve: %s : %s\n", o.RPUrl[start:eq], o.RPUrl[eq+1:end])
		if o.RPUrl[start:eq] != "openid.mode" {
			params = fmt.Sprintf("%s&%s=%s", params, o.RPUrl[start:eq], o.RPUrl[eq+1:end])
		}
		if o.RPUrl[start:eq] == "openid.op_endpoint" {
			url = o.RPUrl[eq+1:end]
		}
	}

	
	fmt.Printf("url: %s\nparams: %s\n", url, params)
	fmt.Printf("\n")
	// We can now send the direct request for verification and check the result
	headers := map[string] string {
		"Content-Type" : "application/x-www-form-urlencoded",
	}
	r,error := post(url, headers, bytes.NewBuffer([]byte(params)))
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
