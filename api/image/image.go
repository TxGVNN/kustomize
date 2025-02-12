// Copyright 2020 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package image

import (
	"regexp"
	"strings"
)

// IsImageMatched returns true if the value of t is identical to the
// image name in the full image name and tag as given by s.
func IsImageMatched(s, t string) bool {
	// Tag values are limited to [a-zA-Z0-9_.{}-].
	// Some tools like Bazel rules_k8s allow tag patterns with {} characters.
	// More info: https://github.com/bazelbuild/rules_k8s/pull/423
	pattern, _ := regexp.Compile("^" + t + "(:[a-zA-Z0-9_.{}-]+)?(@sha256:[a-zA-Z0-9_.{}-]+)?$")
	return pattern.MatchString(s)
}

// Split separates and returns the name and tag parts
// from the image string using either colon `:` or at `@` separators.
// Note that the returned tag keeps its separator.
func Split(imageName string) (name string, tag string) {
	ic := -1
	ia := strings.LastIndex(imageName, "@")
	// check if image name contains a domain
	// if domain is present, ignore domain and check for `:`
	if slashIndex := strings.Index(imageName, "/"); slashIndex < 0 {
		if ia > 0 {
			ic = strings.LastIndex(imageName[:ia], ":")
		} else {
			ic = strings.LastIndex(imageName, ":")
		}

	} else {
		lastIc := -1
		if ia > 0 {
			lastIc = strings.LastIndex(imageName[slashIndex:ia], ":")
		} else {
			lastIc = strings.LastIndex(imageName[slashIndex:], ":")
		}
		// set ic only if `:` is present
		if lastIc > 0 {
			ic = slashIndex + lastIc
		}
	}

	if ic < 0 && ia < 0 {
		return imageName, ""
	}
	i := ic
	if ic < 0 || (ia > 0 && ia < ic) {
		i = ia
	}

	name = imageName[:i]
	tag = imageName[i:]
	return
}
