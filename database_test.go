package main

import "testing"

func TestReadFirebase(t *testing.T) {
	db, _ := ReadFirebase()
	t.Log(db.Version.Newlink)
}
