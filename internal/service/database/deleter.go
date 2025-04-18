package database

type deleter interface {
	DeleteOilField() error
	DeleteProductiveHorizon() error
	DeleteInstrumentType() error
}

func (d dbService) DeleteOilField() error {
	//TODO implement me
	panic("implement me")
}

func (d dbService) DeleteProductiveHorizon() error {
	//TODO implement me
	panic("implement me")
}

func (d dbService) DeleteInstrumentType() error {
	//TODO implement me
	panic("implement me")
}
