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
		var pass, fail int
		w := ansicolor.NewAnsiColorWriter(os.Stdout)
		for _, l := range list {
			var result string
			r, msg := testFile(l)
			if r {
				result = fmt.Sprintf("%5d| @{g}%-12v@{|}| %v\n", l.Id, msg, l.Fullpath)
				pass++
			} else {
				result = fmt.Sprintf("%5d| @{r}%-12v@{|}| %v\n", l.Id, msg, l.Fullpath)
				fail++
			}
			color.Fprintf(w, result)
		}
		color.Fprintf(w, "\n@{g}PASS: %d@{|} / @{r}FAIL: %d@{|}\n", pass, fail)
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
