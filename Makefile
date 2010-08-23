include $(GOROOT)/src/Make.$(GOARCH)

TARG=openid
GOFILES=\
	openid.go\
	yadis.go \
	http.go

include $(GOROOT)/src/Make.pkg

