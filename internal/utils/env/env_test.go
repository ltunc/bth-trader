package env

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
)

func TestGet(t *testing.T) {
	expected := fmt.Sprintf("SomeValue_%d", rand.Int())
	key := "TEST_VALUE_GET"
	if err := os.Setenv(Prefix+key, expected); err != nil {
		t.Error(err)
	}
	fallback := fmt.Sprintf("DefaultValue_%d", rand.Int())
	if actual := Get(key, fallback); expected != actual {
		t.Errorf("Expected '%s', got '%s'", expected, actual)
	}
	if actual := Get("WRONG_KEY", fallback); fallback != actual {
		t.Errorf("Expected fallback value '%s', got '%s'", expected, actual)
	}
}
