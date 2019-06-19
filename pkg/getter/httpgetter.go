/*
Copyright The Helm Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package getter

import (
	"bytes"
	"io"
	"net/http"

	"github.com/pkg/errors"

	"helm.sh/helm/pkg/tlsutil"
	"helm.sh/helm/pkg/urlutil"
)

// HTTPGetter is the efault HTTP(/S) backend handler
type HTTPGetter struct {
	client *http.Client
	opts   options
}

// SetBasicAuth sets the request's Authorization header to use the provided credentials.
func (g *HTTPGetter) SetBasicAuth(username, password string) {
	g.opts.username = username
	g.opts.password = password
}

// SetUserAgent sets the request's User-Agent header to use the provided agent name.
func (g *HTTPGetter) SetUserAgent(userAgent string) {
	g.opts.userAgent = userAgent
}

//Get performs a Get from repo.Getter and returns the body.
func (g *HTTPGetter) Get(href string) (*bytes.Buffer, error) {
	return g.get(href)
}

func (g *HTTPGetter) get(href string) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)

	// Set a helm specific user agent so that a repo server and metrics can
	// separate helm calls from other tools interacting with repos.
	req, err := http.NewRequest("GET", href, nil)
	if err != nil {
		return buf, err
	}
	// req.Header.Set("User-Agent", "Helm/"+strings.TrimPrefix(version.GetVersion(), "v"))
	if g.opts.userAgent != "" {
		req.Header.Set("User-Agent", g.opts.userAgent)
	}

	if g.opts.username != "" && g.opts.password != "" {
		req.SetBasicAuth(g.opts.username, g.opts.password)
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return buf, err
	}
	if resp.StatusCode != 200 {
		return buf, errors.Errorf("failed to fetch %s : %s", href, resp.Status)
	}

	_, err = io.Copy(buf, resp.Body)
	resp.Body.Close()
	return buf, err
}

// newHTTPGetter constructs a valid http/https client as Getter
func newHTTPGetter(options ...Option) (Getter, error) {
	return NewHTTPGetter(options...)
}

// NewHTTPGetter constructs a valid http/https client as HTTPGetter
func NewHTTPGetter(options ...Option) (*HTTPGetter, error) {
	var client HTTPGetter

	for _, opt := range options {
		opt(&client.opts)
	}

	if client.opts.certFile != "" && client.opts.keyFile != "" {
		tlsConf, err := tlsutil.NewClientTLS(client.opts.certFile, client.opts.keyFile, client.opts.caFile)
		if err != nil {
			return &client, errors.Wrap(err, "can't create TLS config for client")
		}
		tlsConf.BuildNameToCertificate()

		sni, err := urlutil.ExtractHostname(client.opts.url)
		if err != nil {
			return &client, err
		}
		tlsConf.ServerName = sni

		client.client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConf,
				Proxy:           http.ProxyFromEnvironment,
			},
		}
	} else {
		client.client = http.DefaultClient
	}

	return &client, nil
}
