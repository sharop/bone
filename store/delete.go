package store

//Delete primitive delete of badger
func (db *Database) Delete(key string) (deleted bool) {
	txn := db.Service.Badger().NewTransaction(true)
	err := txn.Delete(makeKey(key))
	if err != nil {
		return
	}
	err = txn.Commit()
	return err == nil
}

func (db *Database) DeleteBadgerhold(key string) error {
	var result Item
	err := db.Service.Delete(key, result)
	if err != nil {
		return err
	}
	return nil
}
