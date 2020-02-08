package cmd

import (
	//"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"
	"time"
)

type Item struct {
	Id         int         `json:"id"`
	Filename   string      `json:"filename"`
	Fullpath   string      `json:"fullpath"`
	Modtime    time.Time   `json:"modtime"`
	Size       int64       `json:"size"`
	Rw         bool        `json:"rw"`
	Mode       os.FileMode `json:"mode"`
	Modestring string      `json:"modestring"`
	Sha1       string      `json:"sha1"`
	Except     string      `json:"except"`
}

type File struct {
	BaseDir        string   `json:"baseDir"`
	Num            int      `json:"num"`
	WarningOverlay []string `json:"warningOverlay"`
	WarningMaxPath int      `json:"warningMaxPath"`
	Message        string   `json:"message"`
	List           []Item   `json:"list"`
}

func getFileHash(path string) string {
	file, err := os.Open(path)
	if err != nil {
		return "file not found"
	}
	defer file.Close()

	hash := sha1.New()
	//hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		//log.Fatal(err)
		log.Println(err)
	}
	sum := hash.Sum(nil)

	return hex.EncodeToString(sum)
}
