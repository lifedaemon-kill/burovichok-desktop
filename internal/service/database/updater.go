package database

type updater interface {
	UpdateInstrumentType() error
	UpdateProductiveHorizon() error
	UpdateOilField() error
}

func (d dbService) UpdateInstrumentType() error {
	//TODO implement me
	panic("implement me")
}

func (d dbService) UpdateProductiveHorizon() error {
	//TODO implement me
	panic("implement me")
}

func (d dbService) UpdateOilField() error {
	//TODO implement me
	panic("implement me")
}
