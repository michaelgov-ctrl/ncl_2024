package main

import (
	"bufio"
	"errors"
	"io"
	"sync"
)

type appMap struct {
	sync.Mutex
	m     map[string]string
	len   int
	found int
}

func newAppMap() *appMap {
	var am appMap
	am.m = map[string]string{}

	return &am
}

func loadNewAppMap(r io.Reader) *appMap {
	am := newAppMap()

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		am.store(scanner.Text(), "")
	}

	return am
}

func (am *appMap) incFound() {
	am.Lock()
	am.found++
	am.Unlock()
}

func (am *appMap) allFound() bool {
	am.Lock()
	l, f := am.len, am.found
	am.Unlock()

	return l == f
}

func (am *appMap) store(key, val string) {
	am.Lock()
	am.len++
	am.m[key] = val
	am.Unlock()
}

func (am *appMap) exists(key string) bool {
	am.Lock()
	_, ok := am.m[key]
	am.Unlock()

	return ok
}

func (am *appMap) updateOnce(key, val string) error {
	am.Lock()
	v, ok := am.m[key]
	am.Unlock()

	if ok && v != "" {
		return errors.New("non empty entry")
	}

	am.Lock()
	am.m[key] = val
	am.Unlock()

	return nil
}

func (am *appMap) updateOnceIfExists(key, val string) (bool, error) {
	if !am.exists(key) {
		return false, nil
	}

	if err := am.updateOnce(key, val); err != nil {
		return true, err
	}

	return true, nil
}
