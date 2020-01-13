package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add path",
	Long:  `add files test configuration`,
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
			fp, err := filepath.Abs(p)
			if err != nil {
				log.Fatal(err)
			}
			fp = filepath.Clean(fp)
			if err != nil {
				log.Fatal(err)
			}
			f.List = append(f.List, Item{
				i + 1,
				info.Name(),
				fp,
				info.ModTime(),
				info.Size(),
				info.Mode().IsRegular(),
				info.Mode().Perm(),
				info.Mode().String(),
				getFileHash(p),
				"match"})
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
