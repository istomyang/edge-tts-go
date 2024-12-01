package tts

import (
	"edge-tts-go/internal"
	"encoding/binary"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Request a text to Microsoft using Edge browser's speech.
func Request(
	lang string,
	voice string,
	pitch string,
	rate string,
	volume string,
	outputFormat string,
	text string,
) (audio []byte, err error) {
	// You can use a packet grabbing tool, like Charles,
	// to grab how the Edge browser handles audio read aloud
	// and you'll get the flow of operations.

	// The reason for Microsoft using websocket is that
	// the audio data is sent in a continuous stream.
	var conn *websocket.Conn
	if conn, err = buildConnection(); err != nil {
		err = fmt.Errorf("build connection error: %s", err)
		return
	}
	defer conn.Close()

	if err = conn.WriteMessage(websocket.TextMessage, []byte(request1(outputFormat))); err != nil {
		return
	}

	if err = conn.WriteMessage(websocket.TextMessage, []byte(request2(lang, voice, pitch, rate, volume, text))); err != nil {
		return
	}

	for {
		var (
			messageType int

			// Because AI generates Audio Token continuously,
			// the client needs to accept a block of data continuously.
			data []byte
		)
		if messageType, data, err = conn.ReadMessage(); err != nil {
			err = fmt.Errorf("read message error: %s", err)
			continue
		}

		switch messageType {
		case websocket.TextMessage:
			if strings.Contains(string(data), "Path:turn.end") {
				return
			}
		case websocket.BinaryMessage:
			// metadata data length.
			len := uint16(binary.BigEndian.Uint16(data[:2]))
			// You can use `data[2:len]` to get metadata.
			audio = append(audio, data[len+2:]...)
		default:
			return
		}
	}
}

func buildConnection() (conn *websocket.Conn, err error) {
	var header = make(http.Header)
	header.Set("Pragma", "no-cache")
	header.Set("Cache-Control", "no-cache")
	header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 Edg/131.0.0.0")
	header.Set("Origin", "chrome-extension://jdiccldimpdaibmpdkjnbmckianbfold")
	header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")

	conn, _, err = websocket.DefaultDialer.Dial(buildWsUrl(), header)
	return
}

func request1(outputFormat string) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("X-Timestamp:%s\r\n", datetime2String()))
	b.WriteString("Content-Type:application/json; charset=utf-8\r\n")
	b.WriteString("Path:speech.config\r\n\r\n")
	b.WriteString(fmt.Sprintf(`{"context":{"synthesis":{"audio":{"metadataoptions":{"sentenceBoundaryEnabled":"false","wordBoundaryEnabled":"false"},"outputFormat":"%s"}}}}`, outputFormat))
	return b.String()
}

func request2(lang string, voice string, pitch string, rate string, volume string, text string) string {
	u, _ := uuid.NewUUID()
	token := strings.ReplaceAll(u.String(), "-", "")
	var b strings.Builder
	b.WriteString(fmt.Sprintf("X-RequestId:%s\r\n", token))
	b.WriteString("Content-Type:application/ssml+xml\r\n")
	b.WriteString(fmt.Sprintf("X-Timestamp:%s\r\n", datetime2String()))
	b.WriteString("Path:ssml\r\n\r\n")
	b.WriteString(fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?><speak version="1.0" xmlns="http://www.w3.org/2001/10/synthesis"  xml:lang="%s"><voice name="%s"><prosody pitch="%s" rate="%s" volume="%s">%s</prosody></voice></speak>`, lang, voice, pitch, rate, volume, text))
	return b.String()
}

func buildWsUrl() string {
	u, _ := uuid.NewUUID()
	return fmt.Sprintf(
		"wss://speech.platform.bing.com/consumer/speech/synthesize/readaloud/edge/v1?TrustedClientToken=6A5AA1D4EAFF4E9FB37E23D68491D6F4&Sec-MS-GEC=%s&Sec-MS-GEC-Version=1-131.0.2903.51&ConnectionId=%s",
		internal.GenSecMsGec(),
		u.String(),
	)
}

// datetime2String create like this:
//
// ```text
// Sun Feb  3 2018 12:34:56 GMT+0800 (中国标准时间)
// ```
func datetime2String() string {
	now := time.Now().UTC()
	timezone := time.FixedZone("GMT+8", 8*3600)
	datetime := now.In(timezone)
	format := "Mon Jan 02 2006 15:04:05 GMT-07:00 (Z07:00)"
	return datetime.Format(format)
}
