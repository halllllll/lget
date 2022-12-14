package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/halllllll/lget"
)

func main() {
	cd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	loginInfo := &lget.LoginInfo{
		Host:    "xxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		AdminId: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		AdminPw: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
	}
	l_get := lget.NewLget()

	opened_l_get, err := l_get.Login(loginInfo)
	if err != nil {
		panic(err)
	}

	// ex 2022-08-21 10:00:00 -> 1661043600
	// ex 2022-08-21 13:00:00 -> 1661054400
	downloadFileUrl, err := opened_l_get.GetLog(1661043600, 1661054400)
	if err != nil {
		panic(err)
	}
	fmt.Printf("download file url: %s\n", downloadFileUrl)

	logResult := make(chan []byte)

	go func() {
		for {
			rawData, err := opened_l_get.Download(downloadFileUrl)
			if err != nil {
				panic(err)
			}
			logResult <- rawData
		}
	}()

	rawData := <-logResult
	if err := saveFile(rawData, filepath.Join(cd, "done.csv")); err != nil {
		panic(err)
	}
}

func saveFile(data []byte, path string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}
