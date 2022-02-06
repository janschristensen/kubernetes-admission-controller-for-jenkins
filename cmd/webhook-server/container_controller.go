package main

import (
	"strings"
)

func containerImageISOK(image string, conf *imageConf) bool {
	for _, reg := range conf.Allow.Registries {
		if strings.Contains(image, reg) {
			return true
		}
	}

	//TODO: implement option conf.Allow.Empty=="true" and  conf.Allow.Empty=="mutate". Both allow empty.
	return false
}
