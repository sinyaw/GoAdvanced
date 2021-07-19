package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func getMonthsInTheYear(Year int) ([]string, error) {
	var err = errors.New("Invalid Input")
	Months := []string{}
	for i := 1; i <= 12; i++ {
		if time.Now().Month() <= time.Month(i) && time.Now().Year() == Year {
			Months = append(Months, time.Month(i).String())
		} else if time.Now().Year() < Year {
			Months = append(Months, time.Month(i).String())
		} else {
			return nil, err
		}
	}
	return Months, err
}

func getDatesInTheMonth(Month int, Year int) [][]int {
	Dates := [][]int{}
	for i := 1; i <= 31; i++ {
		t := time.Date(Year, time.Month(Month), i, 0, 0, 0, 0, time.UTC)
		if time.Now().Month() == time.Month(Month) && time.Now().Year() == Year {
			if time.Now().Day() <= t.Day() {
				TD := []int{int(t.Weekday()), t.Day()}
				Dates = append(Dates, TD)
			} else {
				TD := []int{int(t.Weekday()), 0}
				Dates = append(Dates, TD)
			}
		} else if t.Month() == time.Month(Month) && t.Year() == Year {
			TD := []int{int(t.Weekday()), t.Day()}
			Dates = append(Dates, TD)
		} else {
			break
		}
	}
	return Dates
}

func printCalender(Month int, Year int) {
	weekday := []string{"S", "M", "T", "W", "T", "F", "S"}
	Dates := getDatesInTheMonth(Month, Year)
	fmt.Println()
	//Print Year and Month
	fmt.Println(Year, time.Month(Month))
	for _, v := range weekday {
		fmt.Printf(" %s\t", v)
	}
	fmt.Println()
	for i := 0; i < Dates[0][0]; i++ {
		fmt.Printf("\t")
	}
	for _, v := range Dates {
		if v[1] != 1 && v[0] == 0 {
			fmt.Println()
		}
		if v[1] == 0 {
			fmt.Printf("\t")
		} else {
			if v[1] < 10 {
				fmt.Printf(" %d\t", v[1])
			} else {
				fmt.Printf("%d\t", v[1])
			}

		}
	}
	fmt.Println()
	fmt.Println()
}

func ReturnADate(Day int, Month int, Year int) time.Time {
	d := time.Date(Year, time.Month(Month), Day, 0, 0, 0, 0, time.UTC)
	return d
}

func timeToString(tm time.Time) string {
	if tm.IsZero() {
		return ""
	} else {
		D := tm
		y, n, d := D.Date()
		m := int(n)
		return strconv.Itoa(y) + "#" + strconv.Itoa(m) + "#" + strconv.Itoa(d)
	}
}

func stringToTime(str string) time.Time {
	if str == "" {
		return time.Time{}
	} else {
		sstr := strings.Split(str, "#")
		y, _ := strconv.Atoi(sstr[0])
		m, _ := strconv.Atoi(sstr[1])
		d, _ := strconv.Atoi(sstr[2])
		return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
	}
}

func checkDateInput(msg string, now time.Time, after time.Time) (int, int, int) {
	var d, m, y int
loop:
	for {
		fmt.Printf(msg)
		input := userInput()
		if input == "*" {
			page = 0
			pages(page)
			return 0, 0, 0
		}

		vSlice := strings.Split(input, "/") // **move to calenders.go

		if len(vSlice) == 3 {
			dd, errdd := checkInt(vSlice[0])
			mm, errmm := checkInt(vSlice[1])
			yyyy, erryyyy := checkInt(vSlice[2])
			if errdd != nil {
				fmt.Println(errdd, "'dd'")
				fmt.Println()

			}
			if errmm != nil {
				fmt.Println(errmm, "'mm'")
				fmt.Println()
			}
			if erryyyy != nil {
				fmt.Println(erryyyy, "'yyyy'")
				fmt.Println()
			}
			if errdd != nil || errmm != nil || erryyyy != nil {
				continue loop
			} else if now.Year() > yyyy || after.Year() < yyyy {
				fmt.Println("Year is out of the range!", "'yyyy'")
				fmt.Println()
				continue loop
			} else if mm <= 0 || mm > 12 {
				fmt.Println("Month is out of the range!", "'mm'")
				fmt.Println()
				continue loop
			} else {
				if dd <= 0 {
					fmt.Println("Invalid input!", "'dd'")
					fmt.Println()
				} else {
					DMrange := getDatesInTheMonth(mm, yyyy)
				loop2:
					for {
						for _, v := range DMrange {
							if dd == v[1] {
								break loop2
							}
						}
						fmt.Println("Date is out of the month range!", "'dd'")
						fmt.Println()
						continue loop
					}
				}
				d = dd
				m = mm
				y = yyyy
			}
		} else {
			fmt.Println("Invalid input!") //error msg
			fmt.Println()
			continue loop
		}
		break loop
	}
	return d, m, y
}
