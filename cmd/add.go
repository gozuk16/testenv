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
		//fmt.Println("add called")

		a, err := cmd.Flags().GetBool("all")
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Printf("-a = %v\n", a)

		if len(args) > 0 {
			//for i, _ := range args {
			//	fmt.Println(args[i])
			//}
			seachFile(args[0], a)
		}
	},
}

func seachFile(path string, optA bool) {
	var f = File{}
	origin, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}
	origin = filepath.Clean(origin)
	if err != nil {
		log.Fatal(err)
	}
	f.BaseDir = filepath.FromSlash(origin)
	f.List = make([]Item, 0)
	f.WarningOverlay = append(f.WarningOverlay, "jar")
	f.WarningOverlay = append(f.WarningOverlay, "dll")
	f.WarningMaxPath = 220
	i := 0
	err = filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		fp, err := filepath.Abs(p)
		if err != nil {
			log.Fatal(err)
		}
		fp = filepath.Clean(fp)
		if err != nil {
			log.Fatal(err)
		}

		if info.IsDir() {
			if !optA && isDot(fp) {
				//fmt.Printf("skip path: %v\n", fp)
				return filepath.SkipDir
			}
		} else {
			if !optA && isDot(fp) {
				//fmt.Printf("skip file: %v\n", fp)
			} else {
				f.List = append(f.List, Item{
					i + 1,                   // id
					info.Name(),             // filename
					filepath.FromSlash(fp),  // fullpath
					info.ModTime(),          // modtime
					info.Size(),             // size
					info.Mode().IsRegular(), // rw
					info.Mode().Perm(),      // mode
					info.Mode().String(),    // modestring
					getFileHash(p),          // sha1
					"match"})                // except
				i++
			}
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

func isDot(path string) bool {
	str := filepath.Base(path)
	for pos, c := range str {
		if pos == 0 && c != '.' {
			return false
		}
		//fmt.Printf("位置: %d 文字: %v\n", pos, string([]rune{c}))
		if pos == 1 && c != '.' {
			return true
		}
	}
	//fmt.Printf("最後まで来た：%v\n", str)
	return true
}

func init() {
	RootCmd.AddCommand(addCmd)

	addCmd.Flags().BoolP("all", "a", false, "Include directory entries whose names begin with a dot (.).")
}
