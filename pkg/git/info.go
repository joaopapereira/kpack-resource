package git

import (
	"github.com/cloudboss/ofcourse/ofcourse"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
)

type Git struct {
	Repository string
	Commit     string
}

type InfoRetriever struct {
	Logger *ofcourse.Logger
}

func (i InfoRetriever) FromPath(path string) (Git, error) {
	i.Logger.Debugf("Checking git repo on: %s", path)
	repository, err := git.PlainOpen(path)
	if err != nil {
		return Git{}, errors.Wrap(err, "reading git folder")
	}
	head, err := repository.Head()
	if err != nil {
		return Git{}, errors.Wrap(err, "reading git HEAD")
	}

	remote, err := repository.Remote("origin")
	if err != nil {
		return Git{}, errors.Wrap(err, "retrieving 'origin' remote")
	}

	return Git{
		Repository: remote.Config().URLs[0],
		Commit:     head.Hash().String(),
	}, nil
}
