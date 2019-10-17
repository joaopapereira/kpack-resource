package git_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/cloudboss/ofcourse/ofcourse"
	"github.com/sclevine/spec"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"

	"kpack-resource/pkg/git"
)

func TestGitInfo(t *testing.T) {
	spec.Run(t, "Test Info", testGitInfo)
}

func testGitInfo(t *testing.T, when spec.G, it spec.S) {
	var (
		gitDir  string
		subject = git.InfoRetriever{
			Logger: ofcourse.NewLogger("warn"),
		}
	)

	it.Before(func() {
		var err error
		gitDir, err = ioutil.TempDir("", "git_dir")
		require.NoError(t, err)

		createRepository(t, gitDir, false)
		addRemote(t, gitDir, "https://github.com/some/path")
	})

	it.After(func() {
		require.NoError(t, os.RemoveAll(gitDir))
	})

	it("returns error when cannot read folder", func() {
		result, err := subject.FromPath("/path/not/found")
		require.EqualError(t, err, "reading git folder: repository does not exist")
		assert.Equal(t, git.Git{}, result)
	})

	it("returns the correct information", func() {
		commitSha := doCommit(t, gitDir)
		result, err := subject.FromPath(gitDir)
		require.NoError(t, err)
		assert.Equal(t, git.Git{
			Repository: "https://github.com/some/path",
			Commit:     commitSha,
		}, result)
	})
}

func createRepository(t *testing.T, dir string, isBare bool) string {
	var cmd *exec.Cmd
	if isBare {
		cmd = exec.Command("git", "init", "--bare", dir)
	} else {
		cmd = exec.Command("git", "init", dir)
	}
	err := cmd.Run()
	require.NoError(t, err)

	return dir
}

func addRemote(t *testing.T, local, remote string) {
	cmd := exec.Command("git", "remote", "add", "origin", remote)
	cmd.Dir = local
	err := cmd.Run()
	require.NoError(t, err)
}

func doCommit(t *testing.T, dir string) string {
	imageConfig, err := ioutil.TempFile(dir, "somefile.go")
	defer os.Remove(imageConfig.Name())
	require.NoError(t, err)

	_, err = imageConfig.WriteString("some string")
	require.NoError(t, err)
	cmd := exec.Command("git", "add", ".")
	cmd.Dir = dir
	output, err := cmd.Output()
	t.Log(string(output))
	cmd = exec.Command("git", "commit", "-am", "\"some message\"")
	cmd.Dir = dir
	output, err = cmd.Output()
	require.NoError(t, err)
	cmd = exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = dir
	output, err = cmd.Output()
	require.NoError(t, err)
	return strings.ReplaceAll(string(output), "\n", "")
}
