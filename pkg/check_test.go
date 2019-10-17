package pkg_test

import (
	"testing"

	"github.com/pivotal/kpack/pkg/apis/build/v1alpha1"
	"github.com/pivotal/kpack/pkg/client/clientset/versioned/fake"
	"github.com/sclevine/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1"

	"kpack-resource/pkg"
	"kpack-resource/pkg/config"
	"kpack-resource/pkg/image"
)

func TestCheck(t *testing.T) {
	spec.Run(t, "Test Resource Check", testCheck)
}

func testCheck(t *testing.T, when spec.G, it spec.S) {
	const imageTag = "some/image:tag"
	const namespace = "some-namespace"
	var (
		kpackClient = fake.NewSimpleClientset()
		subject     = pkg.Check{KpackClient: kpackClient}
		imageName   = image.UniqueName(imageTag)
	)
	it("returns no version when image not configured in kpack", func() {
		newVersion, err := subject.Check(config.Source{}, config.Version{})
		require.NoError(t, err)

		assert.Equal(t, config.Version{}, newVersion)
	})

	when("image exists", func() {
		it("returns no version when no builds as terminated successfully for the image", func() {
			_, err := kpackClient.BuildV1alpha1().Images(namespace).Create(&v1alpha1.Image{
				ObjectMeta: v1.ObjectMeta{
					Name: imageName,
				},
				Spec: v1alpha1.ImageSpec{
					Tag: imageTag,
				},
			})
			require.NoError(t, err)

			newVersion, err := subject.Check(config.Source{
				Image: imageTag,
			}, config.Version{})
			require.NoError(t, err)

			assert.Equal(t, config.Version{}, newVersion)
		})

		it("returns new version when a builds as terminated successfully for the image", func() {
			_, err := kpackClient.BuildV1alpha1().Images(namespace).Create(&v1alpha1.Image{
				ObjectMeta: v1.ObjectMeta{
					Name: imageName,
				},
				Spec: v1alpha1.ImageSpec{
					Tag: imageTag,
				},
				Status: v1alpha1.ImageStatus{
					LatestImage:    imageTag + "@sha256:3c515511cbaa891ab88088b1a12a56cedcd52a9e4a8868380c418a7619c467e0",
					BuildCounter: 10,
				},
			})
			require.NoError(t, err)

			newVersion, err := subject.Check(config.Source{
				Image: imageTag,
			}, config.Version{})
			require.NoError(t, err)

			assert.Equal(t, config.Version{
				Digest: "sha256:3c515511cbaa891ab88088b1a12a56cedcd52a9e4a8868380c418a7619c467e0",
				BuildNumber: 10,
			}, newVersion)
		})
	})
}
