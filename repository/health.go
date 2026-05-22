package repository

func (repoStruct *Repository) CheckDB() error {

	sqlDB, err := repoStruct.db.DB()

	if err != nil {
		return err
	}

	return sqlDB.Ping()

}
