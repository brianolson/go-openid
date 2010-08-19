include $(GOROOT)/src/Make.$(GOARCH)

TARG=openid
GOFILES=\
	openid.go\
	http.go

include $(GOROOT)/src/Make.pkg

