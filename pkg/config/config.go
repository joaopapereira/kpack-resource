package config

import (
	"k8s.io/client-go/rest"
)

func RetrieveLocalConfiguration(apiHost, token, cacert string) (*rest.Config, error) {
	return &rest.Config{
		Host:                apiHost,
		BearerToken:         token,
		BearerTokenFile:     "",
		TLSClientConfig: rest.TLSClientConfig{
			Insecure:   false,
			CAData:     []byte(cacert),
		},
	}, nil
}
