package filemanager

import (
	"errors"
	"os"
	"path/filepath"
)

var ErrEOF = errors.New("EOF")

type FileSaver struct {
	filePath string
	outFile  *os.File
}

func NewFileSaver() *FileSaver {
	return &FileSaver{}
}

func (f *FileSaver) SetFile(fileName string, path string) error {
	f.filePath = filepath.Join(path, fileName)
	file, err := os.Create(f.filePath)
	if err != nil {
		return err
	}
	f.outFile = file
	return nil
}

func (f *FileSaver) IsFileSet() bool {
	return f.filePath != ""
}

func (f *FileSaver) Write(chunk []byte) error {
	if f.outFile == nil {
		return nil
	}
	_, err := f.outFile.Write(chunk)
	return err
}

func (f *FileSaver) Close() error {
	return f.outFile.Close()
}

type FileDeleter struct {
	filePath string
}

func NewFileDeleter() *FileDeleter {
	return &FileDeleter{}
}

func (f *FileDeleter) SetFile(fileName string, path string) {
	f.filePath = filepath.Join(path, fileName)
}

func (f *FileDeleter) Delete() error {
	return os.Remove(f.filePath)
}

type FileReader struct {
	filePath  string
	chunkSize int
	buf       []byte
	file      *os.File
	read      int
}

func NewFileReader(chunkSize int) *FileReader {
	return &FileReader{chunkSize: chunkSize, buf: make([]byte, chunkSize)}
}

func (f *FileReader) SetFile(fileName string, path string) error {
	f.filePath = filepath.Join(path, fileName)
	file, err := os.Open(f.filePath)
	if err != nil {
		return err
	}
	f.file = file
	return nil
}

func (f *FileReader) Next() bool {
	read, err := f.file.Read(f.buf)
	f.read = read
	if err != nil {
		return false
	}
	return true
}

func (f *FileReader) Data() []byte {
	return f.buf[:f.read]
}

func (f *FileReader) Close() error {
	return f.file.Close()
}
