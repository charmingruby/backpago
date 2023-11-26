package users

import (
	"crypto/md5"
	"errors"
	"fmt"
	"time"
)

var (
	ErrPasswordRequired  = errors.New("password is required and can't be blank")
	ErrPasswordMinLength = errors.New("password must be at least 8 characters")
	ErrPasswordMaxLength = errors.New("password must be a maximum of 16 characters")
	ErrNameRequired      = errors.New("name is required and can't be blank")
	ErrNameMinLength     = errors.New("name must be at least 4 characters")
	ErrNameMaxLength     = errors.New("name must be a maximum of 32 characters")
	ErrLoginRequired     = errors.New("login is required and can't be blank")
	ErrLoginMinLength    = errors.New("login must be at least 8 characters")
	ErrLoginMaxLength    = errors.New("login must be a maximum of 24 characters")
)

func New(name, login, password string) (*User, error) {
	u := User{
		Name:       name,
		Login:      login,
		ModifiedAt: time.Now(),
	}

	if err := u.SetPassword(password); err != nil {
		return nil, err
	}

	if err := u.Validate(); err != nil {
		return nil, err
	}

	return &u, nil
}

type User struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	Login      string    `json:"login"`
	Password   string    `json:"password"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
	Deleted    bool      `json:"-"`
	LastLogin  time.Time `json:"last_login"`
}

func (u *User) SetPassword(password string) error {
	if password == "" {
		return ErrPasswordRequired
	}

	if len(password) < 8 {
		return ErrPasswordMinLength
	}

	if len(password) > 16 {
		return ErrPasswordMaxLength
	}

	u.Password = fmt.Sprintf("%x", md5.Sum([]byte(password)))

	return nil
}

func (u *User) Validate() error {
	if u.Name == "" {
		return ErrNameRequired
	}

	if len(u.Name) < 4 {
		return ErrNameMinLength
	}

	if len(u.Name) > 32 {
		return ErrNameMaxLength
	}

	if len(u.Password) < 8 {
		return ErrPasswordMinLength
	}

	if len(u.Password) > 16 {
		return ErrPasswordMaxLength
	}

	if u.Password == fmt.Sprintf("%x", (md5.Sum([]byte("")))) {
		return ErrPasswordRequired
	}

	if u.Login == "" {
		return ErrLoginRequired
	}

	if len(u.Login) < 8 {
		return ErrLoginMinLength
	}

	if len(u.Login) > 24 {
		return ErrLoginMaxLength
	}

	return nil
}
