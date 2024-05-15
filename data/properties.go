package data

const (
	SyncResultsProperty = "sync-results"
	LastSetIPsProperty  = "last-set-ips"
)

func AllProperties() []string {
	return []string{
		SyncResultsProperty,
		LastSetIPsProperty,
	}
}
