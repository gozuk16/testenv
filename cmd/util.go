package cmd

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
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
}

type File struct {
	Title   string `json:"title"`
	Num     int    `json:"num"`
	Message string `json:"message"`
	List    []Item `json:"list"`
}

func getFileHash(path string) string {
	f := strings.NewReader(path)
	hash := sha1.New()
	if _, err := io.Copy(hash, f); err != nil {
		log.Fatal(err)
	}
	sum := hash.Sum(nil)

	return fmt.Sprintf("%x", sum)
}
