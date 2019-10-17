package pkg

import (
	"path/filepath"

	"github.com/cloudboss/ofcourse/ofcourse"
	"github.com/pivotal/kpack/pkg/apis/build/v1alpha1"
	"github.com/pivotal/kpack/pkg/client/clientset/versioned"
	"github.com/pkg/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"kpack-resource/pkg/config"
	"kpack-resource/pkg/git"
	"kpack-resource/pkg/image"
)

type GitInfoRetrieve interface {
	FromPath(path string) (git.Git, error)
}

type Out struct {
	KpackClient versioned.Interface
	Logger      *ofcourse.Logger
	GitInfo     GitInfoRetrieve
}

func (o *Out) Put(inputDirectory string, source config.Source, params config.OutParams) (v1alpha1.Image, error) {
	imgName := image.UniqueName(source.Image)

	sourceConfig, _ := o.generateSourceConfig(inputDirectory, params)

	o.Logger.Debugf("Check image '%s' on namespace: '%s'", imgName, source.Namespace)
	img, err := o.KpackClient.BuildV1alpha1().Images(source.Namespace).Get(imgName, v1.GetOptions{})
	if err != nil && !k8serrors.IsNotFound(err) {
		return v1alpha1.Image{}, errors.Wrap(err, "retrieving the image")
	} else if k8serrors.IsNotFound(err) {
		o.Logger.Debugf("Image not found, going to create a new one")
		createdImage, err := o.KpackClient.BuildV1alpha1().Images(source.Namespace).Create(&v1alpha1.Image{
			ObjectMeta: v1.ObjectMeta{
				Name:      imgName,
				Namespace: source.Namespace,
				Annotations: map[string]string{
					"kpack-resource.joaopapereira.io/managed": "concourse",
				},
			},
			Spec: v1alpha1.ImageSpec{
				Tag: source.Image,
				Builder: v1alpha1.ImageBuilder{
					TypeMeta: v1.TypeMeta{
						Kind: v1alpha1.ClusterBuilderKind,
					},
					Name: "cluster-sample-builder",
				},
				ServiceAccount:       source.ServiceAccount,
				Source:               sourceConfig,
				ImageTaggingStrategy: v1alpha1.BuildNumber,
			},
		})
		if err != nil {
			return v1alpha1.Image{}, errors.Wrap(err, "creating an image")
		}
		return *createdImage, nil
	}
	o.Logger.Debugf("Updating existing image")

	img.Spec.Source = sourceConfig
	img.Spec.ServiceAccount = source.ServiceAccount

	updatedImage, err := o.KpackClient.BuildV1alpha1().Images(source.Namespace).Update(img)
	if err != nil {
		return v1alpha1.Image{}, errors.Wrap(err, "updating an image")
	}

	return *updatedImage, nil
}

func (o *Out) generateSourceConfig(inputDirectory string, params config.OutParams) (v1alpha1.SourceConfig, error) {
	if params.BlobPath != "" {
		o.Logger.Infof("Generating Image from Blob at: %s", params.BlobPath)
		return v1alpha1.SourceConfig{
			Blob: &v1alpha1.Blob{URL: params.BlobPath},
		}, nil
	} else if params.Git != "" {
		info, err := o.GitInfo.FromPath(filepath.Join(inputDirectory, params.Git))
		if err != nil {
			return v1alpha1.SourceConfig{}, err
		}

		o.Logger.Infof("Generating Image from Git repository: %s on the commit %s", info.Repository, info.Commit)
		return v1alpha1.SourceConfig{
			Git: &v1alpha1.Git{
				URL:      info.Repository,
				Revision: info.Commit,
			},
		}, nil
	}
	return v1alpha1.SourceConfig{}, nil
}
