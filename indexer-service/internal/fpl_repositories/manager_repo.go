package fpl_repositories

import (
	"database/sql"

	"github.com/imadbelkat1/shared/models"
)

type ManagerRepo struct {
	db             *sql.DB
	Entry          *models.Entry
	EntryPicks     *models.EntryEventPicks
	EntryTransfers *models.EntryTransfers
	EntryHistory   *models.EntryHistory
}

func NewManagerRepo(db *sql.DB, entry *models.Entry, entryPicks *models.EntryEventPicks, entryTransfers *models.EntryTransfers, entryHistory *models.EntryHistory) *ManagerRepo {
	return &ManagerRepo{
		db:             db,
		Entry:          entry,
		EntryPicks:     entryPicks,
		EntryTransfers: entryTransfers,
		EntryHistory:   entryHistory,
	}
}
