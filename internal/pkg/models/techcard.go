package models

const (
	FileExtension_XSLX = iota
	FileExtension_CSV
)

const (
	ContentType_ZABOY_PRESSURE_TEMPERATURE = iota
	ContentType_PIPE_CASING_LINEAR_PRESSURE
	ContentType_DEBIT
	ContentType_INCLINOMETRY
)

message UploadRequest {
FileExtension file_extension = 1;
ContentType content_type = 2;
string table = 3;
}

enum ResponseStatus{
ResponseStatus_OK = 0;
ResponseStatus_BAD_REQUSET = 1;
ResponseStatus_INTERNAL_ERROR = 2;
}

message UploadResponse {
ResponseStatus status = 1;
optional string message = 2;
}
