package scraper

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
