package util

import (
	"os"
	"io"
)

type FileWriter struct {
	Reader io.Reader
	FileName string
	file *os.File
	fileOpened bool
}

func (fw *FileWriter) Read(out []byte) (i int, e error) {
	i, e = fw.Reader.Read(out)
	if i > 0 && e == nil {
		if fw.fileOpened == false {
			os.Remove(fw.FileName)
			fw.file, e = os.OpenFile(fw.FileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			fw.fileOpened = e == nil
		}
		if fw.fileOpened {
			fw.file.Write(out[:i])
		}
	}
	return
}

func (fw *FileWriter) Close() {
	if fw.fileOpened && fw.file != nil {
		fw.file.Close()
		fw.file = nil
		fw.fileOpened = false
	}
}
