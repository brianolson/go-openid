include $(GOROOT)/src/Make.$(GOARCH)

TARG=openid
GOFILES=\
	openid.go\
	yadis.go \
	xrds.go \
	http.go

include $(GOROOT)/src/Make.pkg

