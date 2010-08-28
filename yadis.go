package openid

import (
	"io"
	"fmt"
	"os"
	"xml"
	"strings"
)

func searchHTMLMeta(r io.Reader) (string, os.Error) {
	parser := xml.NewParser(r)
	var token xml.Token
	var err os.Error
	for {
		token, err = parser.Token();
		if (token == nil || err != nil) {
			if err == os.EOF {
				break;
			}
			return "", err
		}

		switch token.(type) {
		case xml.StartElement:
			if token.(xml.StartElement).Name.Local == "meta" {
				// Found a meta token. Verify that it is a X-XRDS-Location and return the content
				var content string
				var contentE bool
				var httpEquivOK bool
				contentE = false
				httpEquivOK = false
				for _,v := range token.(xml.StartElement).Attr {
					if v.Name.Local == "http-equiv" && v.Value == "X-XRDS-Location" {
						httpEquivOK = true
					}
					if v.Name.Local == "content" {
						content = v.Value
						contentE = true
					}
				}
				if contentE && httpEquivOK {
					return fmt.Sprint(content), nil
				}
			}
		}
	}
	return "",os.ErrorString("Value not found")
}

func Yadis(url string) (io.Reader, os.Error) {
	fmt.Printf("Search: %s\n",url)
	headers := map[string] string {
		"Accept": "application/xrds+xml",
	}
	r, err := get (url, headers)
	if (err != nil || r == nil) {
		fmt.Printf("Yadis: Error in GET\n")
		return nil, err
	}

	// If it is an XRDS document, parse it and return URI
	content, ok := r.Header["Content-Type"]
	if ok && strings.HasPrefix(content, "application/xrds+xml") {
		fmt.Printf("Document XRDS found\n")
		return r.Body, nil
	}
	
	// If it is an HTML doc search for meta tags
	content, ok = r.Header["Content-Type"]
	if ok && content == "text/html" {
		fmt.Printf("Document HTML found\n")
		url, err := searchHTMLMeta(r.Body)
		if err != nil {
			return nil, err
		}
		return Yadis(url)
	}
	

	// If the response contain an X-XRDS-Location header
	xrds, ok := r.Header["X-Xrds-Location"]
	if ok {
		return Yadis(xrds)
	}

	// If nothing is found try to parse it as a XRDS doc
	return nil, nil
}
