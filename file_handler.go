package tesis

import (
	"os"
)

func NewFileHandler(f string) (h *FileHandler, e error) {
	h = new(FileHandler)
	h.file, e = os.Open(f)
	if os.IsNotExist(e) {
		h.file, e = os.Create(f)
	}
	if e == nil {
		h.wfle, e = os.Create(f + "~")
	}
	if h.file != nil && e != nil {
		h.file.Close()
	}
	return
}

type FileHandler struct {
	file *os.File
	wfle *os.File
}

func (h *FileHandler) Read(p []byte) (n int, e error) {
	n, e = h.file.Read(p)
	return
}

func (h *FileHandler) Write(p []byte) (n int, e error) {
	n, e = h.wfle.Write(p)
	h.wfle.Sync()
	return
}

func (h *FileHandler) Close() (e error) {
	h.file.Close()
	h.wfle.Close()
	e = os.Rename(h.wfle.Name(), h.file.Name())
	return
}
