package folders

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/charmingruby/backpago/internal/files"
	"github.com/go-chi/chi/v5"
)

func (h *handler) Delete(rw http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = deleteFolderContent(h.db, int64(id))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = Delete(h.db, int64(id))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")

}

func deleteFolderContent(db *sql.DB, folderID int64) error {
	err := deleteFiles(db, folderID)
	if err != nil {
		return err
	}

	return deleteSubfolders(db, folderID)
}

func deleteSubfolders(db *sql.DB, folderID int64) error {
	subfolders, err := getSubFolders(db, folderID)
	if err != nil {
		return err
	}

	removedFolders := make([]Folder, 0, len(subfolders))

	for _, sf := range subfolders {
		err := Delete(db, sf.Id)
		if err != nil {
			break
		}

		err = deleteFolderContent(db, sf.Id)
		if err != nil {
			Update(db, sf.Id, &sf)
			break
		}

		removedFolders = append(removedFolders, sf)
	}

	if len(subfolders) != len(removedFolders) {
		for _, rf := range removedFolders {
			Update(db, rf.Id, &rf)

		}
	}

	return nil
}

func deleteFiles(db *sql.DB, folderID int64) error {
	f, err := files.List(db, int64(folderID))
	if err != nil {
		return err
	}

	removedFiles := make([]files.File, 0, len(f))

	for _, file := range f {
		file.Deleted = true
		err = files.Update(db, file.Id, &file)
		if err != nil {
			break
		}

		removedFiles = append(removedFiles, file)
	}

	if len(f) != len(removedFiles) {
		for _, file := range removedFiles {
			file.Deleted = false
			files.Update(db, int64(folderID), &file)
		}

		return err
	}

	return nil
}

func Delete(db *sql.DB, id int64) error {
	stmt := `update "folders" set "modified_at"=$1, "deleted"=true where id=$2`
	_, err := db.Exec(stmt, time.Now(), id)

	return err

}
