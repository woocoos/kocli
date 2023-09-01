package schema

type UtilItem struct {
	Name    string  `json:"name"`
	Type    string  `json:"type"`
	Content NpmInfo `json:"content"`
}
