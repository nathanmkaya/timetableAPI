package parser

import (
	"fmt"

	"regexp"

	"strings"

	"github.com/tealeg/xlsx"
)

func getDates(sheet *xlsx.Sheet) {
	fmt.Println(sheet.Name)
	for row, _ := range sheet.Rows {
		for col, _ := range sheet.Cols {
			if col != 1 {
				cell := sheet.Cell(row, col)
				pattern := "(Mo(n(day)?)?|Tu(e(sday)?)?|We(d(nesday)?)?|Th(u(r(s(day)?)?)?)?|Fr(i(day)?)?|Sa(t(urday)?)?|Su(n(day)?)?) (3[01]|[12][0-9]|0[1-9])/(1[0-2]|0[1-9])/[0-9]{2}"
				r, _ := regexp.Compile(pattern)
				date := strings.Title(strings.ToLower(cell.Value))
				if r.MatchString(date) {
					fmt.Println(cell)
				}
			}
		}
	}
}

func getShift(sheet xlsx.Sheet) string {
	return ""
}

func getDate(sheet *xlsx.Sheet, row int, col int) string {
	pattern := "(Mo(n(day)?)?|Tu(e(sday)?)?|We(d(nesday)?)?|Th(u(r(s(day)?)?)?)?|Fr(i(day)?)?|Sa(t(urday)?)?|Su(n(day)?)?) (3[01]|[12][0-9]|0[1-9])/(1[0-2]|0[1-9])/[0-9]{2}"
	r, _ := regexp.Compile(pattern)
	for row := row; row >= 0; row-- {
		for col := col; col >= 0; col-- {
			date := strings.Title(strings.ToLower(sheet.Cell(row, col).Value))
			if r.MatchString(date) {
				return date
			}
		}
	}
	return ""
}

func getColSpan(sheet xlsx.Sheet) int {
	var count int = 0
	for row, _ := range sheet.Rows {
		if sheet.Cell(row, 1).Value != "" {
			count++
		}
	}
	return count
}

func getTime(sheet *xlsx.Sheet, row int, col int) string {
	pattern := "([1-9]|1[0-2]).([0-5][0-9]|60)(PM|AM)-([1-9]|1[0-2]).([0-5][0-9]|60)(PM|AM)"
	r, _ := regexp.Compile(pattern)
	for row := row; row >= 0; row-- {
		for col := col; col >= 0; col-- {
			date := strings.Title(strings.ToLower(sheet.Cell(row, col).Value))
			if r.MatchString(date) {
				return date
			}
		}
	}
	return ""
}

func ParseExams(filename string) {
	xlFile, err := xlsx.OpenFile(filename)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	for _, v := range xlFile.Sheets {
		fmt.Println(v.Name)
		for row, _ := range v.Rows {
			for col, _ := range v.Cols {
				if col == 0 || v.Cell(row, col).Value == "" || v.Cell(row, 1).Value == "ROOM" || v.Cell(row, col).Value == "CHAPEL" {
					continue
				} else {
					fmt.Printf("%s : %s\n", v.Cell(row, col).Value, getDate(v, row, col))
				}
			}
		}
	}
}
