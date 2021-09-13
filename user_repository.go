package main

import (
	"errors"
	"sync"
)

type InMemoryUserStorage struct {
	lock    sync.RWMutex
	storage map[string]User
}

func NewInMemoryUserStorage() *InMemoryUserStorage {
	return &InMemoryUserStorage{
		lock:    sync.RWMutex{},
		storage: make(map[string]User),
	}
}

func (memStore *InMemoryUserStorage) Add(login string, user User) error {
	memStore.lock.RLock()
	_, ok := memStore.storage[login]
	memStore.lock.RUnlock()
	if ok {
		return errors.New("User with such login already exists")
	}
	memStore.lock.Lock()
	memStore.storage[login] = user
	memStore.lock.Unlock()
	return nil
}

func (memStore *InMemoryUserStorage) Update(login string, user User) error {
	memStore.lock.RLock()
	_, ok := memStore.storage[login]
	memStore.lock.RUnlock()
	if !ok {
		return errors.New("User with such login doesn't exist")
	}

	memStore.lock.Lock()
	memStore.storage[login] = user
	memStore.lock.Unlock()
	return nil
}

func (memStore *InMemoryUserStorage) Delete(login string) (User, error) {
	memStore.lock.RLock()
	usr, ok := memStore.storage[login]
	memStore.lock.RUnlock()
	if !ok {
		return User{}, errors.New("There is no such user")
	}

	delete(memStore.storage, login)
	return usr, nil
}

func (memStore *InMemoryUserStorage) Get(login string) (User, error) {
	memStore.lock.RLock()
	usr, ok := memStore.storage[login]
	memStore.lock.RUnlock()
	if !ok {
		return User{}, errors.New("There is no such user")
	}

	return usr, nil
}
