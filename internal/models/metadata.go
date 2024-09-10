package models

import (
	"database/sql"
	"fmt"

	"github.com/JuHaNi654/cms/internal/database"
)

type Metadata struct {
	Id    int
	Ready bool
}

func IsServceInitialized(db *database.SqlClient) (bool, error) {
  var ready bool 
  q := `
  SELECT
   ready
  FROM metadata 
  WHERE id = 1
  `

  err := db.Db.QueryRow(q).Scan(&ready)
  if err != nil {
    if err == sql.ErrNoRows {
      return false, fmt.Errorf("IsServceInitialized: metadata not found")
    }

    return false, fmt.Errorf("IsServceInitialized: %v", err)
  }

  return ready, nil
}
