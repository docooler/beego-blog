package loveletter

import (
	"github.com/syndtr/goleveldb/leveldb"
	"strconv"
)

const (
	DB_LOVE_LETTER_FILE = "../db/loveletter"
)

var (
	db, _   = leveldb.OpenFile(DB_LOVE_LETTER_FILE, nil)
	NrLeter = 0
)

func init() {
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		NrLeter += 1
	}
}

func Add(letter string) {
	NrLeter += 1

	db.Put([]byte(strconv.Itoa(NrLeter)), []byte(letter), nil)
}
func GetNrLetter() int {
	return NrLeter
}

func Get(index int) string {
	data, _ := db.Get([]byte(strconv.Itoa(index)), nil)
	return string(data)
}

func Delete(index int) {
	db.Delete([]byte(strconv.Itoa(index)), nil)
}
