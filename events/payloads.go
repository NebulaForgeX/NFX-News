package events

type SourceFetchedEvent struct {
	SourceTopic
	SourceID string           `json:"source_id"`
	Items    []SourceNewsItem `json:"items"`
}

type SourceNewsItem struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	MobileURL string `json:"mobile_url"`
	PubDate   int64  `json:"pub_date"`
	ExtraJSON string `json:"extra_json"`
}

type ReportGeneratedEvent struct {
	ReportTopic
	ReportID  string `json:"report_id"`
	Mode      string `json:"mode"`
	Title     string `json:"title"`
	ItemCount int    `json:"item_count"`
	Payload   string `json:"payload_json"`
}
