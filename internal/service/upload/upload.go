package upload

import "github.com/xuri/excelize/v2"

type Service interface {
	UploadBlockOne([]byte data, int contentType, fileExtension) ([]byte, error)
}

type service struct {
}

func NewService() Service {
	return &service{}
}

//Принимает на вход предобработанный первый блок данных
//Возвращает этот же блок с вычисленными данными
func (s *service) UploadBlockOne(data []byte, int contentType, fileExtension) ([]byte, error) {

}
