package main

import (
	"errors"
	"testing"
)

func TestPut(t *testing.T) {
	var key = "create_key"
	var val = "create_val"

	var value interface{}
	var contains bool

	defer delete(store, key)

	_, contains = store[key]
	if contains {
		t.Error("key/value already exist")
	}

	err := Put(key, val)
	if err != nil {
		t.Error(err)
	}

	value, contains = store[key]
	if !contains {
		t.Error("create failed")
	}

	if value != val {
		t.Error("value mismatch")
	}
}

func TestGet(t *testing.T) {
	const key = "read_key"
	const val = "read_val"

	var value interface{}
	var err error

	defer delete(store, key)

	_, err = Get(key)
	if err == nil {
		t.Error("expected an error")
	}
	if !errors.Is(err, ErrorNoSuchKey) {
		t.Error("unexpected error: ", err)
	}

	store[key] = val

	value, err = Get(key)
	if err != nil {
		t.Error("unexpected error: ", err)
	}

	if value != val {
		t.Error("val/value mismatch")
	}
}

func TestDelete(t *testing.T) {
	var key = "delete_key"
	var val = "delete_val"

	var contains bool

	defer delete(store, key)

	store[key] = val

	_, contains = store[key]
	if !contains {
		t.Error("key/value doesn't exist")
	}

	Delete(key)

	_, contains = store[key]
	if contains {
		t.Error("delete failed")
	}
}
