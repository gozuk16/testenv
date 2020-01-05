package cmd

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
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
	//List    map[int]Item `json:"list"`
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add called")
		if len(args) > 0 {
			for i, _ := range args {
				fmt.Println(args[i])
			}
			seachFile(args[0])
		}
	},
}

func seachFile(path string) {
	//fmt.Println(filepath.Dir(filepath.Clean(path)))
	var f = File{}
	origin, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}
	origin = filepath.Clean(origin)
	if err != nil {
		log.Fatal(err)
	}
	f.Title = origin
	f.List = make([]Item, 0)
	i := 0
	err = filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			//fmt.Printf("%s %s %d %s %s\n", info.Name(), p, info.Size(), info.ModTime(), getFileHash(p))
			fp, err := filepath.Abs(p)
			if err != nil {
				log.Fatal(err)
			}
			fp = filepath.Clean(fp)
			if err != nil {
				log.Fatal(err)
			}
			//f.List[i] = Item{
			f.List = append(f.List, Item{
				i + 1,
				info.Name(),
				fp,
				info.ModTime(),
				info.Size(),
				info.Mode().IsRegular(),
				info.Mode().Perm(),
				info.Mode().String(),
				getFileHash(p)})
			i++
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	if i > 0 {
		f.Num = i
		bytes, err := json.MarshalIndent(&f, "", "    ")
		if err != nil {
			fmt.Println("Err: ", err)
		}
		jsonstring := string(bytes)
		fmt.Println(jsonstring)
	}
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

func init() {
	RootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
