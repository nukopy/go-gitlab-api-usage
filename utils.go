package main

import "time"

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

var layout = "2006-01-02 15:04:05"

func toJST(t time.Time) time.Time {
	tJST := t.In(time.FixedZone("Asia/Tokyo", 9*60*60))
	return tJST
}

func TimeToString(t time.Time) string {
	str := t.Format(layout)
	return str
}

func TimeToJSTString(t time.Time) string {
	tJST := toJST(t)
	str := tJST.Format(layout)
	return str
}
