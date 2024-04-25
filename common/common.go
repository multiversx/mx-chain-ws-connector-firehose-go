package common

// DBMode defines db mode type
type DBMode string

const (
	// FullPersisterDBMode defines the db mode where each event will be saved
	// to cache and persister
	FullPersisterDBMode DBMode = "full-persister"

	// ImportDBMode defines the db mode where the events will be saved only to cache
	ImportDBMode DBMode = "import-db"

	// OptimizedPersisterDBMode defines the db mode where the events will be saved to
	// cache, and they will be dumped to persister when necessary
	OptimizedPersisterDBMode DBMode = "optimized-persister"
)
