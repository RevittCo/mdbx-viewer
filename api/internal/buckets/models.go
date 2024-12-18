package buckets

const MaxPageLen = 200

type OpenReq struct {
	Path string `json:"path"`
}
