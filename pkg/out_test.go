package pkg_test

import (
	"testing"

	"github.com/cloudboss/ofcourse/ofcourse"
	"github.com/pivotal/kpack/pkg/apis/build/v1alpha1"
	"github.com/pivotal/kpack/pkg/client/clientset/versioned/fake"
	"github.com/sclevine/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"kpack-resource/pkg"
	"kpack-resource/pkg/config"
	"kpack-resource/pkg/git"
	"kpack-resource/pkg/image"
)

func TestOut(t *testing.T) {
	spec.Run(t, "Test Resource Out", testOut)
}

func testOut(t *testing.T, when spec.G, it spec.S) {
	const imageTag = "some/image:tag"
	const namespace = "some-namespace"
	const serviceAccount = "some-service-account"
	var (
		kpackClient = fake.NewSimpleClientset()
		fakeGitInfo = &fakeGitInfo{}
		subject     = pkg.Out{
			KpackClient: kpackClient,
			Logger:      ofcourse.NewLogger("warn"),
			GitInfo:     fakeGitInfo,
		}
		imageName = image.UniqueName(imageTag)
		source    = config.Source{
			Image:          imageTag,
			Kpack:          config.Kpack{},
			Namespace:      namespace,
			ServiceAccount: serviceAccount,
		}
	)

	when("want to put a blob", func() {
		it("creates a new build", func() {
			params := config.OutParams{
				BlobPath: "/some/path/blob.jar",
			}
			_, err := subject.Put("", source, params)
			require.NoError(t, err)
			img, err := kpackClient.BuildV1alpha1().Images(namespace).Get(imageName, v1.GetOptions{})
			require.NoError(t, err)

			assert.Equal(t, v1alpha1.ImageSpec{
				Tag: imageTag,
				Builder: v1alpha1.ImageBuilder{
					TypeMeta: v1.TypeMeta{
						Kind: v1alpha1.ClusterBuilderKind,
					},
					Name: "cluster-sample-builder",
				},
				ServiceAccount: serviceAccount,
				Source: v1alpha1.SourceConfig{
					Blob: &v1alpha1.Blob{URL: "/some/path/blob.jar"},
				},
				ImageTaggingStrategy: v1alpha1.BuildNumber,
			}, img.Spec)
		})

		it("updates existing build", func() {
			_, err := kpackClient.BuildV1alpha1().Images(namespace).Create(&v1alpha1.Image{
				ObjectMeta: v1.ObjectMeta{
					Name:      imageName,
					Namespace: namespace,
				},
				Spec: v1alpha1.ImageSpec{
					Tag: source.Image,
					Builder: v1alpha1.ImageBuilder{
						TypeMeta: v1.TypeMeta{
							Kind: v1alpha1.ClusterBuilderKind,
						},
						Name: "cluster-sample-builder",
					},
					ServiceAccount: source.ServiceAccount,
					Source: v1alpha1.SourceConfig{
						Blob: &v1alpha1.Blob{
							URL: "/some/other/blob.jar",
						},
					},
					ImageTaggingStrategy: v1alpha1.BuildNumber,
				},
			})
			require.NoError(t, err)

			params := config.OutParams{
				BlobPath: "/some/path/blob.jar",
			}
			source.ServiceAccount = "some-other-service-account"
			_, err = subject.Put("", source, params)
			require.NoError(t, err)

			img, err := kpackClient.BuildV1alpha1().Images(namespace).Get(imageName, v1.GetOptions{})
			require.NoError(t, err)

			assert.Equal(t, v1alpha1.ImageSpec{
				Tag: imageTag,
				Builder: v1alpha1.ImageBuilder{
					TypeMeta: v1.TypeMeta{
						Kind: v1alpha1.ClusterBuilderKind,
					},
					Name: "cluster-sample-builder",
				},
				ServiceAccount: "some-other-service-account",
				Source: v1alpha1.SourceConfig{
					Blob: &v1alpha1.Blob{
						URL: "/some/path/blob.jar",
					},
				},
				ImageTaggingStrategy: v1alpha1.BuildNumber,
			}, img.Spec)
		})
	})

	when("want to put with a git repository", func() {
		it("creates a new build", func() {
			fakeGitInfo.fromPathReturns = git.Git{
				Repository: "github.com/some/repo",
				Commit:     "ad897234",
			}
			params := config.OutParams{
				Git: "/some/path/to/repo",
			}
			_, err := subject.Put("", source, params)
			require.NoError(t, err)
			img, err := kpackClient.BuildV1alpha1().Images(namespace).Get(imageName, v1.GetOptions{})
			require.NoError(t, err)

			assert.Equal(t, v1alpha1.ImageSpec{
				Tag: imageTag,
				Builder: v1alpha1.ImageBuilder{
					TypeMeta: v1.TypeMeta{
						Kind: v1alpha1.ClusterBuilderKind,
					},
					Name: "cluster-sample-builder",
				},
				ServiceAccount: serviceAccount,
				Source: v1alpha1.SourceConfig{
					Git: &v1alpha1.Git{
						URL:      "github.com/some/repo",
						Revision: "ad897234",
					},
				},
				ImageTaggingStrategy: v1alpha1.BuildNumber,
			}, img.Spec)
		})

		it("updates existing build", func() {
			fakeGitInfo.fromPathReturns = git.Git{
				Repository: "github.com/some/repo",
				Commit:     "ad897234",
			}
			params := config.OutParams{
				Git: "/some/path/to/repo",
			}

			_, err := kpackClient.BuildV1alpha1().Images(namespace).Create(&v1alpha1.Image{
				ObjectMeta: v1.ObjectMeta{
					Name:      imageName,
					Namespace: namespace,
				},
				Spec: v1alpha1.ImageSpec{
					Tag: source.Image,
					Builder: v1alpha1.ImageBuilder{
						TypeMeta: v1.TypeMeta{
							Kind: v1alpha1.ClusterBuilderKind,
						},
						Name: "cluster-sample-builder",
					},
					ServiceAccount: source.ServiceAccount,
					Source: v1alpha1.SourceConfig{
						Blob: &v1alpha1.Blob{
							URL: "/some/other/blob.jar",
						},
					},
					ImageTaggingStrategy: v1alpha1.BuildNumber,
				},
			})
			require.NoError(t, err)

			source.ServiceAccount = "some-other-service-account"
			_, err = subject.Put("", source, params)
			require.NoError(t, err)

			img, err := kpackClient.BuildV1alpha1().Images(namespace).Get(imageName, v1.GetOptions{})
			require.NoError(t, err)

			assert.Equal(t, v1alpha1.ImageSpec{
				Tag: imageTag,
				Builder: v1alpha1.ImageBuilder{
					TypeMeta: v1.TypeMeta{
						Kind: v1alpha1.ClusterBuilderKind,
					},
					Name: "cluster-sample-builder",
				},
				ServiceAccount: "some-other-service-account",
				Source: v1alpha1.SourceConfig{
					Git: &v1alpha1.Git{
						URL:      "github.com/some/repo",
						Revision: "ad897234",
					},
				},
				ImageTaggingStrategy: v1alpha1.BuildNumber,
			}, img.Spec)
		})
	})
}

type fakeGitInfo struct {
	fromPathReturns git.Git
}

func (f fakeGitInfo) FromPath(path string) (git.Git, error) {
	return f.fromPathReturns, nil
}
