package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/shiena/ansicolor"
	"github.com/wsxiaoys/terminal/color"
)

// testCmd represents the add command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "test",
	Long:  `test`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("test called")
		if len(args) > 0 {
			for i, _ := range args {
				fmt.Println(args[i])
			}
			testFiles(args[0])
		}
	},
}

func testFiles(path string) {

	list := readFile(path)
	if len(list) > 0 {
		for _, l := range list {
			/*
					fmt.Printf("\n  Id: %d\n", l.Id)
					fmt.Printf("  Filename: %v\n", l.Filename)
					fmt.Printf("  Fullpath: %v\n", l.Fullpath)
					fmt.Printf("  Modtime: %v\n", l.Modtime)
					fmt.Printf("  Size: %d\n", l.Size)
					fmt.Printf("  Rw: %v\n", l.Rw)
					fmt.Printf("  Mode: %d\n", l.Mode)
					fmt.Printf("  Modestring: %v\n", l.Modestring)
					fmt.Printf("  Sha1: %v\n", l.Sha1)

				fmt.Printf("%5d: %v, %v\n", l.Id, testFile(l), l.Fullpath)
			*/

			var result string
			t, m := testFile(l)
			if t {
				result = fmt.Sprintf("%5d| @{g}%-12v@{|}| %v\n", l.Id, m, l.Fullpath)
			} else {
				result = fmt.Sprintf("%5d| @{r}%-12v@{|}| %v\n", l.Id, m, l.Fullpath)
			}
			w := ansicolor.NewAnsiColorWriter(os.Stdout)
			color.Fprintf(w, result)

		}
	}

}

func readFile(path string) []Item {
	p, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	p = filepath.Clean(p)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	raw, err := ioutil.ReadFile(p)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	var f = File{}
	json.Unmarshal(raw, &f)

	fmt.Printf("\nTitle: %v\n", f.Title)
	fmt.Printf("Num: %d\n", f.Num)
	fmt.Printf("Message: %v\n", f.Message)
	if len(f.List) > 0 {
		return f.List
	}
	return nil
}

func testFile(item Item) (bool, string) {
	if !isExist(item.Fullpath) {
		return false, "exist: false"
	}
	if !isMatch(item.Fullpath, item.Sha1) {
		return false, "match: false"
	}
	return true, "match: true"
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func isMatch(path string, except string) bool {
	sha1 := getFileHash(path)
	if except != sha1 {
		//fmt.Printf("except: %v, given: %v\n", except, sha1)
		return false
	}
	return true
}

func init() {
	RootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
