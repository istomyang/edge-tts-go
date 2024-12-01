package internal_test

import (
	"edge-tts-go/internal"
	"testing"
)

func TestGenSecMsGec(t *testing.T) {
	token := internal.GenSecMsGec()
	t.Log(token)
}
