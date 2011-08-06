Go-OpenID
=========

About
-----

Go-OpenID is an implementation of OpenID authentication protocol in Golang.

For now, the implementation does not manage XRI identifier, and can only check authentication with a direct request.

Here are the specifications used:

* http://openid.net/specs/openid-authentication-2_0.html
* http://yadis.org/wiki/Yadis_1.0_%28HTML%29

Install
-------

	git clone http://github.com/fduraffourg/go-openid.git && cd go-openid && make && make install

or
	goinstall github.com/fduraffourg/go-openid


Usage
-----

        url := openid.GetRedirectURL("Identifier", "http://www.realm.com", "/loginCheck")

Now you have to redirect the user to the url returned. The OP will then forward the user back to you, after authenticating him.

To check the identity, do that:

   	 grant, id, err := openid.Verify(URL)

URL is the url the user was redirected to.  grant will be true if the
user was correctly authenticated, false otherwise.  If the user was
authenticated, id contains its identifier.
