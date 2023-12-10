package folders

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *handler) Get(rw http.ResponseWriter, r *http.Request) {
	folderID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	f, err := GetFolder(h.db, int64(folderID))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	c, err := GetFolderContent(h.db, int64(folderID))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	fc := FolderContent{
		Folder:  *f,
		Content: c,
	}

	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(fc)

}

func GetFolder(db *sql.DB, folderID int64) (*Folder, error) {
	stmt := `select * from "folders" where id=$1`

	row := db.QueryRow(stmt, folderID)

	var f Folder
	err := row.Scan(
		&f.Id,
		&f.ParentId,
		&f.Name,
		&f.CreatedAt,
		&f.ModifiedAt,
		&f.Deleted,
	)
	if err != nil {
		return nil, err
	}

	return &f, nil
}

func getSubFolder(db *sql.DB, folderID int64) ([]Folder, error) {
	stmt := `select * from "folders" where "parent_id"=$1 and "deleted"=false`

	rows, err := db.Query(stmt, folderID)
	if err != nil {
		return nil, err
	}

	f := make([]Folder, 0)
	for rows.Next() {
		var folder Folder
		err = rows.Scan(
			&folder.Id,
			&folder.ParentId,
			&folder.Name,
			&folder.CreatedAt,
			&folder.ModifiedAt,
			&folder.Deleted,
		)
		if err != nil {
			continue
		}

		f = append(f, folder)

	}

	return f, nil
}

func GetFolderContent(db *sql.DB, folderID int64) ([]FolderResource, error) {
	subfolders, err := getSubFolder(db, folderID)
	if err != nil {
		return nil, err
	}

	fr := make([]FolderResource, 0, len(subfolders))
	for _, sf := range subfolders {
		r := FolderResource{
			ID:         sf.Id,
			Name:       sf.Name,
			Type:       "directory",
			CreatedAt:  sf.CreatedAt,
			ModifiedAt: sf.ModifiedAt,
		}

		fr = append(fr, r)
	}

	return fr, nil
}
