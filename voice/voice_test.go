package voice

import (
	"log"
	"testing"
)

func TestListVoices(t *testing.T) {
	voices, err := ListVoices()
	if err != nil {
		t.Fatal(err)
	}
	log.Println(voices)
}
