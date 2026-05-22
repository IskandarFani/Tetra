package services

func (servStruct *Services) CheckDB() error {
	return servStruct.repo.CheckDB()
}
