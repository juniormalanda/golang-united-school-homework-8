package file

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	rootDir = "data/"
)

type NotFoundError error

type File struct {
	name string
	file *os.File
}

func NewFile(fileName string) *File {
	file := &File{name: fileName}
	file.open()

	return file
}

func (f *File) Close() {
	f.file.Close()
}

func (f *File) open() error {
	if f.file == nil {
		file, err := os.OpenFile(rootDir+f.name, os.O_RDWR|os.O_CREATE, 0755)

		if err != nil {
			return nil
		}

		f.file = file
	}

	return nil
}

func (f *File) AddUser(item string) (err error) {
	var user User

	err = json.Unmarshal([]byte(item), user)

	if err != nil {
		return
	}

	users, err := f.Users()

	if err != nil {
		return
	}

	users = append(users, user)

	err = f.putUsers(users)

	return
}

func (f *File) Users() (users []User, err error) {
	data, err := f.List()

	if err != nil {
		return
	}

	err = json.Unmarshal(data, users)

	return
}

func (f *File) putUsers(users []User) (err error) {
	data, err := json.Marshal(users)

	if err != nil {
		return
	}

	_, err = f.file.Write(data)

	return
}

func (f *File) List() ([]byte, error) {
	var data []byte

	f.file.Read(data)

	return data, nil
}

func (f *File) Remove(id string) (succeed bool, err error) {
	users, err := f.Users()

	if err != nil {
		return
	}

	for i, user := range users {
		if user.Id == id {
			users = append(users[:i], users[i+1:]...)
			f.putUsers(users)
			succeed = true
			return
		}
	}

	err = fmt.Errorf("Item with id %s not found", id)

	return
}

func (f *File) FindById(id string) (data []byte, err error) {
	users, err := f.Users()

	if err != nil {
		return
	}

	for _, user := range users {
		if user.Id == id {
			data, err = json.Marshal(user)
			return
		}
	}

	return
}
