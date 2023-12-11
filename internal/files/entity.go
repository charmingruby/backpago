package files

import (
	"errors"
	"time"
)

var (
	ErrNameRequired  = errors.New("name is required and can't be blank")
	ErrNameMinLength = errors.New("name must be at least 2 characters")
	ErrNameMaxLength = errors.New("name must be a maximum of 32 characters")
	ErrOwnerRequired = errors.New("owner is required and can't be blank")
	ErrTypeRequired  = errors.New("type is required and can't be blank")
	ErrPathRequired  = errors.New("path is required and can't be blank")
)

func New(ownerID int64, name, fileType, path string) (*File, error) {
	f := File{
		OwnerId:    ownerID,
		Name:       name,
		Type:       fileType,
		Path:       path,
		ModifiedAt: time.Now(),
	}

	err := f.Validate()
	if err != nil {
		return nil, err
	}

	return &f, nil

}

type File struct {
	Id         int64     `json:"id"`
	FolderId   int64     `json:"-"`
	OwnerId    int64     `json:"owner_id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Path       string    `json:"-"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
	Deleted    bool      `json:"-"`
}

func (f *File) Validate() error {
	if f.Name == "" {
		return ErrNameRequired
	}

	if len(f.Name) < 2 {
		return ErrNameMinLength
	}

	if len(f.Name) > 32 {
		return ErrNameMaxLength
	}

	if f.OwnerId == 0 {
		return ErrOwnerRequired
	}

	if f.Type == "" {
		return ErrTypeRequired
	}

	if f.Path == "" {
		return ErrPathRequired
	}

	return nil
}
