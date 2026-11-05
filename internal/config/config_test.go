package config

import "testing"

func TestDefaultConfigurationValidates(t *testing.T) {
	c := Default()
	if err := Validate(&c); err != nil {
		t.Fatal(err)
	}
}
