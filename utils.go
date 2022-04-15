package main

import "time"

var layoutDefault = "2006-01-02 15:04:05"
var layoutForFilename = "20060102-150405"

func toJST(t time.Time) time.Time {
	tJST := t.In(time.FixedZone("Asia/Tokyo", 9*60*60))
	return tJST
}

func timeToString(t time.Time, layout string) string {
	str := t.Format(layout)
	return str
}

func timeToJSTString(t time.Time, layout string) string {
	tJST := toJST(t)
	str := tJST.Format(layout)
	return str
}
