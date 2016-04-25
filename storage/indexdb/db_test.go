// Copyright 2015 Eleme Inc. All rights reserved.

package indexdb

import (
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/util"
	"github.com/syndtr/goleveldb/leveldb"
	"os"
	"testing"
)

func TestOpen(t *testing.T) {
	fileName := "db-testing"
	db, err := Open(fileName)
	util.Must(t, err == nil)
	util.Must(t, util.IsFileExist(fileName))
	db.Close()
	os.RemoveAll(fileName)
}

func TestLoad(t *testing.T) {
	fileName := "db-testing"
	db, _ := Open(fileName)
	defer os.RemoveAll(fileName)
	defer db.Close()
	idx := &models.Index{Name: "foo", Stamp: 1450430839, Score: 0.7, Average: 78.5}
	// Add one
	db.Put(idx)
	util.Must(t, db.tr.Has(idx.Name))
	// Clear cache
	db.tr.Clear()
	util.Must(t, db.tr.Len() == 0)
	util.Must(t, !db.tr.Has(idx.Name))
	// Reload
	db.load()
	// Must not empty and idx in cache
	util.Must(t, db.tr.Len() == 1)
	util.Must(t, db.tr.Has(idx.Name))
}

func TestPut(t *testing.T) {
	// Open db.
	fileName := "db-testing"
	db, _ := Open(fileName)
	defer os.RemoveAll(fileName)
	defer db.Close()
	// Test.
	idx := &models.Index{Name: "foo", Stamp: 1450430837, Score: 1.2, Average: 109.5}
	err := db.Put(idx)
	util.Must(t, err == nil)
	// Must in cache
	util.Must(t, db.tr.Has(idx.Name))
	// Must in db file
	v, err := db.db.Get([]byte(idx.Name), nil)
	util.Must(t, err == nil)
	idx1 := &models.Index{}
	decode(v, idx1)
	util.Must(t, idx1.Stamp == idx.Stamp)
	util.Must(t, idx1.Score == idx.Score)
	util.Must(t, idx1.Average == idx.Average)
}

func TestGet(t *testing.T) {
	// Open db.
	fileName := "db-testing"
	db, _ := Open(fileName)
	defer os.RemoveAll(fileName)
	defer db.Close()
	// Not found.
	_, err := db.Get("Not-exist")
	util.Must(t, ErrNotFound == err)
	// Put one.
	idx := &models.Index{Name: "foo", Stamp: 1450430837, Score: 0.3, Average: 100}
	db.Put(idx)
	// Get it from cache.
	i, err := db.Get(idx.Name)
	util.Must(t, nil == err)
	util.Must(t, i.Equal(idx))
}

func TestDelete(t *testing.T) {
	// Open db.
	fileName := "db-testing"
	db, _ := Open(fileName)
	defer os.RemoveAll(fileName)
	defer db.Close()
	// Add one.
	idx := &models.Index{Name: "foo", Stamp: 1450430837, Score: 0.3, Average: 100}
	db.Put(idx)
	// Must in cache.
	util.Must(t, db.tr.Has(idx.Name))
	// Delete it.
	err := db.Delete(idx.Name)
	util.Must(t, err == nil)
	// Must not exist in cache
	util.Must(t, !db.tr.Has(idx.Name))
	// Must not in db.
	_, err = db.db.Get([]byte(idx.Name), nil)
	util.Must(t, err == leveldb.ErrNotFound)
	// Cant get again.
	_, err = db.Get(idx.Name)
	util.Must(t, ErrNotFound == err)
}

func TestFilter(t *testing.T) {
	// Open db.
	fileName := "db-testing"
	db, _ := Open(fileName)
	defer os.RemoveAll(fileName)
	defer db.Close()
	// Add indexes.
	excludeName := "a.b.c.d.e.f"
	db.Put(&models.Index{Name: "a.b.c.e.f.g"})
	db.Put(&models.Index{Name: "a.b.c.d.f.g"})
	db.Put(&models.Index{Name: excludeName})
	// Filter
	l := db.Filter("a.b.c.*.f.*")
	util.Must(t, len(l) == 2)
	util.Must(t, l[0].Name != excludeName && l[1].Name != excludeName)
}

func TestLen(t *testing.T) {
	// Open db.
	fileName := "db-testing"
	db, _ := Open(fileName)
	defer os.RemoveAll(fileName)
	defer db.Close()
	// Add indexes.
	db.Put(&models.Index{Name: "abc"})
	db.Put(&models.Index{Name: "efg"})
	db.Put(&models.Index{Name: "fgh"})
	// Len
	util.Must(t, db.Len() == 3)
}
