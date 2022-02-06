package main

import (
	"testing"
)

func TestVerifyContainerImage(t *testing.T) {
	name := "bogusImageNameThatShouldFail"
	want := false
	conf, _ := readConf("image_conf.yaml")

	if !want == containerImageISOK(name, conf) {
		t.Fatalf(`%q should fail`, name)
	}
}
