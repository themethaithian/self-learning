package config

import "testing"

func TestDeliberatelyFailingForCIProof(t *testing.T) {
	t.Fatal("intentional failure: proving CI turns red")
}
