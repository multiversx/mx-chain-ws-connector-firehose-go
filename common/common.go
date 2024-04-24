package common

const (
	// FullPersisterDBMode defines the db mode where each event will be saved
	// to cache and persister
	FullPersisterDBMode = "full-persister"

	// ImportDBMode defines the db mode where the events will be saved only to cache
	ImportDBMode = "import-db"

	// OptimizedPersisterDBMode defines the db mode where the events will be saved to
	// cache, and they will be dumped to persister when necessary
	OptimizedPersisterDBMode = "optimized-persister"
)
