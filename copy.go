package main

import (
	"os"
)

type copyaction struct {
	file *os.File
	md5ac md5action
}

func (a *copyaction) Init() {
	var err error
	a.file, err = os.OpenFile(g_copyfilename, os.O_CREATE | os.O_RDWR, 0755)
	if err != nil {
		uiPrintf2(0, "%s\n", err)
	}
	a.md5ac.Init()
	return
}

func (a *copyaction) Run(buf[] byte, abspos, tcnt, n uint64, lastbit bool) (uint64, error) {
	l, err := a.md5ac.Run(buf, abspos, tcnt, n, lastbit)
	if err != nil {
		a.file.Close()
		return 0, err
	}
	_, err = a.file.Write(buf[0:l])
	if err  != nil {
		a.file.Close()
		return 0, err
	}
	if lastbit {
		a.file.Close()
	}
	return l, nil
}
