package main

import (
	"crypto/rand"
	"encoding/hex"

	consul "github.com/hashicorp/consul/api"
)

type User struct {
	Name   string `json:"name"`
	Secret []byte `json:"secret"`
}

var KV = initKV()

func initKV() *consul.KV {
	client, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		panic(err)
	}
	kv := client.KV()
	return kv
}

func (u User) Create() (User, error) {
	u.Secret = generateSecret()
	p := &consul.KVPair{Key: u.Name, Value: u.Secret}
	_, err := KV.Put(p, nil)
	return u, err
}

func (u User) Exist() bool {
	res, _, _ := KV.Get(u.Name, nil)
	if res != nil {
		return true
	}
	return false
}

func (u User) Delete() error {
	_, err := KV.Delete(u.Name, nil)
	return err
}

func (u User) GetAll() ([]User, error) {
	var users []User

	list, _, err := KV.List("", nil)

	if err != nil {
		return users, err
	}

	for _, u := range list {
		users = append(users, User{u.Key, u.Value})
	}

	return users, err

}

func InitSecrets() error {
	users, err := User{}.GetAll()
	if err != nil {
		return err
	}
	for _, u := range users {
		if err != nil {
			return err
		}
		Secrets[u.Name] = u.Secret
	}
	return err
}

func generateSecret() string {
	const len = 16
	bytes := make([]byte, len)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
