package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"httpfromtcp/rootmod/internal/config"
	"io"
	"net/http"
	"time"
)

type DrugsComClient struct {
	baseUrl   string
	caCerts   *x509.CertPool
	rootCerts *x509.CertPool
}

func NewDrugsClient(baseUrl string) (*DrugsComClient, error) {
	caCerts, err := config.NewCerts(false)
	if err != nil {
		return nil, err
	}
	rootCa, err := config.NewCerts(true)
	if err != nil {
		return nil, err
	}

	return &DrugsComClient{
		baseUrl:   baseUrl,
		caCerts:   caCerts,
		rootCerts: rootCa,
	}, nil
}

func (client DrugsComClient) GetData(c context.Context, param string) ([]byte, error) {

	url := client.baseUrl + param

	tlsConfig := &tls.Config{
		RootCAs: client.rootCerts,
	}
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
		Timeout: 120 * time.Second,
	}

	response, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
