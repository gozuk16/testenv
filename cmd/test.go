package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

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

var overlayExtensions []string
var maxPathThreshold int

func testFiles(path string) {
	baseDir, list := readFile(path)
	if len(list) > 0 {
		var pass, fail, warning int
		w := ansicolor.NewAnsiColorWriter(os.Stdout)
		for _, l := range list {
			var result string
			r, msg := testFile(l)
			if r {
				result = fmt.Sprintf("%5d| @{g}%-16v@{|}| %v\n", l.Id, msg, l.Fullpath)
				pass++
			} else {
				result = fmt.Sprintf("%5d| @{r}%-16v@{|}| %v\n", l.Id, msg, l.Fullpath)
				fail++
			}
			color.Fprintf(w, result)

			resMsg := formatTestResultMessage("overlay") + ": warn"
			if shouldOverlayTest(l.Filename) {
				msgs := testOverlay(baseDir, l)
				if msgs != nil {
					for i, msg := range msgs {
						if i == len(msgs)-1 {
							msg = " └─ " + msg
						} else {
							msg = " ├─ " + msg
						}
						result = fmt.Sprintf("%5v| @{y}%-16v@{|}| %v\n", "", resMsg, msg)
						warning++
						color.Fprintf(w, result)
					}
				}
			}

			resMsg = formatTestResultMessage("max path") + ": warn"
			if maxPathThreshold > 0 {
				r, len := testMaxPath(l.Fullpath)
				msg := fmt.Sprintf(" └─ path len:%-3d, over: %-3d", len, len-maxPathThreshold)
				if !r {
					result = fmt.Sprintf("%5v| @{y}%-16v@{|}| %v\n", "", resMsg, msg)
					warning++
					color.Fprintf(w, result)
				}

			}
		}
		color.Fprintf(w, "\n@{g}PASS: %d@{|} / @{r}FAIL: %d@{|} / @{y}WARNING: %d@{|}\n", pass, fail, warning)
	}
}

func readFile(path string) (string, []Item) {
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
	overlayExtensions = f.WarningOverlay
	fmt.Printf("WarningOverlay: %v\n", strings.Join(overlayExtensions, ","))
	maxPathThreshold = f.WarningMaxPath
	fmt.Printf("WarningMaxPath: %d\n", maxPathThreshold)
	fmt.Printf("Message: %v\n", f.Message)
	if len(f.List) > 0 {
		return f.Title, f.List
	}
	return "", nil
}

func testFile(item Item) (bool, string) {
	switch item.Except {
	case "exist":
		s := formatTestResultMessage("exist")
		if !isExist(item.Fullpath) {
			return false, s + ": false"
		}
		return true, s + ": true "
	case "not_exist":
		s := formatTestResultMessage("not exist")
		if isExist(item.Fullpath) {
			return false, s + ": false"
		}
		return true, s + ": true "
	case "match":
		s := formatTestResultMessage("match")
		if !isMatch(item.Fullpath, item.Sha1) {
			return false, s + ": false"
		}
		return true, s + ": true "
	case "newer":
		s := formatTestResultMessage("newer")
		if !isNewer(item.Fullpath, item.Modtime) {
			return false, s + ": false"
		}
		return true, s + ": true "
	}

	return true, formatTestResultMessage("nomatch")
}

func formatTestResultMessage(msg string) string {
	return fmt.Sprintf("%9v", msg)
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

func isNewer(path string, except time.Time) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	if except.Unix() > info.ModTime().Unix() {
		//fmt.Printf("except: %v, given: %v\n", except, info.ModTime())
		return false
	}
	return true
}

func shouldOverlayTest(filename string) bool {
	for _, ext := range overlayExtensions {
		if ext == strings.TrimLeft(filepath.Ext(filename), ".") {
			return true
		}
	}
	return false
}

func testOverlay(baseDir string, item Item) []string {
	var s []string
	path := item.Fullpath
	except := filepath.Base(path[:len(path)-len(filepath.Ext(path))])
	ext := strings.TrimLeft(filepath.Ext(path), ".")
	err := filepath.Walk(filepath.Dir(baseDir), func(p string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if info.Name() != except+"."+ext {
				if isOverlay(info.Name(), ext, except) {
					s = append(s, p)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	if len(s) > 0 {
		return s
	} else {
		return nil
	}
}

func isOverlay(filename string, ext string, except string) bool {
	if strings.TrimLeft(filepath.Ext(filename), ".") == ext {
		if strings.Contains(filename, except) {
			return true
		}
	}
	return false
}

func testMaxPath(path string) (bool, int) {
	num := utf8.RuneCountInString(path)
	if num > maxPathThreshold {
		return false, num
	}
	return true, num
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
