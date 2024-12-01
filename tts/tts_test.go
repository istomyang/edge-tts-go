package tts_test

import (
	"edge-tts-go/tts"
	"os"
	"testing"
)

func TestTts(t *testing.T) {
	data, err := tts.Request(
		"en-US",
		"Microsoft Server Speech Text to Speech Voice (en-US, JennyNeural)",
		"+0Hz",
		"+25%",
		"+0%",
		"audio-24khz-48kbitrate-mono-mp3",
		"In February of 2016 I began to experience two separate realities at the same time.",
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(data) == 0 {
		t.Fatal("data len is 0")
	}

	_ = os.WriteFile("text.mp3", data, 0644)
}
