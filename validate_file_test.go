package main

import (
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
			result, length := maxPath(c.path)
			//fmt.Println(length)
			if result != c.except {
				t.Error(c.path)
			}
		})
	}
}

func setMaxPathData(cases *[]TestMaxPathData) {
	*cases = []TestMaxPathData{
		{path: "/123456789/123456789/123456789/123456789/123456789", except: true, msg: "50"},
		{path: "/123456789/123456789/123456789/123456789/123456789/", except: false, msg: "51"},
		{path: "C:\\123456789\\123456789\\123456789\\123456789\\1234567", except: true, msg: "50 Windows"},
		{path: "C:\\123456789\\123456789\\123456789\\123456789\\12345678", except: false, msg: "51 Windows"},
		{path: "/123456789/123456789/123456789/123456789/12345日本語𩸽", except: true, msg: "50 日本語サロゲートペア入り"},
		{path: "/123456789/123456789/123456789/123456789/12345日本語𩸽/", except: false, msg: "51 日本語サロゲートペア入り"},
	}
}
