package messages

const (
	MKReport = "report"
	MKNotify = "notify"
)

type ReportTopic struct{}

func (ReportTopic) RoutingKey() string { return MKReport }

type NotifyTopic struct{}

func (NotifyTopic) RoutingKey() string { return MKNotify }
