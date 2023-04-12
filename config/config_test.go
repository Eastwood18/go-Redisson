package config

import (
	"fmt"
	"testing"
)

func TestFromYaml(t *testing.T) {
	c := new(Config)
	err := c.FromYaml("test.yaml")
	fmt.Println(c, err)
}
