package rixxdb

import "errors"

// These errors can occur when opening or calling methods on a DB.
var (
	// ErrDbClosed occurs when a DB is accessed after it is closed.
	ErrDbClosed = errors.New("DB is not open")

	// ErrDbMemoryOnly occurs when persisting an in-memory DB.
	ErrDbMemoryOnly = errors.New("DB is memory-only")

	// ErrDbAlreadySyncing occurs when a sync is already in progress.
	ErrDbAlreadySyncing = errors.New("DB is already syncing")

	// ErrDbAlreadyShrinking occurs when a shrink is already in progress.
	ErrDbAlreadyShrinking = errors.New("DB is already shrinking")

	// ErrDbFileContentsInvalid occurs when the file is invalid or currupted.
	ErrDbFileContentsInvalid = errors.New("DB file contents invalid")

	// ErrDbInvalidEncryptionKey occurs when the provided encryption key is invalid.
	ErrDbInvalidEncryptionKey = errors.New("DB encryption key invalid")
)

// These errors can occur when beginning or calling methods on a TX.
var (
	// ErrTxClosed occurs when cancelling or committing a closed transaction.
	ErrTxClosed = errors.New("TX is closed")

	// ErrTxNotWritable occurs when writing or committing a read-only transaction.
	ErrTxNotWritable = errors.New("TX is not writable")

	// ErrTxNotEditable occurs when calling manually closing a managed transaction.
	ErrTxNotEditable = errors.New("TX is not editable")

	// ErrKvNotExpectedValue occurs when using a nil key in put, select, or delete methods.
	ErrTxKeyCanNotBeNil = errors.New("TX key can not be nil")

	// ErrKvNotExpectedValue occurs when conditionally putting or deleting a key-value item.
	ErrTxNotExpectedValue = errors.New("KV val is not expected value")
)
