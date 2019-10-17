package image

import (
	"crypto/md5"
	"encoding/hex"
	"regexp"
)

const imagePrefix = "i-"

func UniqueName(imageName string) string {
	last15 := last15(imageName)

	reg, err := regexp.Compile("[^A-Za-z0-9_-]+")
	if err != nil {
		return ""
	}
	escaped := reg.ReplaceAllString(last15, "-")

	hasher := md5.New()
	hasher.Write([]byte(imageName))
	sha := hex.EncodeToString(hasher.Sum(nil))

	return imagePrefix + escaped + "-" + sha
}

func last15(imageName string) string {
	if len(imageName) < 15 {
		return imageName
	}

	return imageName[len(imageName)-15:]
}
