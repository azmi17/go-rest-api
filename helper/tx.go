package helper

import "database/sql"

func CommitOrRollBack(tx *sql.Tx) {

	err := recover()

	if err != nil {
		errorRollback := tx.Rollback()
		PanicIfError(errorRollback)
		panic(err)
	} else {
		errorCommit := tx.Commit()
		PanicIfError(errorCommit)
	}
	/*
		defer func() {

			statement dibawah di wrap ke helper saja
			supaya tidak diterapkan di semua method
			karena terkesan verbose dan rendundant
			err := recover()
			if err != nil {
				tx.Rollback()
					panic(err)
				} else {
					tx.Commit()
			}

		}()
	*/
}
