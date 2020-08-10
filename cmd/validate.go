package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/spf13/cobra"

	"github.com/shiena/ansicolor"
	"github.com/wsxiaoys/terminal/color"
)

// validateCmd represents the add command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "validate config.json",
	Long:  `execute validate`,
	Run: func(cmd *cobra.Command, args []string) {
		optE, err := cmd.Flags().GetBool("expand")
		if err != nil {
			fmt.Println(err)
		}
		optV, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			fmt.Println(err)
		}

		if len(args) > 0 {
			//for i, _ := range args {
			//	fmt.Println(args[i])
			//}
			validateFiles(args[0], optE, optV)
		}
	},
}

var overlayExtensions []string
var maxPathThreshold int

func validateFiles(path string, optE bool, optV bool) {
	baseDir, list := readFile(path)
	if len(list) > 0 {
		startTime := time.Now()
		var pass, fail, warning int
		var id int
		var fullpath string
		var hasOverlay bool
		targetFiles := getOverlayTargetFiles(baseDir)
		w := ansicolor.NewAnsiColorWriter(os.Stdout)
		for _, l := range list {
			hasOverlay = false
			validateResult, msg := validateFile(l)
			if validateResult {
				if optV {
					color.Fprintf(w, "%5d| @{g}%-16v@{|}| %v\n", l.ID, msg, l.Fullpath)
				} else {
					id = l.ID
					fullpath = l.Fullpath
				}
				pass++
			} else {
				color.Fprintf(w, "%5d| @{r}%-16v@{|}| %v\n", l.ID, msg, l.Fullpath)
				fail++
			}

			resMsg := formatResultMessage("overlay") + ": warn"
			if shouldOverlay(l.Filename) {
				msgs := overlay(baseDir, targetFiles, l, optE)
				if msgs != nil {
					if !optV && validateResult {
						color.Fprintf(w, "%5d| @{g}%-16v@{|}| %v\n", id, msg, fullpath)
						hasOverlay = true
					}
					for i, m := range msgs {
						if i == len(msgs)-1 {
							m = " └─ " + m
						} else {
							m = " ├─ " + m
						}
						color.Fprintf(w, "%5v| @{y}%-16v@{|}| %v\n", "", resMsg, m)
					}
					warning++
				}
			}

			resMsg = formatResultMessage("max path") + ": warn"
			if maxPathThreshold > 0 {
				r, size := maxPath(l.Fullpath)
				if !r {
					if !optV && validateResult && !hasOverlay {
						color.Fprintf(w, "%5d| @{g}%-16v@{|}| %v\n", id, msg, fullpath)
					}
					m := fmt.Sprintf(" └─ path len:%-3d, over: %-3d", size, size-maxPathThreshold)
					color.Fprintf(w, "%5v| @{y}%-16v@{|}| %v\n", "", resMsg, m)
					warning++
				}

			}
		}
		color.Fprintf(w, "\n@{g}PASS: %d@{|} / @{r}FAIL: %d@{|} / @{y}WARNING: %d@{|}\n", pass, fail, warning)
		endTime := time.Now()
		fmt.Printf("elapsed time: %fs\n", (endTime.Sub(startTime)).Seconds())
	}
}

func readFile(path string) (string, []Item) {
	p, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	p = filepath.Clean(p)
	raw, err := ioutil.ReadFile(p)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	var f = File{}
	if err := json.Unmarshal(raw, &f); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nBaseDir: %v\n", f.BaseDir)
	fmt.Printf("Num: %d\n", f.Num)
	overlayExtensions = f.WarningOverlay
	fmt.Printf("WarningOverlay: %v\n", strings.Join(overlayExtensions, ","))
	maxPathThreshold = f.WarningMaxPath
	fmt.Printf("WarningMaxPath: %d\n", maxPathThreshold)
	fmt.Printf("Message: %v\n", f.Message)
	if len(f.List) > 0 {
		return f.BaseDir, f.List
	}
	return "", nil
}

func validateFile(item Item) (bool, string) {
	switch item.Except {
	case "exist":
		s := formatResultMessage("exist")
		if !isExist(item.Fullpath) {
			return false, s + ": false"
		}
		return true, s + ": true "
	case "not_exist":
		s := formatResultMessage("not exist")
		if isExist(item.Fullpath) {
			return false, s + ": false"
		}
		return true, s + ": true "
	case "match":
		s := formatResultMessage("match")
		if !isMatch(item.Fullpath, item.Sha1) {
			return false, s + ": false"
		}
		return true, s + ": true "
	case "newer":
		s := formatResultMessage("newer")
		if !isNewer(item.Fullpath, item.Modtime) {
			return false, s + ": false"
		}
		return true, s + ": true "
	}

	return true, formatResultMessage("nomatch")
}

func formatResultMessage(msg string) string {
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

func shouldOverlay(filename string) bool {
	for _, ext := range overlayExtensions {
		if ext == strings.TrimLeft(filepath.Ext(filename), ".") {
			return true
		}
	}
	return false
}

func getOverlayTargetFiles(baseDir string) []string {
	var s []string
	err := filepath.Walk(filepath.Clean(baseDir), func(p string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			s = append(s, p)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	if len(s) > 0 {
		return s
	}
	return nil
}

func overlay(baseDir string, targetFiles []string, item Item, optE bool) []string {
	var s []string
	path := filepath.Dir(item.Fullpath)
	filename := item.Filename
	except := filepath.Base(filename[:len(filename)-len(filepath.Ext(filename))])
	ext := strings.TrimLeft(filepath.Ext(filename), ".")

	for _, file := range targetFiles {
		f, _ := filepath.Abs(file)
		if !optE {
			if filepath.Dir(f) != path {
				continue
			}
		}
		if item.Fullpath != f {
			if isOverlay(file, ext, except) {
				if strings.HasPrefix(file, baseDir) {
					file = strings.TrimPrefix(file, baseDir+filepath.FromSlash("/"))
				}
				s = append(s, file)
			}
		}
	}

	if len(s) > 0 {
		return s
	}
	return nil
}

func isOverlay(filename string, ext string, except string) bool {
	separator := [...]string{"-", "_", "."}
	shortFilename := filepath.Base(filename)
	if strings.HasSuffix(shortFilename, "."+ext) {
		target := strings.TrimRight(shortFilename, "."+ext)

		// 先頭が"_"ならチェック対象に含める
		if target[:1] == "_" {
			target = strings.TrimLeft(target, "_")
		}
		// 先頭がチェック対象のファイル名と一致ならチェック対象に含める
		if strings.HasPrefix(target, except) {
			s := strings.TrimLeft(target, except)
			// 完全一致ならtrue、後ろになにかついてるなら引き続きチェック
			if s != "" {
				//c1 := string([]rune(s)[:1])
				//c2 := string([]rune(s)[1:2])
				// 対象はAscii領域でUTF-8の2Byte目にはこの領域が来ることはないので高速なByteでチェックで良い
				c1 := s[:1]
				var c2 string = ""
				if len(s) > 1 {
					c2 = s[1:2]
				}
				for _, sp := range separator {
					// 1文字目がセパレーターなら引き続きチェック
					if c1 == sp {
						// 2文字目が数字ならtrue
						_, err := strconv.Atoi(c2)
						if err != nil {
							return false
						}
						return true
					}
				}
				return false
			}
			return true
		}
	}
	return false
}

func maxPath(path string) (bool, int) {
	num := utf8.RuneCountInString(path)
	if num > maxPathThreshold {
		return false, num
	}
	return true, num
}

func init() {
	RootCmd.AddCommand(validateCmd)

	validateCmd.Flags().BoolP("verbose", "v", false, "Be verbose when validating, showing them as they are validated.")
	validateCmd.Flags().BoolP("expand", "e", false, "Be expand when checking for duplicates, include other directories")
}
