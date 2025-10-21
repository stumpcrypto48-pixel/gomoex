package config

import (
	"testing"
)

func TestConfig(t *testing.T) {

	t.Run("Test CA creation", func(t *testing.T) {
		_, err := NewCerts(true)
		if err != nil {
			t.Errorf("error while creating certs :: %v", err)
		}
	})

}
