package main

import (
	"github.com/cloudboss/ofcourse/ofcourse"

	"kpack-resource/pkg/resource"
)

func main() {
	ofcourse.Check(&resource.Resource{})
}
