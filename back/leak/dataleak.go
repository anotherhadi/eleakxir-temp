package leak

import (
	"database/sql"
	_ "github.com/marcboeker/go-duckdb"
)

type Dataleaks struct {
	Dataleaks          []Dataleak
	DataleaksDirectory string
	CacheDirectory     string
	Duckdb             *sql.DB

	// Stats
	TotalRows      uint64
	TotalDataleaks uint64
	TotalSize      uint64 // MB
}

type Dataleak struct {
	Path    string
	Name    string
	Columns []string
	Length  uint64
	Size    uint64 // In MB
}

func OpenDataleaks(dataleaksDir, cacheDirectory string) (*Dataleaks, error) {
	db, err := sql.Open("duckdb", "")
	if err != nil {
		return nil, err
	}

	d := &Dataleaks{
		DataleaksDirectory: dataleaksDir,
		CacheDirectory:     cacheDirectory,
		Duckdb:             db,
	}

	err = d.getCache()
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (d *Dataleaks) CloseDataleaks() error {
	if d.Duckdb != nil {
		if err := d.Duckdb.Close(); err != nil {
			return err
		}
	}
	return nil
}
