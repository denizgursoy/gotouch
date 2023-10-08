package auth

import (
	"net/http"

	"github.com/jdx/go-netrc"
)

var _ http.RoundTripper = (*NetrcTransport)(nil)

func NewAuthenticatedHTTPClient() *http.Client {
	return &http.Client{
		Transport: NewNetrcTransport(&http.Transport{}),
	}
}

func NewNetrcTransport(underlying http.RoundTripper) http.RoundTripper {
	if underlying == nil {
		underlying = http.DefaultTransport
	}

	creds, err := ParseNetrc()
	if err != nil {
		return underlying
	}

	return &NetrcTransport{
		credentials: creds,
		Underlying:  underlying,
	}
}

type NetrcTransport struct {
	credentials *netrc.Netrc
	Underlying  http.RoundTripper
}

func (n NetrcTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if mc := n.credentials.Machine(request.Host); mc != nil {
		request.SetBasicAuth(mc.Get(netrcLoginKey), mc.Get(netrcPasswordKey))
	}

	if n.Underlying != nil {
		return n.Underlying.RoundTrip(request)
	}

	return http.DefaultTransport.RoundTrip(request)
}
