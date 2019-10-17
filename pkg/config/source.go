package config

import (
	"errors"
	"strconv"

	"github.com/cloudboss/ofcourse/ofcourse"
)

func NewSource(source map[string]interface{}) (Source, error) {
	return Source{
		Image: fromSource(source, "image"),
		Kpack: Kpack{
			Token:  fromSource(source, "k8s_token"),
			Domain: fromSource(source, "k8s_api"),
			CaCert: fromSource(source, "k8s_ca_cert"),
		},
		ServiceAccount: fromSource(source, "service_account"),
		Namespace:      fromSource(source, "namespace"),
		Username:       fromSource(source, "username"),
		Password:       fromSource(source, "password"),
	}, nil
}

type Source struct {
	Image          string `json:"image"`
	Kpack          Kpack  `json:"kpack"`
	Namespace      string `json:"namespace"`
	ServiceAccount string `json:"serviceAccount"`
	Username       string `json:"username"`
	Password       string `json:"password"`
}

func (s Source) Valid() (bool, error) {
	if s.Image == "" {
		return false, errors.New("image is a mandatory field")
	}

	return true, nil
}

func NewOutParams(params ofcourse.Params) OutParams {
	return OutParams{
		SourcePath: fromSource(params, "source_path"),
		Git:        fromSource(params, "git_path"),
		BlobPath:   fromSource(params, "blob_path"),
	}
}

type OutParams struct {
	SourcePath string `json:"source_path"`
	Git        string `json:"git_file"`
	BlobPath   string `json:"blob_path"`
}

type GitFile struct {
	URL      string `json:"url"`
	Revision string `json:"revision"`
}

type Kpack struct {
	Token  string `json:"token"`
	Domain string `json:"network_address"`
	CaCert string `json:"cacert"`
}

func NewVersion(version map[string]string) Version {
	buildNumber, _ := strconv.Atoi(fromVersion(version, "buildNumber"))
	return Version{
		Digest:      fromVersion(version, "digest"),
		BuildNumber: int64(buildNumber),
	}
}

type Version struct {
	Digest      string `json:"digest"`
	BuildNumber int64  `json:"buildNumber"`
}

func (v Version) ToMap() map[string]string {
	if v.Digest == "" {
		return nil
	}
	return map[string]string{
		"digest":      v.Digest,
		"buildNumber": strconv.Itoa(int(v.BuildNumber)),
	}
}

func fromSource(source map[string]interface{}, key string) string {
	return fromSourceDefault(source, key, "")
}

func fromSourceDefault(source map[string]interface{}, key, defaultValue string) string {
	if result, ok := source[key]; !ok {
		return defaultValue
	} else {
		return result.(string)
	}
}

func fromVersion(version map[string]string, key string) string {
	return fromVersionDefault(version, key, "")
}

func fromVersionDefault(version map[string]string, key, defaultValue string) string {
	if result, ok := version[key]; !ok {
		return defaultValue
	} else {
		return result
	}
}
