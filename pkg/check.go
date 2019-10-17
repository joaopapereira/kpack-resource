package pkg

import (
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/pivotal/kpack/pkg/client/clientset/versioned"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"kpack-resource/pkg/config"
	"kpack-resource/pkg/image"
)

type Check struct {
	KpackClient versioned.Interface
}

func (c *Check) Check(source config.Source, latestVersion config.Version) (config.Version, error) {
	imageName := image.UniqueName(source.Image)
	img, err := c.KpackClient.BuildV1alpha1().Images(source.Namespace).Get(imageName, v1.GetOptions{})
	if errors.IsNotFound(err) {
		return latestVersion, nil
	} else if err != nil {
		return latestVersion, err
	}

	if img.Status.LatestImage == "" {
		return latestVersion, nil
	}
	imgRef, err := name.ParseReference(img.Status.LatestImage)
	if err != nil {
		return latestVersion, err
	}

	return config.Version{
		Digest: imgRef.Identifier(),
		BuildNumber: img.Status.BuildCounter,
	}, nil
}
