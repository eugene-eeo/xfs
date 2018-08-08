package libxfs

import bolt "github.com/coreos/bbolt"

var hashes_bucket = []byte("xfs-path-to-hash")

type Path []byte
type Hash []byte

var EmptyHash = Hash("")

func InitDB(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(hashes_bucket)
		return err
	})
}

func SaveHash(db *bolt.DB, path Path, hash Hash) error {
	return db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(hashes_bucket).Put(path, hash)
	})
}

func GetHash(db *bolt.DB, path Path) (Hash, error) {
	hash := EmptyHash
	err := db.View(func(tx *bolt.Tx) error {
		v := tx.Bucket(hashes_bucket).Get(path)
		if len(v) > 0 {
			hash = v
		}
		return nil
	})
	return hash, err
}
