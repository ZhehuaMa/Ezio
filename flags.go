package Ezio

import "flag"

var qps int64

func init() {
	flag.Int64Var(&qps, "qps", 1, "Quests per second.")
	flag.Parse()
}
