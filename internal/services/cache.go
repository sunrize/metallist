package services

import (
    "encoding/json"
    "errors"
    "fmt"
    "os"

    "github.com/dgraph-io/badger"
)

// SaveJSON saves a JSON object to the database under the given key.
func SaveJSON(db *badger.DB, key string, value interface{}) error {
    data, err := json.Marshal(value)
    if err != nil {
        return errors.New("Error marshaling JSON: " + err.Error())
    }

    err = db.Update(func(txn *badger.Txn) error {
        item, err := txn.Get([]byte(key))
        if err == badger.ErrKeyNotFound {
            // Item does not exist, create it
            err = txn.Set([]byte(key), data)
        } else if err != nil {
            return err
        } else {
            // Item exists, update it
            err = item.Delete()
            if err != nil {
                return err
            }
            err = txn.Set([]byte(key), data)
        }
        return err
    })
	
    if err != nil {
        fmt.Println("Error saving data: ", err)
        return err
    }

    return nil
}

// LoadJSON loads a JSON object from the database under the given key.
func LoadJSON(db *badger.DB, key string, target interface{}) error {
    var data []byte
    err := db.View(func(txn *badger.Txn) error {
        var err error
        data, err = txn.Get([]byte(key))
        return err
    })
    if err != nil {
        if errors.Is(err, badger.ErrKeyNotFound) {
            fmt.Fprintf(os.Stderr, "Key '%s' not found in database\n", key)
            os.Exit(1)
        }
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    if len(data) == 0 {
        // Key exists but is empty
        return nil
    }

    // Unmarshal JSON into target
    return json.Unmarshal(data, target)
}

// OpenBadger opens a new Badger database at the given path.
func OpenBadger(path string) (*badger.DB, error) {
    return badger.Open(badger.DefaultOptions(path))
}