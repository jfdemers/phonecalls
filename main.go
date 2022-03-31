package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var path = flag.String("file", "", "path of phone calls logs")

func parseDate(d string) (time.Time, error) {
	var date = time.Now()

	parser := regexp.MustCompile(`(\d{4})-(\d{2})-(\d{2}) (\d{1,2}):(\d{2}):(\d{2}) (a\.m\.|p\.m\.)`)
	result := parser.FindStringSubmatch(d)

	if len(result) <= 0 {
		return date, errors.New("invalid date")
	}

	year, err := strconv.Atoi(result[1])
	if err != nil {
		return date, err
	}

	month, err := strconv.Atoi(result[2])
	if err != nil {
		return date, err
	}

	m := time.Month(month)

	day, err := strconv.Atoi(result[3])
	if err != nil {
		return date, err
	}

	hour, err := strconv.Atoi(result[4])
	if err != nil {
		return date, err
	}

	minute, err := strconv.Atoi(result[5])
	if err != nil {
		return date, err
	}

	second, err := strconv.Atoi(result[6])
	if err != nil {
		return date, err
	}

	if result[7] == "p.m." && hour != 12 {
		hour += 12
	}

	return time.Date(year, m, day, hour, minute, second, 0, time.Local), nil
}

func displayTitle() {
	fmt.Printf("%-20s %10s %10s %10s %10s %10s %10s %10s \n", "Date", "Avant 9h", "9-12", "12-13", "13-17", "17-20", "20+", "Total")
}

func displayStats(currentDay string, total int, before9 int, day9to12 int, day12to13 int, day13to17 int, day17to20 int, after20 int) {
	//log.Println(currentDay, ":", before9, day9to12, day12to13, day13to17, day17to20, after20, total)
	fmt.Printf("%-20s %10d %10d %10d %10d %10d %10d %10d \n", currentDay, before9, day9to12, day12to13, day13to17, day17to20, after20, total)
}

func main() {
	flag.Parse()

	file, err := os.Open(*path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	lineParser := regexp.MustCompile(`(.*),(.*),(.*),(.*),(.*),(.*),(.*),(.*),(.*),(.*),(.*)`)
	scanner := bufio.NewScanner(file)

	total := 0
	currentDay := ""
	dayTotal := 0
	dayBefore9 := 0
	day9to12 := 0
	day12to13 := 0
	day13to17 := 0
	day17to20 := 0
	dayAfter20 := 0

	displayTitle()

	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		t := scanner.Text()
		result := lineParser.FindStringSubmatch(t)

		if len(result) <= 0 {
			continue
		}

		if len(result[1]) > 0 && strings.Contains(result[4], "RG General") {
			d, err := parseDate(result[1])
			if err != nil {
				continue
			}

			day := d.Weekday().String() + " " + strconv.Itoa(d.Year()) + "-" + strconv.Itoa(int(d.Month())) + "-" + strconv.Itoa(d.Day())
			if currentDay != day {
				if total > 0 {
					displayStats(currentDay, dayTotal, dayBefore9, day9to12, day12to13, day13to17, day17to20, dayAfter20)
				}

				currentDay = day
				dayTotal = 0
				dayBefore9 = 0
				day9to12 = 0
				day12to13 = 0
				day13to17 = 0
				day17to20 = 0
				dayAfter20 = 0
			}

			if d.Hour() < 9 {
				dayBefore9++
			} else if d.Hour() < 12 {
				day9to12++
			} else if d.Hour() < 13 {
				day12to13++
			} else if d.Hour() < 17 {
				day13to17++
			} else if d.Hour() < 20 {
				day17to20++
			} else if d.Hour() >= 20 {
				dayAfter20++
			}

			dayTotal++
			total++
		}
	}

	fmt.Println("Appels total: ", total)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
