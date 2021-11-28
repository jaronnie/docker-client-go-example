package util

import (
	yaml2 "k8s.io/apimachinery/pkg/util/yaml"
)

func GetContainerIDByMachineUUID(uuid string) (string, error) {
	return "", nil
}

func Yaml2Json(content []byte) ([]byte, error) {
	data, err := yaml2.ToJSON(content)
	if err != nil {
		return nil, err
	}
	return data, nil
}


