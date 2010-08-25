package openid

import (
	"io"
	"fmt"
)

func Yadis(url string) io.Reader {
	fmt.Printf("Search: %s\n",url)
	headers := map[string] string {
		"Accept": "application/xrds+xml",
	}
	r, err := get (url, headers)
	if (err != nil || r == nil) { 
		fmt.Printf("Error in GET\n")
		return nil
	}

	// If it is an XRDS document, parse it and return URI
	content, ok := r.Header["Content-Type"]
	if ok && content == "application/xrds+xml" {
		fmt.Printf("Document XRDS found\n")
		return r.Body
	}
	

	// If it is an HTML doc search for meta tags

	// If the response contain an X-XRDS-Location header
	xrds, ok := r.Header["X-Xrds-Location"]
	if ok {
		return Yadis(xrds)
	}

	// If nothing is found try to parse it as a XRDS doc
	return r.Body
	
}
