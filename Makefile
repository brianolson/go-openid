include $(GOROOT)/src/Make.inc

TARG=openid
GOFILES=\
	authrequest.go\
	xrds.go\
	yadis.go

include $(GOROOT)/src/Make.pkg

