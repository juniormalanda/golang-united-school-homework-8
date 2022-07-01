package file

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

var NotFoundError = errors.New("Not found")
var ItemExistsError = errors.New("Item already exists")

type File struct {
	name string
}

func NewFile(fileName string) *File {
	return &File{name: fileName}
}

func (f *File) AddUser(item string) (user User, err error) {
	err = json.Unmarshal([]byte(item), &user)

	if err != nil {
		return
	}

	users, err := f.Users()

	if err != nil {
		return
	}

	for _, u := range users {
		if u.Id == user.Id {
			return user, ItemExistsError
		}
	}

	users = append(users, user)

	err = f.put(users)

	return
}

func (f *File) Users() (users []User, err error) {
	data, err := f.List()

	if err != nil {
		return
	}

	if len(data) == 0 {
		return
	}

	err = json.Unmarshal(data, &users)

	return
}

func (f *File) put(users []User) (err error) {
	data, err := json.Marshal(users)

	if err != nil {
		return
	}

	file, err := os.OpenFile(f.name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)

	if err != nil {
		return nil
	}

	defer file.Close()
	_, err = file.Write(data)

	return
}

func (f *File) List() ([]byte, error) {
	file, err := os.OpenFile(f.name, os.O_RDWR|os.O_CREATE, 0755)

	if err != nil {
		return nil, err
	}

	defer file.Close()
	return ioutil.ReadAll(file)
}

func (f *File) Remove(id string) (err error) {
	users, err := f.Users()

	if err != nil {
		return
	}

	for i, user := range users {
		if user.Id == id {
			users = append(users[:i], users[i+1:]...)
			f.put(users)
			return
		}
	}

	return NotFoundError
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
