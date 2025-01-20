package buckets

const MaxPageLen = 200

type OpenReq struct {
	Path string `json:"path"`
}

type TreeElement struct {
	Id           string        `json:"id"`
	Name         string        `json:"name"`
	FullPath     string        `json:"fullPath"`
	IsSelectable bool          `json:"isSelectable"`
	Children     []TreeElement `json:"children"`
}

type DataSource struct {
	Source       string        `json:"source"`
	TreeElements []TreeElement `json:"treeElements"`
}
