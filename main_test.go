package rixxdb

import (
	"fmt"
	"os"
	"sync"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var test = []byte("test")

func TestPersistence(t *testing.T) {

	var db *DB
	var err error

	Convey("No persistence", t, func() {
		db, err = Open("memory", &Config{EncryptionKey: nil})
		So(err, ShouldBeNil)
		So(db, ShouldNotBeNil)
		So(db.Flush(), ShouldEqual, ErrDbMemoryOnly)
		So(db.Close(), ShouldBeNil)
	})

	Convey("No persistence", t, func() {
		db, err = Open("memory", &Config{EncryptionKey: []byte("1234567890123456")})
		So(err, ShouldBeNil)
		So(db, ShouldNotBeNil)
		So(db.Flush(), ShouldEqual, ErrDbMemoryOnly)
		So(db.Close(), ShouldBeNil)
	})

	Convey("Invalid key", t, func() {
		db, err = Open("memory", &Config{EncryptionKey: []byte("12345678901234567890")})
		So(err, ShouldEqual, ErrDbInvalidEncryptionKey)
	})

	Convey("Path persistence", t, func() {
		db, err = Open("test.db", &Config{})
		So(err, ShouldBeNil)
		So(db, ShouldNotBeNil)
		So(db.Flush(), ShouldBeNil)
		So(db.Close(), ShouldBeNil)
		db, err = Open("test.db", &Config{})
		So(err, ShouldBeNil)
		So(db, ShouldNotBeNil)
		So(db.Close(), ShouldBeNil)
		os.RemoveAll("test.db")
	})

	Convey("File persistence", t, func() {
		db, err = Open("file://test.db", &Config{})
		So(err, ShouldBeNil)
		So(db, ShouldNotBeNil)
		So(db.Flush(), ShouldBeNil)
		So(db.Close(), ShouldBeNil)
		db, err = Open("file://test.db", &Config{})
		So(err, ShouldBeNil)
		So(db, ShouldNotBeNil)
		So(db.Close(), ShouldBeNil)
		os.RemoveAll("test.db")
	})

}

func TestTransactions(t *testing.T) {

	Convey("Check transaction errors", t, func() {

		var db *DB
		var tx *TX
		var err error

		db, err = Open("memory", &Config{})

		tx, err = db.Begin(false)

		// Test writing to a read transaction

		_, err = tx.Clr(test)
		So(err, ShouldEqual, ErrTxNotWritable)
		_, err = tx.ClrL(test, 0)
		So(err, ShouldEqual, ErrTxNotWritable)
		_, err = tx.ClrP(test, 0)
		So(err, ShouldEqual, ErrTxNotWritable)
		_, err = tx.ClrR(test, test, 0)
		So(err, ShouldEqual, ErrTxNotWritable)

		_, err = tx.Del(0, test)
		So(err, ShouldEqual, ErrTxNotWritable)
		_, err = tx.DelC(0, test, nil)
		So(err, ShouldEqual, ErrTxNotWritable)
		_, err = tx.DelL(0, test, 0)
		So(err, ShouldEqual, ErrTxNotWritable)
		_, err = tx.DelP(0, test, 0)
		So(err, ShouldEqual, ErrTxNotWritable)
		_, err = tx.DelR(0, test, test, 0)
		So(err, ShouldEqual, ErrTxNotWritable)

		_, err = tx.Put(0, test, nil)
		So(err, ShouldEqual, ErrTxNotWritable)
		_, err = tx.PutC(0, test, nil, nil)
		So(err, ShouldEqual, ErrTxNotWritable)
		_, err = tx.PutL(0, test, nil, 0)
		So(err, ShouldEqual, ErrTxNotWritable)
		_, err = tx.PutP(0, test, nil, 0)
		So(err, ShouldEqual, ErrTxNotWritable)
		_, err = tx.PutR(0, test, test, nil, 0)
		So(err, ShouldEqual, ErrTxNotWritable)

		err = tx.Commit()
		So(err, ShouldEqual, ErrTxNotWritable)

		// Try altering a closed transaction

		err = tx.Cancel()
		So(err, ShouldEqual, ErrTxClosed)

		err = tx.Commit()
		So(err, ShouldEqual, ErrTxClosed)

		_, err = tx.Clr(test)
		So(err, ShouldEqual, ErrTxClosed)
		_, err = tx.ClrL(test, 0)
		So(err, ShouldEqual, ErrTxClosed)
		_, err = tx.ClrP(test, 0)
		So(err, ShouldEqual, ErrTxClosed)
		_, err = tx.ClrR(test, test, 0)
		So(err, ShouldEqual, ErrTxClosed)

		_, err = tx.Get(0, test)
		So(err, ShouldEqual, ErrTxClosed)
		_, err = tx.GetL(0, test, 0)
		So(err, ShouldEqual, ErrTxClosed)
		_, err = tx.GetP(0, test, 0)
		So(err, ShouldEqual, ErrTxClosed)
		_, err = tx.GetR(0, test, test, 0)
		So(err, ShouldEqual, ErrTxClosed)

		_, err = tx.Del(0, test)
		So(err, ShouldEqual, ErrTxClosed)
		_, err = tx.DelC(0, test, nil)
		So(err, ShouldEqual, ErrTxClosed)
		_, err = tx.DelL(0, test, 0)
		So(err, ShouldEqual, ErrTxClosed)
		_, err = tx.DelP(0, test, 0)
		So(err, ShouldEqual, ErrTxClosed)
		_, err = tx.DelR(0, test, test, 0)
		So(err, ShouldEqual, ErrTxClosed)

		_, err = tx.Put(0, test, test)
		So(err, ShouldEqual, ErrTxClosed)
		_, err = tx.PutC(0, test, test, nil)
		So(err, ShouldEqual, ErrTxClosed)
		_, err = tx.PutL(0, test, test, 0)
		So(err, ShouldEqual, ErrTxClosed)
		_, err = tx.PutP(0, test, test, 0)
		So(err, ShouldEqual, ErrTxClosed)
		_, err = tx.PutR(0, test, test, nil, 0)
		So(err, ShouldEqual, ErrTxClosed)

		tx.Cancel()

		// Try passing invalid arguments

		tx, err = db.Begin(true)

		_, err = tx.Clr(nil)
		So(err, ShouldEqual, ErrTxKeyCanNotBeNil)
		_, err = tx.ClrL(nil, 0)
		So(err, ShouldEqual, ErrTxKeyCanNotBeNil)
		_, err = tx.ClrP(nil, 0)
		So(err, ShouldEqual, ErrTxKeyCanNotBeNil)
		_, err = tx.ClrR(nil, test, 0)
		So(err, ShouldEqual, ErrTxKeyCanNotBeNil)
		_, err = tx.ClrR(test, nil, 0)
		So(err, ShouldEqual, ErrTxKeyCanNotBeNil)

		_, err = tx.Get(0, nil)
		So(err, ShouldEqual, ErrTxKeyCanNotBeNil)
		_, err = tx.GetL(0, nil, 0)
		So(err, ShouldEqual, ErrTxKeyCanNotBeNil)
		_, err = tx.GetP(0, nil, 0)
		So(err, ShouldEqual, ErrTxKeyCanNotBeNil)
		_, err = tx.GetR(0, nil, test, 0)
		So(err, ShouldEqual, ErrTxKeyCanNotBeNil)
		_, err = tx.GetR(0, test, nil, 0)
		So(err, ShouldEqual, ErrTxKeyCanNotBeNil)

		_, err = tx.Del(0, nil)
		So(err, ShouldEqual, ErrTxKeyCanNotBeNil)
		_, err = tx.DelC(0, nil, nil)
		So(err, ShouldEqual, ErrTxKeyCanNotBeNil)
		_, err = tx.DelL(0, nil, 0)
		So(err, ShouldEqual, ErrTxKeyCanNotBeNil)
		_, err = tx.DelP(0, nil, 0)
		So(err, ShouldEqual, ErrTxKeyCanNotBeNil)
		_, err = tx.DelR(0, nil, test, 0)
		So(err, ShouldEqual, ErrTxKeyCanNotBeNil)
		_, err = tx.DelR(0, test, nil, 0)
		So(err, ShouldEqual, ErrTxKeyCanNotBeNil)

		_, err = tx.Put(0, nil, test)
		So(err, ShouldEqual, ErrTxKeyCanNotBeNil)
		_, err = tx.PutC(0, nil, test, nil)
		So(err, ShouldEqual, ErrTxKeyCanNotBeNil)
		_, err = tx.PutL(0, nil, test, 0)
		So(err, ShouldEqual, ErrTxKeyCanNotBeNil)
		_, err = tx.PutP(0, nil, test, 0)
		So(err, ShouldEqual, ErrTxKeyCanNotBeNil)
		_, err = tx.PutR(0, nil, test, nil, 0)
		So(err, ShouldEqual, ErrTxKeyCanNotBeNil)
		_, err = tx.PutR(0, test, nil, nil, 0)
		So(err, ShouldEqual, ErrTxKeyCanNotBeNil)

		tx.Cancel()

		// Test managed transaction errors

		db.View(func(tx *TX) error {
			return fmt.Errorf("Test error")
		})

		db.Update(func(tx *TX) error {
			return fmt.Errorf("Test error")
		})

		db.View(func(tx *TX) error {
			err = tx.Cancel()
			So(err, ShouldEqual, ErrTxNotEditable)
			err = tx.Commit()
			So(err, ShouldEqual, ErrTxNotEditable)
			return nil
		})

		db.Update(func(tx *TX) error {
			err = tx.Cancel()
			So(err, ShouldEqual, ErrTxNotEditable)
			err = tx.Commit()
			So(err, ShouldEqual, ErrTxNotEditable)
			return nil
		})

		// Test managed transaction panics

		func() {

			defer func() { recover() }()

			db.View(func(tx *TX) error {
				panic("test")
			})

		}()

		func() {

			defer func() { recover() }()

			db.Update(func(tx *TX) error {
				panic("test")
			})

		}()

	})

	Convey("Check cancelling transactions", t, func() {

		var db *DB
		var tx *TX
		var kv *KV
		var ok bool
		var err error

		db, err = Open("memory", &Config{})
		So(err, ShouldBeNil)

		tx, err = db.Begin(true)
		So(err, ShouldBeNil)
		So(tx, ShouldNotBeNil)

		ok = tx.Closed()
		So(ok, ShouldBeFalse)

		kv, err = tx.Put(0, []byte("test"), []byte("tester"))
		So(err, ShouldBeNil)
		So(kv, ShouldNotBeNil)

		kv, err = tx.Get(0, []byte("test"))
		So(err, ShouldBeNil)
		So(kv, ShouldNotBeNil)
		So(kv.Exi(), ShouldBeTrue)

		err = tx.Cancel()
		So(err, ShouldBeNil)

		ok = tx.Closed()
		So(ok, ShouldBeTrue)

		tx, err = db.Begin(false)
		So(err, ShouldBeNil)
		So(tx, ShouldNotBeNil)

		kv, err = tx.Get(0, []byte("test"))
		So(err, ShouldBeNil)
		So(kv, ShouldNotBeNil)
		So(kv.Exi(), ShouldBeFalse)

		err = tx.Cancel()
		So(err, ShouldBeNil)

	})

	Convey("Check committing transactions", t, func() {

		var db *DB
		var tx *TX
		var kv *KV
		var ok bool
		var err error

		db, err = Open("memory", &Config{})
		So(err, ShouldBeNil)

		tx, err = db.Begin(true)
		So(err, ShouldBeNil)
		So(tx, ShouldNotBeNil)

		ok = tx.Closed()
		So(ok, ShouldBeFalse)

		kv, err = tx.Put(0, []byte("test"), []byte("tester"))
		So(err, ShouldBeNil)
		So(kv, ShouldNotBeNil)

		kv, err = tx.Get(0, []byte("test"))
		So(err, ShouldBeNil)
		So(kv, ShouldNotBeNil)
		So(kv.Exi(), ShouldBeTrue)

		err = tx.Commit()
		So(err, ShouldBeNil)

		ok = tx.Closed()
		So(ok, ShouldBeTrue)

		tx, err = db.Begin(false)
		So(err, ShouldBeNil)
		So(tx, ShouldNotBeNil)

		kv, err = tx.Get(0, []byte("test"))
		So(err, ShouldBeNil)
		So(kv, ShouldNotBeNil)
		So(kv.Exi(), ShouldBeTrue)
		So(kv.Ver(), ShouldEqual, 0)
		So(kv.Key(), ShouldResemble, []byte("test"))
		So(kv.Val(), ShouldResemble, []byte("tester"))

		err = tx.Cancel()
		So(err, ShouldBeNil)

	})

	Convey("Write many keys concurrently", t, func() {

		var db *DB
		var tx *TX
		var kv *KV
		var ok bool
		var err error
		var max = 100

		db, err = Open("memory", &Config{})
		So(err, ShouldBeNil)

		tx, err = db.Begin(true)
		So(err, ShouldBeNil)
		So(tx, ShouldNotBeNil)

		ok = tx.Closed()
		So(ok, ShouldBeFalse)

		var w sync.WaitGroup

		w.Add(3)

		go func() {
			for i := 1; i < max; i++ {
				tx.Put(0, []byte(fmt.Sprint(i)), []byte(fmt.Sprint(i)))
			}
			w.Done()
		}()

		go func() {
			for i := 1; i < max; i++ {
				tx.Put(0, []byte(fmt.Sprint(i)), []byte(fmt.Sprint(i)))
			}
			w.Done()
		}()

		go func() {
			for i := 1; i < max; i++ {
				tx.Put(0, []byte(fmt.Sprint(i)), []byte(fmt.Sprint(i)))
			}
			w.Done()
		}()

		w.Wait()

		err = tx.Commit()
		So(err, ShouldBeNil)

		tx, err = db.Begin(false)
		So(err, ShouldBeNil)
		So(tx, ShouldNotBeNil)

		for i := 1; i < max; i++ {
			kv, err = tx.Get(0, []byte(fmt.Sprint(i)))
			So(err, ShouldBeNil)
			So(kv, ShouldNotBeNil)
			So(kv.Exi(), ShouldBeTrue)
			So(kv.Val(), ShouldResemble, []byte(fmt.Sprint(i)))
		}

		err = tx.Cancel()
		So(err, ShouldBeNil)

	})

}

func TestOperations(t *testing.T) {

	var db *DB
	var tx *TX
	var kv *KV
	var kvs []*KV
	var err error

	Convey("Check complex iterations", t, func() {

		db, err = Open("memory", &Config{})

		tx, err = db.Begin(true)
		So(err, ShouldBeNil)
		So(tx, ShouldNotBeNil)

		defer tx.Cancel()

		_, err = tx.Put(5, []byte("/kv"), []byte("KV"))
		So(err, ShouldBeNil)
		_, err = tx.Put(5, []byte("/kv/ns"), []byte("NS"))
		So(err, ShouldBeNil)
		_, err = tx.Put(5, []byte("/kv/ns/db"), []byte("DB"))
		So(err, ShouldBeNil)
		_, err = tx.Put(5, []byte("/kv/ns/db/tb1"), []byte("TB1"))
		So(err, ShouldBeNil)
		_, err = tx.Put(5, []byte("/kv/ns/db/tb2"), []byte("TB2"))
		So(err, ShouldBeNil)
		_, err = tx.Put(5, []byte("/kv/ns/db/tb3"), []byte("TB3"))
		So(err, ShouldBeNil)

		kv, err = tx.Get(10, []byte("/kv"))
		So(kv.Val(), ShouldResemble, []byte("KV"))
		kv, err = tx.Get(10, []byte("/kv/ns"))
		So(kv.Val(), ShouldResemble, []byte("NS"))
		kv, err = tx.Get(10, []byte("/kv/ns/db"))
		So(kv.Val(), ShouldResemble, []byte("DB"))

		kvs, err = tx.GetL(10, []byte("/"), 0)
		So(kvs, ShouldHaveLength, 1)
		So(kvs[0].Val(), ShouldResemble, []byte("KV"))

		kvs, err = tx.GetL(10, []byte("/kv"), 0)
		So(kvs, ShouldHaveLength, 1)
		So(kvs[0].Val(), ShouldResemble, []byte("NS"))

		kvs, err = tx.GetL(10, []byte("/kv/ns"), 0)
		So(kvs, ShouldHaveLength, 1)
		So(kvs[0].Val(), ShouldResemble, []byte("DB"))

		kvs, err = tx.GetL(10, []byte("/kv/ns/db"), 0)
		So(kvs, ShouldHaveLength, 3)
		So(kvs[0].Val(), ShouldResemble, []byte("TB1"))
		So(kvs[1].Val(), ShouldResemble, []byte("TB2"))
		So(kvs[2].Val(), ShouldResemble, []byte("TB3"))

		kvs, err = tx.GetP(10, []byte("/"), 0)
		So(kvs, ShouldHaveLength, 6)
		So(kvs[0].Val(), ShouldResemble, []byte("KV"))
		So(kvs[1].Val(), ShouldResemble, []byte("NS"))
		So(kvs[2].Val(), ShouldResemble, []byte("DB"))
		So(kvs[3].Val(), ShouldResemble, []byte("TB1"))
		So(kvs[4].Val(), ShouldResemble, []byte("TB2"))
		So(kvs[5].Val(), ShouldResemble, []byte("TB3"))

		kvs, err = tx.GetP(10, []byte("/k"), 0)
		So(kvs, ShouldHaveLength, 6)
		So(kvs[0].Val(), ShouldResemble, []byte("KV"))
		So(kvs[1].Val(), ShouldResemble, []byte("NS"))
		So(kvs[2].Val(), ShouldResemble, []byte("DB"))
		So(kvs[3].Val(), ShouldResemble, []byte("TB1"))
		So(kvs[4].Val(), ShouldResemble, []byte("TB2"))
		So(kvs[5].Val(), ShouldResemble, []byte("TB3"))

		kvs, err = tx.GetP(10, []byte("/kv"), 0)
		So(kvs, ShouldHaveLength, 6)
		So(kvs[0].Val(), ShouldResemble, []byte("KV"))
		So(kvs[1].Val(), ShouldResemble, []byte("NS"))
		So(kvs[2].Val(), ShouldResemble, []byte("DB"))
		So(kvs[3].Val(), ShouldResemble, []byte("TB1"))
		So(kvs[4].Val(), ShouldResemble, []byte("TB2"))
		So(kvs[5].Val(), ShouldResemble, []byte("TB3"))

		kvs, err = tx.GetP(10, []byte("/kv/ns/db"), 0)
		So(kvs, ShouldHaveLength, 4)
		So(kvs[0].Val(), ShouldResemble, []byte("DB"))
		So(kvs[1].Val(), ShouldResemble, []byte("TB1"))
		So(kvs[2].Val(), ShouldResemble, []byte("TB2"))
		So(kvs[3].Val(), ShouldResemble, []byte("TB3"))

		kvs, err = tx.GetP(10, []byte("/kv/ns/db/"), 0)
		So(kvs, ShouldHaveLength, 3)
		So(kvs[0].Val(), ShouldResemble, []byte("TB1"))
		So(kvs[1].Val(), ShouldResemble, []byte("TB2"))
		So(kvs[2].Val(), ShouldResemble, []byte("TB3"))

		kvs, err = tx.GetP(10, []byte("/kv/ns/db/tb"), 0)
		So(kvs, ShouldHaveLength, 3)
		So(kvs[0].Val(), ShouldResemble, []byte("TB1"))
		So(kvs[1].Val(), ShouldResemble, []byte("TB2"))
		So(kvs[2].Val(), ShouldResemble, []byte("TB3"))

		kvs, err = tx.GetR(10, []byte("/"), []byte("/kv/ns/db/tb~"), 0)
		So(kvs, ShouldHaveLength, 6)
		So(kvs[0].Val(), ShouldResemble, []byte("KV"))
		So(kvs[1].Val(), ShouldResemble, []byte("NS"))
		So(kvs[2].Val(), ShouldResemble, []byte("DB"))
		So(kvs[3].Val(), ShouldResemble, []byte("TB1"))
		So(kvs[4].Val(), ShouldResemble, []byte("TB2"))
		So(kvs[5].Val(), ShouldResemble, []byte("TB3"))

		kvs, err = tx.GetR(10, []byte("/kv/ns/db/tb~"), []byte("/"), 0)
		So(kvs, ShouldHaveLength, 6)
		So(kvs[0].Val(), ShouldResemble, []byte("TB3"))
		So(kvs[1].Val(), ShouldResemble, []byte("TB2"))
		So(kvs[2].Val(), ShouldResemble, []byte("TB1"))
		So(kvs[3].Val(), ShouldResemble, []byte("DB"))
		So(kvs[4].Val(), ShouldResemble, []byte("NS"))
		So(kvs[5].Val(), ShouldResemble, []byte("KV"))

		kvs, err = tx.GetR(10, []byte("/~"), []byte("/"), 0)
		So(kvs, ShouldHaveLength, 6)
		So(kvs[0].Val(), ShouldResemble, []byte("TB3"))
		So(kvs[1].Val(), ShouldResemble, []byte("TB2"))
		So(kvs[2].Val(), ShouldResemble, []byte("TB1"))
		So(kvs[3].Val(), ShouldResemble, []byte("DB"))
		So(kvs[4].Val(), ShouldResemble, []byte("NS"))
		So(kvs[5].Val(), ShouldResemble, []byte("KV"))

		kvs, err = tx.GetR(10, []byte("/kv/ns/db/tb"), []byte("/kv/ns/db/tb~"), 0)
		So(kvs, ShouldHaveLength, 3)
		So(kvs[0].Val(), ShouldResemble, []byte("TB1"))
		So(kvs[1].Val(), ShouldResemble, []byte("TB2"))
		So(kvs[2].Val(), ShouldResemble, []byte("TB3"))

		kvs, err = tx.GetR(10, []byte("/kv/ns/db/tb~"), []byte("/kv/ns/db/tb"), 0)
		So(kvs, ShouldHaveLength, 3)
		So(kvs[0].Val(), ShouldResemble, []byte("TB3"))
		So(kvs[1].Val(), ShouldResemble, []byte("TB2"))
		So(kvs[2].Val(), ShouldResemble, []byte("TB1"))

		kv, err = tx.PutC(15, []byte("/kv/ns/db/tb1"), []byte("TB-1"), []byte("TB1"))
		So(kv.Val(), ShouldResemble, []byte("TB1"))

		kv, err = tx.PutC(15, []byte("/kv/ns/db/tb2"), []byte("TB-2"), []byte("TB2"))
		So(kv.Val(), ShouldResemble, []byte("TB2"))

		kv, err = tx.PutC(15, []byte("/kv/ns/db/tb3"), []byte("TB-3"), []byte("TB3"))
		So(kv.Val(), ShouldResemble, []byte("TB3"))

		kv, err = tx.PutC(15, []byte("/kv/ns/db/tb1"), []byte("TB-4"), []byte("TB4"))
		So(err, ShouldEqual, ErrTxNotExpectedValue)

		kv, err = tx.PutC(15, []byte("/kv/ns/db/tb4"), []byte("TB-4"), nil)
		So(kv.Val(), ShouldResemble, []byte(nil))

		kv, err = tx.PutC(15, []byte("/kv/ns/db/tb5"), []byte("TB-5"), nil)
		So(kv.Val(), ShouldResemble, []byte(nil))

		kvs, err = tx.PutL(20, []byte("/"), []byte("KV-test"), 0)
		So(kvs, ShouldHaveLength, 1)
		So(kvs[0].Key(), ShouldResemble, []byte("/kv"))
		So(kvs[0].Val(), ShouldResemble, []byte("KV"))

		kvs, err = tx.PutL(20, []byte("/kv"), []byte("NS-test"), 0)
		So(kvs, ShouldHaveLength, 1)
		So(kvs[0].Key(), ShouldResemble, []byte("/kv/ns"))
		So(kvs[0].Val(), ShouldResemble, []byte("NS"))

		kvs, err = tx.PutL(20, []byte("/kv/ns"), []byte("DB-test"), 0)
		So(kvs, ShouldHaveLength, 1)
		So(kvs[0].Key(), ShouldResemble, []byte("/kv/ns/db"))
		So(kvs[0].Val(), ShouldResemble, []byte("DB"))

		kvs, err = tx.PutL(25, []byte("/kv/ns/db"), []byte("TB-test"), 0)
		So(kvs, ShouldHaveLength, 5)
		So(kvs[0].Key(), ShouldResemble, []byte("/kv/ns/db/tb1"))
		So(kvs[0].Val(), ShouldResemble, []byte("TB-1"))
		So(kvs[1].Key(), ShouldResemble, []byte("/kv/ns/db/tb2"))
		So(kvs[1].Val(), ShouldResemble, []byte("TB-2"))
		So(kvs[2].Key(), ShouldResemble, []byte("/kv/ns/db/tb3"))
		So(kvs[2].Val(), ShouldResemble, []byte("TB-3"))
		So(kvs[3].Key(), ShouldResemble, []byte("/kv/ns/db/tb4"))
		So(kvs[3].Val(), ShouldResemble, []byte("TB-4"))
		So(kvs[4].Key(), ShouldResemble, []byte("/kv/ns/db/tb5"))
		So(kvs[4].Val(), ShouldResemble, []byte("TB-5"))

		kvs, err = tx.DelL(25, []byte("/kv/ns/db"), 0)
		So(kvs, ShouldHaveLength, 5)
		So(kvs[0].Val(), ShouldResemble, []byte("TB-test"))
		So(kvs[1].Val(), ShouldResemble, []byte("TB-test"))
		So(kvs[2].Val(), ShouldResemble, []byte("TB-test"))
		So(kvs[3].Val(), ShouldResemble, []byte("TB-test"))
		So(kvs[4].Val(), ShouldResemble, []byte("TB-test"))

		kvs, err = tx.GetP(30, []byte("/"), 0)
		So(kvs, ShouldHaveLength, 3)
		So(kvs[0].Val(), ShouldResemble, []byte("KV-test"))
		So(kvs[1].Val(), ShouldResemble, []byte("NS-test"))
		So(kvs[2].Val(), ShouldResemble, []byte("DB-test"))

		kvs, err = tx.DelP(30, []byte("/"), 0)
		So(kvs, ShouldHaveLength, 8)
		So(kvs[0].Val(), ShouldResemble, []byte("KV-test"))
		So(kvs[1].Val(), ShouldResemble, []byte("NS-test"))
		So(kvs[2].Val(), ShouldResemble, []byte("DB-test"))
		So(kvs[3].Val(), ShouldResemble, []byte(nil))
		So(kvs[4].Val(), ShouldResemble, []byte(nil))
		So(kvs[5].Val(), ShouldResemble, []byte(nil))
		So(kvs[6].Val(), ShouldResemble, []byte(nil))
		So(kvs[7].Val(), ShouldResemble, []byte(nil))

		kvs, err = tx.GetP(30, []byte("/"), 0)
		So(kvs, ShouldHaveLength, 0)

		kvs, err = tx.PutP(35, []byte("/"), []byte("TEST"), 0)
		So(kvs, ShouldHaveLength, 8)
		So(kvs[0].Val(), ShouldResemble, []byte(nil))
		So(kvs[1].Val(), ShouldResemble, []byte(nil))
		So(kvs[2].Val(), ShouldResemble, []byte(nil))
		So(kvs[3].Val(), ShouldResemble, []byte(nil))
		So(kvs[4].Val(), ShouldResemble, []byte(nil))
		So(kvs[5].Val(), ShouldResemble, []byte(nil))
		So(kvs[6].Val(), ShouldResemble, []byte(nil))
		So(kvs[7].Val(), ShouldResemble, []byte(nil))

		kvs, err = tx.DelP(35, []byte("/"), 0)
		So(kvs, ShouldHaveLength, 8)
		So(kvs[0].Val(), ShouldResemble, []byte("TEST"))
		So(kvs[1].Val(), ShouldResemble, []byte("TEST"))
		So(kvs[2].Val(), ShouldResemble, []byte("TEST"))
		So(kvs[3].Val(), ShouldResemble, []byte("TEST"))
		So(kvs[4].Val(), ShouldResemble, []byte("TEST"))
		So(kvs[5].Val(), ShouldResemble, []byte("TEST"))
		So(kvs[6].Val(), ShouldResemble, []byte("TEST"))
		So(kvs[7].Val(), ShouldResemble, []byte("TEST"))

		kvs, err = tx.PutR(35, []byte("/"), []byte("/kv/ns/db/tb~"), []byte("TEMP"), 0)
		So(kvs, ShouldHaveLength, 8)
		So(kvs[0].Val(), ShouldResemble, []byte(nil))
		So(kvs[1].Val(), ShouldResemble, []byte(nil))
		So(kvs[2].Val(), ShouldResemble, []byte(nil))
		So(kvs[3].Val(), ShouldResemble, []byte(nil))
		So(kvs[4].Val(), ShouldResemble, []byte(nil))
		So(kvs[5].Val(), ShouldResemble, []byte(nil))
		So(kvs[6].Val(), ShouldResemble, []byte(nil))
		So(kvs[7].Val(), ShouldResemble, []byte(nil))

		kvs, err = tx.DelR(35, []byte("/"), []byte("/kv/ns/db/tb~"), 0)
		So(kvs, ShouldHaveLength, 8)
		So(kvs[0].Val(), ShouldResemble, []byte("TEMP"))
		So(kvs[1].Val(), ShouldResemble, []byte("TEMP"))
		So(kvs[2].Val(), ShouldResemble, []byte("TEMP"))
		So(kvs[3].Val(), ShouldResemble, []byte("TEMP"))
		So(kvs[4].Val(), ShouldResemble, []byte("TEMP"))
		So(kvs[5].Val(), ShouldResemble, []byte("TEMP"))
		So(kvs[6].Val(), ShouldResemble, []byte("TEMP"))
		So(kvs[7].Val(), ShouldResemble, []byte("TEMP"))

		kvs, err = tx.PutR(35, []byte("/kv/ns/db/tb~"), []byte("/"), []byte("TEMP"), 0)
		So(kvs, ShouldHaveLength, 8)
		So(kvs[0].Val(), ShouldResemble, []byte(nil))
		So(kvs[1].Val(), ShouldResemble, []byte(nil))
		So(kvs[2].Val(), ShouldResemble, []byte(nil))
		So(kvs[3].Val(), ShouldResemble, []byte(nil))
		So(kvs[4].Val(), ShouldResemble, []byte(nil))
		So(kvs[5].Val(), ShouldResemble, []byte(nil))
		So(kvs[6].Val(), ShouldResemble, []byte(nil))
		So(kvs[7].Val(), ShouldResemble, []byte(nil))

		kvs, err = tx.DelR(35, []byte("/kv/ns/db/tb~"), []byte("/"), 0)
		So(kvs, ShouldHaveLength, 8)
		So(kvs[0].Val(), ShouldResemble, []byte("TEMP"))
		So(kvs[1].Val(), ShouldResemble, []byte("TEMP"))
		So(kvs[2].Val(), ShouldResemble, []byte("TEMP"))
		So(kvs[3].Val(), ShouldResemble, []byte("TEMP"))
		So(kvs[4].Val(), ShouldResemble, []byte("TEMP"))
		So(kvs[5].Val(), ShouldResemble, []byte("TEMP"))
		So(kvs[6].Val(), ShouldResemble, []byte("TEMP"))
		So(kvs[7].Val(), ShouldResemble, []byte("TEMP"))

		kvs, err = tx.GetP(35, []byte("/"), 0)
		So(kvs, ShouldHaveLength, 0)

		kv, err = tx.Put(40, []byte("/kv"), []byte("KV-test"))
		So(kv.Val(), ShouldResemble, []byte(nil))

		kv, err = tx.Get(40, []byte("/kv"))
		So(kv.Val(), ShouldResemble, []byte("KV-test"))

		kv, err = tx.Del(40, []byte("/kv"))
		So(kv.Val(), ShouldResemble, []byte("KV-test"))

		kv, err = tx.Put(40, []byte("/kv"), []byte("KV-test"))
		So(kv.Val(), ShouldResemble, []byte(nil))

		kv, err = tx.DelC(40, []byte("/kv"), []byte("KV-tester"))
		So(err, ShouldResemble, ErrTxNotExpectedValue)

		kv, err = tx.DelC(40, []byte("/kv"), []byte("KV-test"))
		So(kv.Val(), ShouldResemble, []byte("KV-test"))

		kv, err = tx.Get(40, []byte("/kv"))
		So(kv.Val(), ShouldResemble, []byte(nil))

		kv, err = tx.Put(45, []byte("/kv/ns/db/tbx"), []byte("TBX"))
		So(kv.Val(), ShouldResemble, []byte(nil))

		kv, err = tx.Get(45, []byte("/kv/ns/db/tbx"))
		So(kv.Val(), ShouldResemble, []byte("TBX"))

		kv, err = tx.Clr([]byte("/kv/ns/db/tbx"))
		So(kv.Val(), ShouldResemble, []byte("TBX"))

		kv, err = tx.Put(45, []byte("/kv/ns/db/tbx"), []byte("TBX"))
		So(kv.Val(), ShouldResemble, []byte(nil))

		kv, err = tx.Put(45, []byte("/kv/ns/db/tbx"), []byte("TBX"))
		So(kv.Val(), ShouldResemble, []byte("TBX"))

		kv, err = tx.Put(45, []byte("/kv/ns/db/tbx"), []byte("TBX"))
		So(kv.Val(), ShouldResemble, []byte("TBX"))

		kvs, err = tx.ClrL([]byte("/kv/ns"), 0)
		So(kvs, ShouldHaveLength, 1)
		So(kvs[0].Val(), ShouldResemble, []byte(nil))

		kvs, err = tx.ClrP([]byte("/kv/ns/db/tb"), 0)
		So(kvs, ShouldHaveLength, 6)
		So(kvs[0].Val(), ShouldResemble, []byte(nil))
		So(kvs[1].Val(), ShouldResemble, []byte(nil))
		So(kvs[2].Val(), ShouldResemble, []byte(nil))

		kvs, err = tx.ClrP([]byte("/"), 0)
		So(kvs, ShouldHaveLength, 2)
		So(kvs[0].Val(), ShouldResemble, []byte(nil))

		_, err = tx.Put(5, []byte("/kv"), []byte("KV"))
		So(err, ShouldBeNil)
		_, err = tx.Put(5, []byte("/kv/ns"), []byte("NS"))
		So(err, ShouldBeNil)
		_, err = tx.Put(5, []byte("/kv/ns/db"), []byte("DB"))
		So(err, ShouldBeNil)
		_, err = tx.Put(5, []byte("/kv/ns/db/tb1"), []byte("TB1"))
		So(err, ShouldBeNil)
		_, err = tx.Put(5, []byte("/kv/ns/db/tb2"), []byte("TB2"))
		So(err, ShouldBeNil)
		_, err = tx.Put(5, []byte("/kv/ns/db/tb3"), []byte("TB3"))
		So(err, ShouldBeNil)

		kvs, err = tx.ClrR([]byte("/"), []byte("/kv/ns/db/tb~"), 0)
		So(kvs, ShouldHaveLength, 6)
		So(kvs[0].Val(), ShouldResemble, []byte("KV"))
		So(kvs[1].Val(), ShouldResemble, []byte("NS"))
		So(kvs[2].Val(), ShouldResemble, []byte("DB"))
		So(kvs[3].Val(), ShouldResemble, []byte("TB1"))
		So(kvs[4].Val(), ShouldResemble, []byte("TB2"))
		So(kvs[5].Val(), ShouldResemble, []byte("TB3"))

		_, err = tx.Put(5, []byte("/kv"), []byte("KV"))
		So(err, ShouldBeNil)
		_, err = tx.Put(5, []byte("/kv/ns"), []byte("NS"))
		So(err, ShouldBeNil)
		_, err = tx.Put(5, []byte("/kv/ns/db"), []byte("DB"))
		So(err, ShouldBeNil)
		_, err = tx.Put(5, []byte("/kv/ns/db/tb1"), []byte("TB1"))
		So(err, ShouldBeNil)
		_, err = tx.Put(5, []byte("/kv/ns/db/tb2"), []byte("TB2"))
		So(err, ShouldBeNil)
		_, err = tx.Put(5, []byte("/kv/ns/db/tb3"), []byte("TB3"))
		So(err, ShouldBeNil)

		kvs, err = tx.ClrR([]byte("/kv/ns/db/tb~"), []byte("/"), 0)
		So(kvs, ShouldHaveLength, 6)
		So(kvs[0].Val(), ShouldResemble, []byte("TB3"))
		So(kvs[1].Val(), ShouldResemble, []byte("TB2"))
		So(kvs[2].Val(), ShouldResemble, []byte("TB1"))
		So(kvs[3].Val(), ShouldResemble, []byte("DB"))
		So(kvs[4].Val(), ShouldResemble, []byte("NS"))
		So(kvs[5].Val(), ShouldResemble, []byte("KV"))

		So(tx.Commit(), ShouldBeNil)

	})

}
