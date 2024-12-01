package voice

import (
	"edge-tts-go/internal"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// / List all voices from Microsoft.
// /
// / You can use its name to [`crate::request`] voice name.
func ListVoices() (r []Voice, err error) {
	req, _ := http.NewRequest(http.MethodGet, buildUrl(), nil)

	req.Header.Set("sec-ch-ua-platform", "macOS")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 Edg/131.0.0.0")
	req.Header.Set("sec-ch-ua", `"Microsoft Edge";v="131", "Chromium";v="131", "Not_A Brand";v="24"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("X-Edge-Shopping-Flag", "1")
	req.Header.Set("Sec-MS-GEC", internal.GenSecMsGec())
	req.Header.Set("Sec-MS-GEC-Version", "1-131.0.2903.70")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")

	var resp *http.Response
	if resp, err = http.DefaultClient.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("request error: %s", resp.Status)
		return
	}

	var body []byte
	if body, err = io.ReadAll(resp.Body); err != nil {
		return
	}

	if err = json.Unmarshal(body, &r); err != nil {
		return
	}

	return
}

func buildUrl() string {
	return fmt.Sprintf(
		"https://speech.platform.bing.com/consumer/speech/synthesize/readaloud/voices/list?trustedclienttoken=6A5AA1D4EAFF4E9FB37E23D68491D6F4&Sec-MS-GEC=%s&Sec-MS-GEC-Version=1-131.0.2903.70", internal.GenSecMsGec())
}

// pub fn list_voices() -> Result<Vec<Voice>, Error> {
//     let client = reqwest::blocking::Client::new();
//     let res = client.get(build_url())

//     if res.status() != StatusCode::OK {
//         return Err(Error::from(res.text().unwrap()));
//     }

//     res.json().map_err(|e| Error::from(e.to_string()))
// }

// fn build_url() -> String {
//     format!("https://speech.platform.bing.com/consumer/speech/synthesize/readaloud/voices/list?trustedclienttoken=6A5AA1D4EAFF4E9FB37E23D68491D6F4&Sec-MS-GEC={}&Sec-MS-GEC-Version=1-131.0.2903.70", gen_sec_ms_gec())
// }

type Voice struct {
	Name           string   `json:"Name"`
	ShortName      string   `json:"ShortName"`
	Gender         string   `json:"Gender"`
	Locale         string   `json:"Locale"`
	SuggestedCodec string   `json:"SuggestedCodec"`
	FriendlyName   string   `json:"FriendlyName"`
	Status         string   `json:"Status"`
	VoiceTag       VoiceTag `json:"VoiceTag"`
}

type VoiceTag struct {
	ContentCategories  []string `json:"ContentCategories"`
	VoicePersonalities []string `json:"VoicePersonalities"`
}
