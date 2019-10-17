package resource

import (
	"errors"
	"fmt"

	oc "github.com/cloudboss/ofcourse/ofcourse"
	"github.com/pivotal/kpack/pkg/client/clientset/versioned"

	"kpack-resource/pkg"
	"kpack-resource/pkg/config"
	"kpack-resource/pkg/git"
)

type Resource struct{}

func (r *Resource) Check(source oc.Source, version oc.Version, env oc.Environment,
	logger *oc.Logger) ([]oc.Version, error) {
	resourceSource, err := config.NewSource(source)
	if err != nil {
		logger.Errorf("validating source: %s", err)
		return nil, err
	}

	kpackClient, err := initk8s(resourceSource)
	if err != nil {
		logger.Errorf("create kpack client: %s", err)
		return nil, err
	}

	checker := &pkg.Check{
		KpackClient: kpackClient,
	}
	newVersion, err := checker.Check(resourceSource, config.NewVersion(version))
	if err != nil {
		logger.Errorf("error retrieving current version: %s", err.Error())
		return nil, err
	}

	if newVersion.Digest == "" {
		return nil, errors.New("no images were built using this resource")
	}

	versions := []oc.Version{newVersion.ToMap()}

	return versions, nil
}

func (r *Resource) In(outputDirectory string, source oc.Source, params oc.Params, version oc.Version,
	env oc.Environment, logger *oc.Logger) (oc.Version, oc.Metadata, error) {
	logger.Debugf("Starting to get image")
	resourceSource, err := config.NewSource(source)
	if err != nil {
		logger.Errorf("validating source: %s", err)
		return nil, nil, err
	}
	localVersion := config.NewVersion(version)

	in := pkg.In{
		Logger: logger,
	}
	err = in.Fetch(outputDirectory, resourceSource, localVersion)
	if err != nil {
		logger.Errorf("Failed to get docker image: %s", err.Error())
		return nil, nil, err
	}
	metadata := oc.Metadata{}
	return version, metadata, nil
}

func (r *Resource) Out(inputDirectory string, source oc.Source, params oc.Params,
	env oc.Environment, logger *oc.Logger) (oc.Version, oc.Metadata, error) {
	logger.Debugf("Starting to create/update image")
	resourceSource, err := config.NewSource(source)
	if err != nil {
		logger.Errorf("validating source: %s", err)
		return nil, nil, err
	}

	kpackClient, err := initk8s(resourceSource)
	if err != nil {
		logger.Errorf("create kpack client: %s", err)
		return nil, nil, err
	}

	out := pkg.Out{
		KpackClient: kpackClient,
		Logger:      logger,
		GitInfo: git.InfoRetriever{
			Logger: logger,
		},
	}
	image, err := out.Put(inputDirectory, resourceSource, config.NewOutParams(params))

	if err != nil {
		logger.Errorf("Failed to create/update image: %s", err.Error())
		return nil, nil, err
	}

	metadata := oc.Metadata{
		oc.NameVal{
			Name:  "imageResourceName",
			Value: image.Name,
		},
	}

	return nil, metadata, nil
}

func initk8s(cfg config.Source) (versioned.Interface, error) {
	k8sConfig, err := config.RetrieveLocalConfiguration(cfg.Kpack.Domain, cfg.Kpack.Token, cfg.Kpack.CaCert)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve k8s configuration: %s", err)
	}

	kpackClient, err := versioned.NewForConfig(k8sConfig)
	if err != nil {
		return nil, fmt.Errorf("could not get kpack clientset: %s", err.Error())
	}

	return kpackClient, err
}
