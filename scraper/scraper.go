package scraper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Query struct {
	OperationName string `json:"operationName"`
	Variables     struct {
		Slug     string `json:"slug"`
		Platform string `json:"platform"`
	} `json:"variables"`
	Extensions struct {
		PersistedQuery struct {
			Version    int    `json:"version"`
			Sha256Hash string `json:"sha256Hash"`
		} `json:"persistedQuery"`
	} `json:"extensions"`
}

type GQLResponse struct {
	Data struct {
		Clip Clip `json:"clip"`
	} `json:"data"`
	Extensions struct {
		DurationMilliseconds int    `json:"durationMilliseconds"`
		OperationName        string `json:"operationName"`
		RequestID            string `json:"requestID"`
	} `json:"extensions"`
}

type Clip struct {
	ID                  string              `json:"id"`
	PlaybackAccessToken PlaybackAccessToken `json:"playbackAccessToken"`
	VideoQualities      []VideoQuality      `json:"videoQualities"`
	Typename            string              `json:"__typename"`
}

type PlaybackAccessToken struct {
	Signature string `json:"signature"`
	Value     string `json:"value"`
	Typename  string `json:"__typename"`
}

type VideoQuality struct {
	FrameRate float32 `json:"frameRate"`
	Quality   string  `json:"quality"`
	SourceURL string  `json:"sourceURL"`
	Typename  string  `json:"__typename"`
}

type GQLError struct {
	Errors []error `json:"errors"`
}

func (gqle *GQLError) printErrors() {
	for _, error := range gqle.Errors {
		log.Fatal("\n", error)
	}
}

// Hardcoded lmfaoo
const CLIENT_ID = `kimne78kx3ncx6brgo4mv6wki5h1ko`
const SHA256_HASH = `6fd3af2b22989506269b9ac02dd87eb4a6688392d67d94e41a6886f1e9f5c00f`

func check(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func fetch(req *http.Request) []byte {
	client := &http.Client{}
	resp, err := client.Do(req)
	check(err)

	respBody, err := io.ReadAll(resp.Body)
	check(err)

	return respBody
}

func authenticatedPost(url string, body io.Reader) []byte {
	req, err := http.NewRequest("POST", url, body)
	check(err)

	req.Header.Add("Client-ID", CLIENT_ID)

	respBody := fetch(req)

	//fmt.Println(string(respBody))

	return respBody
}

func gqlPersistedQuery(query Query) GQLResponse {
	url := "https://gql.twitch.tv/gql"
	queryByte, err := json.Marshal(query)
	check(err)
	//fmt.Printf("%s\n", string(queryByte))

	respBody := authenticatedPost(url, bytes.NewBuffer(queryByte))

	var gqlResp GQLResponse
	check(json.Unmarshal(respBody, &gqlResp))

	return gqlResp
}

func GetClipAccessToken(slug string) GQLResponse {
	query := Query{
		OperationName: "VideoAccessToken_Clip",
		Variables: struct {
			Slug     string "json:\"slug\""
			Platform string "json:\"platform\""
		}{
			Slug:     slug,
			Platform: "web",
		},
		Extensions: struct {
			PersistedQuery struct {
				Version    int    "json:\"version\""
				Sha256Hash string "json:\"sha256Hash\""
			} "json:\"persistedQuery\""
		}{
			struct {
				Version    int    "json:\"version\""
				Sha256Hash string "json:\"sha256Hash\""
			}{
				Version:    1,
				Sha256Hash: SHA256_HASH,
			},
		},
	}

	resp := gqlPersistedQuery(query)
	return resp
}

func BuildDownloadURL(vq VideoQuality, pbat PlaybackAccessToken) string {
	params := url.Values{}
	params.Add("sig", pbat.Signature)
	params.Add("token", pbat.Value)
	downloadURL := fmt.Sprintf("%s?%s", vq.SourceURL, params.Encode())
	return downloadURL
}
