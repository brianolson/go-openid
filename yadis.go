package openid

import (
	"io"
	"regexp"
	"fmt"
	"strings"
)

func ParseXRDS(r io.Reader) string {
	var buffer = make([]byte,1024)
	io.ReadFull(r, buffer)
	URIRegex := regexp.MustCompile("<URI[^>]*>.*</URI>")
	uris := URIRegex.MatchStrings(string(buffer))
	if len(uris) < 1 {
		return ""
	}
	uri := uris[0]
	start := strings.Index(uri,">") + 1
	end := strings.Index(uri,"</")
	uri = uri[start:end]
	fmt.Printf("%s\n",uri)
	return uri
}

func Yadis(url string) string{
	fmt.Printf("Search: %s\n",url)
	headers := map[string] string {
		"Accept": "application/xrds+xml",
	}
	r, err := get (url, headers)
	if (err != nil || r == nil) { 
		fmt.Printf("Error in GET\n")
		return "" }

	// If it is an XRDS document, parse it and return URI
	content, ok := r.Header["Content-Type"]
	if ok && content == "application/xrds+xml" {
		fmt.Printf("Document XRDS found\n")
		return ParseXRDS(r.Body)
	}
	

	// If it is an HTML doc search for meta tags

	// If the response contain an X-XRDS-Location header
	xrds, ok := r.Header["X-Xrds-Location"]
	if ok {
		return Yadis(xrds)
	}

	// If nothing is found try to parse it as a XRDS doc
	return ParseXRDS(r.Body)
	
}
