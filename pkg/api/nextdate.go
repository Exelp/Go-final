package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DATE_Layout = "20060102"

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	var now time.Time
	var err error

	nowStart := r.FormValue("now")
	dStart := r.FormValue("date")
	repeat := r.FormValue("repeat")

	if nowStart == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DATE_Layout, nowStart)
		if err != nil {
			http.Error(w, "invalid date format", http.StatusBadRequest)
			return
		}
	}
	if dStart == "" || repeat == "" {
		http.Error(w, "invalid date format", http.StatusBadRequest)
		return
	}

	nextDate, err := NextDate(now, dStart, repeat)
	if err != nil {
		http.Error(w, "invalid date format", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))

}
func NextDate(now time.Time, dStart string, repeat string) (string, error) {
	dateStart, err := time.Parse(DATE_Layout, dStart)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %w", err)
	}
	parts := strings.Split(repeat, " ")
	rule := parts[0]

	switch rule {
	case "d":
		if len(parts) < 2 {
			return "", fmt.Errorf("invalid rule 'd': %w", err)
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("invalid rule 'd': %w", err)
		}
		if days > 400 || days <= 0 {
			return "", fmt.Errorf("invalid number of days rules 'd': %d", days)
		}
		for {
			dateStart = dateStart.AddDate(0, 0, days)
			if dateStart.After(now) {
				return dateStart.Format(DATE_Layout), nil
			}
		}
	case "y":
		for {
			dateStart = dateStart.AddDate(1, 0, 0)
			if dateStart.After(now) {
				return dateStart.Format(DATE_Layout), nil
			}
		}
	case "w":
		if len(parts) < 2 {
			return "", fmt.Errorf("invalid rule 'w': %w", err)
		}
		var weekDays [8]bool

		weekDayStr := strings.Split(parts[1], ",")
		for _, wd := range weekDayStr {
			day, err := strconv.Atoi(wd)
			if err != nil || day < 1 || day > 7 {
				return "", fmt.Errorf("invalid rule 'w': %s", wd)
			}
			weekDays[day] = true
		}
		for {
			dateStart = dateStart.AddDate(0, 0, 1)
			if dateStart.After(now) {
				wd := int(dateStart.Weekday())
				if wd == 0 {
					wd = 7
				}
				if weekDays[wd] {
					return dateStart.Format(DATE_Layout), nil
				}
			}
		}

	case "m":
		if len(parts) < 2 {
			return "", fmt.Errorf("invalid rule 'm': %w", err)
		}
		var dayMonth [32]bool
		var monthYears [13]bool

		daysMonthStr := strings.Split(parts[1], ",")
		last := false
		preLast := false

		for _, dm := range daysMonthStr {
			d, err := strconv.Atoi(dm)
			if err != nil {
				return "", fmt.Errorf("invalid rule 'm': %s", dm)
			}

			switch {
			case d >= 1 && d <= 31:
				dayMonth[d] = true
			case d == -1:
				last = true
			case d == -2:
				preLast = true
			default:
				return "", fmt.Errorf("invalid rule 'm': %d", d)
			}

		}

		if len(parts) > 2 {
			monthStrs := strings.Split(parts[2], ",")
			for _, monthStr := range monthStrs {
				m, err := strconv.Atoi(monthStr)
				if err != nil || m < 1 || m > 12 {
					return "", fmt.Errorf("invalid month value: %s", monthStr)
				}
				monthYears[m] = true
			}
		} else {
			for i := 1; i <= 12; i++ {
				monthYears[i] = true
			}
		}

		for {
			dateStart = dateStart.AddDate(0, 0, 1)
			if !dateStart.After(now) {
				continue
			}
			day := dateStart.Day()
			month := dateStart.Month()
			year := dateStart.Year()
			lastDayOfMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
			preLastDay := lastDayOfMonth - 1

			if (dayMonth[day]) ||
				(last && day == lastDayOfMonth) ||
				(preLast && day == preLastDay) {
				if monthYears[int(month)] {

					return dateStart.Format(DATE_Layout), nil

				}

			}

		}

	}
	return "", fmt.Errorf("invalid date format")
}
