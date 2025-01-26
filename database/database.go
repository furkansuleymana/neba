package database

import (
	"encoding/json"
	"fmt"

	"github.com/furkansuleymana/neba/database/models"
	"go.etcd.io/bbolt"
)

// Open opens a BoltDB database at the specified path and ensures that the specified bucket exists.
// If the database or bucket does not exist, they will be created.
//
// Parameters:
//   - databasePath: The file path to the BoltDB database.
//   - bucketName: The name of the bucket to ensure exists.
//
// Returns:
//   - *bbolt.DB: A pointer to the opened BoltDB database.
//   - error: An error if the database or bucket could not be opened or created.
func Open(databasePath string, bucketName string) (*bbolt.DB, error) {
	db, err := bbolt.Open(databasePath, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("open database, %v", err)
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("set up bucket, %v", err)
	}

	return db, nil
}

// Update updates the specified bucket in the bbolt database with the given AxisDevice.
// It marshals the AxisDevice to JSON and stores it in the bucket using the device's SerialNumber as the key.
// If the bucket does not exist, it returns an error.
//
// Parameters:
//   - db: A pointer to the bbolt.DB instance.
//   - bucketName: The name of the bucket to update.
//   - device: The AxisDevice to store in the bucket.
//
// Returns:
//   - error: An error if the update operation fails, otherwise nil.
func Update(db *bbolt.DB, bucketName string, device models.AxisDevice) error {
	return db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", bucketName)
		}
		encoded, err := json.Marshal(device)
		if err != nil {
			return fmt.Errorf("marshal device to JSON: %v", err)
		}
		return bucket.Put([]byte(device.SerialNumber), encoded)
	})
}

// View retrieves an AxisDevice from the specified bucket in the BoltDB database
// using the provided serial number. It returns the AxisDevice and any error encountered.
//
// Parameters:
//   - db: A pointer to the BoltDB database.
//   - bucketName: The name of the bucket to retrieve the device from.
//   - serialNumber: The serial number of the device to retrieve.
//
// Returns:
//   - A pointer to the AxisDevice if found, or nil if not found.
//   - An error if the bucket or device is not found, or if there is an issue unmarshalling the JSON data.
func View(db *bbolt.DB, bucketName string, serialNumber string) (*models.AxisDevice, error) {
	var device models.AxisDevice
	err := db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", bucketName)
		}
		value := bucket.Get([]byte(serialNumber))
		if value == nil {
			return fmt.Errorf("device %s not found", serialNumber)
		}
		if err := json.Unmarshal(value, &device); err != nil {
			return fmt.Errorf("unmarshal JSON: %v", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &device, nil
}

// CloseDB closes the given bbolt database.
// It ensures that all database resources are properly released.
//
// Parameters:
//   - db: A pointer to the bbolt.DB instance that needs to be closed.
func CloseDB(db *bbolt.DB) {
	db.Close()
}
