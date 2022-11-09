package chart

const MINUTE = 60
const HOUR = MINUTE * 60
const DAY = HOUR * 24
const WEEK = DAY * 7
const MONTH = DAY * 30

type HighLow struct {
	Time  int
	Low   float64
	High  float64
	First float64
	Last  float64
}
