package events

import "nfxnews/pkgs/kafkax/eventbus"

const (
	TKSource    eventbus.TopicKey = "source"
	TKSourceDLQ eventbus.TopicKey = "source_poison"

	TKNews    eventbus.TopicKey = "news"
	TKNewsDLQ eventbus.TopicKey = "news_poison"

	TKCrawl    eventbus.TopicKey = "crawl"
	TKCrawlDLQ eventbus.TopicKey = "crawl_poison"

	TKReport    eventbus.TopicKey = "report"
	TKReportDLQ eventbus.TopicKey = "report_poison"

	TKNotify    eventbus.TopicKey = "notify"
	TKNotifyDLQ eventbus.TopicKey = "notify_poison"

	TKMCP    eventbus.TopicKey = "mcp"
	TKMCPDLQ eventbus.TopicKey = "mcp_poison"
)

type SourceTopic struct{}

func (SourceTopic) TopicKey() eventbus.TopicKey { return TKSource }

type NewsTopic struct{}

func (NewsTopic) TopicKey() eventbus.TopicKey { return TKNews }

type CrawlTopic struct{}

func (CrawlTopic) TopicKey() eventbus.TopicKey { return TKCrawl }

type ReportTopic struct{}

func (ReportTopic) TopicKey() eventbus.TopicKey { return TKReport }

type NotifyTopic struct{}

func (NotifyTopic) TopicKey() eventbus.TopicKey { return TKNotify }

type MCPTopic struct{}

func (MCPTopic) TopicKey() eventbus.TopicKey { return TKMCP }
