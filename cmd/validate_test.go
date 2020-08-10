package cmd

import (
	"fmt"
	"runtime"
	"testing"
)

type TestMaxPathData struct {
	path   string
	except bool
	msg    string
}

func TestMaxPathLength(t *testing.T) {

	cases := []TestMaxPathData{}
	setMaxPathData(&cases)

	//maxPathTreshold := 50
	for _, c := range cases {
		t.Run(c.msg, func(t *testing.T) {
			result, length := testMaxPath(c.path)
			fmt.Println(length)
			if result != c.except {
				t.Error(c.path)
			}
		})
	}
}

func setMaxPathData(cases *[]TestMaxPathData) {
	if runtime.GOOS == "windows" {
		*cases = []TestMaxPathData{
			{path: "C:\\var\\test\\testfile", except: false, msg: "フルパスの通常ファイル(Windows)"},
			{path: "C:\\var\\test\\.", except: true, msg: "フルパスのカレントディレクトリ(Windows)"},
			{path: "C:\\var\\test\\.test", except: true, msg: "フルパスのドットファイル(Windows)"},
		}
	} else {
		*cases = []TestMaxPathData{
			{path: ".", except: true, msg: "カレントディレクトリ(相対)"},
			{path: "..", except: true, msg: "上位ディレクトリ(相対)"},
			{path: "testfile", except: false, msg: "通常ファイル(相対)"},
			{path: "a.b", except: false, msg: "ドットが2文字目に来る(相対)"},
			{path: ".a", except: true, msg: "ドットファイル(相対)"},
			{path: ".ab", except: true, msg: "ドットファイル(相対)"},
			{path: "..a", except: true, msg: "ドット2個ファイル(相対)"},
			{path: "..ab", except: true, msg: "ドット2個ファイル(相対)"},
			{path: "...a", except: true, msg: "ドット3個ファイル(相対)"},
			{path: "...ab", except: true, msg: "ドット3個ファイル(相対)"},
			{path: "/var/test/testfile", except: false, msg: "フルパスの通常ファイル(UNIX)"},
			{path: "/var/test/.", except: true, msg: "フルパスのカレントディレクトリ(UNIX)"},
			{path: "/var/test/.test", except: true, msg: "フルパスのドットファイル(UNIX)"},
			//{path: "C:\\var\\test\\testfile", except: false, msg: "フルパスの通常ファイル(Windows)"},
			//{path: "C:\\var\\test\\.", except: true, msg: "フルパスのカレントディレクトリ(Windows)"},
			//{path: "C:\\var\\test\\.test", except: true, msg: "フルパスのドットファイル(Windows)"},
			{path: "日本語𩸽", except: false, msg: "日本語ファイル(サロゲートペア入り)"},
		}
	}
}
