package pkg

import (
	"io/ioutil"
	"path/filepath"

	"github.com/cloudboss/ofcourse/ofcourse"
	"github.com/fatih/color"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/tarball"

	"kpack-resource/pkg/config"
)

type In struct {
	Logger *ofcourse.Logger
}

func (i *In) Fetch(dest string, source config.Source, version config.Version) error {
	if version.Digest == "" {
		i.Logger.Errorf("Image not built at this point in time: %s", source.Image)
		return nil
	}
	ref := source.Image + "@" + version.Digest

	imgRef, err := name.ParseReference(ref, name.WeakValidation)
	if err != nil {
		i.Logger.Errorf("failed to resolve name: %s", err)
		return err
	}

	i.Logger.Infof("fetching %s@%s\n", color.GreenString(source.Image), color.YellowString(version.Digest))

	auth := &authn.Basic{
		Username: source.Username,
		Password: source.Password,
	}

	var imageOpts []remote.Option

	if auth.Username != "" && auth.Password != "" {
		imageOpts = append(imageOpts, remote.WithAuth(auth))
	}

	image, err := remote.Image(imgRef, imageOpts...)
	if err != nil {
		i.Logger.Errorf("failed to locate remote image: %s", err)
		return err
	}

	tag, err := name.NewTag(imgRef.Name(), name.WeakValidation)
	if err != nil {
		i.Logger.Errorf("failed to construct tag reference: %s", err)
		return err
	}

	err = tarball.WriteToFile(filepath.Join(dest, "image.tar"), tag, image)
	if err != nil {
		i.Logger.Errorf("failed to write OCI image: %s", err)
		return err
	}

	err = ioutil.WriteFile(filepath.Join(dest, "tag"), []byte(tag.TagStr()), 0644)
	if err != nil {
		i.Logger.Errorf("failed to save image tag: %s", err)
		return err
	}

	return nil
}
