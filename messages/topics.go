package messages

const (
	MKReport = "report"
	MKNotify = "notify"
	MKSystem = "system"
)

type ReportTopic struct{}

func (ReportTopic) RoutingKey() string { return MKReport }

type NotifyTopic struct{}

func (NotifyTopic) RoutingKey() string { return MKNotify }

type SystemTopic struct{}

func (SystemTopic) RoutingKey() string { return MKSystem }
