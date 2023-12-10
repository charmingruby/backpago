package folders

import (
	"errors"
	"time"
)

var (
	ErrNameRequired  = errors.New("name is required and can't be blank")
	ErrNameMinLength = errors.New("name must be at least 2 characters")
	ErrNameMaxLength = errors.New("name must be a maximum of 32 characters")
)

func New(parentId int64, name string) (*Folder, error) {
	f := Folder{
		ParentId: parentId,
		Name:     name,
	}

	if err := f.Validate(); err != nil {
		return nil, err
	}

	return &f, nil
}

type Folder struct {
	Id         int64     `json:"id"`
	ParentId   int64     `json:"parent_id"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
	Deleted    bool      `json:"-"`
}

func (f *Folder) Validate() error {
	if f.Name == "" {
		return ErrNameRequired
	}

	if len(f.Name) < 2 {
		return ErrNameMinLength
	}

	if len(f.Name) > 32 {
		return ErrNameMaxLength
	}

	return nil
}

type FolderContent struct {
	Folder  Folder           `json:"folder"`
	Content []FolderResource `json:"content"`
}

type FolderResource struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}
