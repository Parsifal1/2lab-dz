package fileMeta

import "os"

//FileMeta - единица передачи метаданных файла
type FileMetaStruct struct {
	Name string // `json: "filename"`
}

func FileMeta(info os.FileInfo) (meta *FileMetaStruct, err error) {
	meta = &FileMetaStruct{
		Name: info.Name(),
	}
	return
}
