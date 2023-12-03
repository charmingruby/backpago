package files

import (
	"database/sql"
	"time"
)

func Update(db *sql.DB, id int64, f *File) error {
	f.ModifiedAt = time.Now()

	stmt := `update "files" set "name"=$1, "deleted"=$2, "modified_at"=$3 where id=$4`
	_, err := db.Exec(stmt, f.Name, f.Deleted, f.ModifiedAt, id)
	if err != nil {
		return err
	}

	return nil
}
