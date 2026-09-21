package constants

const (
	ServiceSource = "source"
	ServiceNews   = "news"
	ServiceCrawl  = "crawl"
	ServiceReport = "report"
	ServiceNotify = "notify"
	ServiceMCP    = "mcp"
)

func AllServices() []string {
	return []string{
		ServiceSource,
		ServiceNews,
		ServiceCrawl,
		ServiceReport,
		ServiceNotify,
		ServiceMCP,
	}
}
