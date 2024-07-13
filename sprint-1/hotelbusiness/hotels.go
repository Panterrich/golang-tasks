//go:build !solution

package hotelbusiness

import (
	"cmp"
	"slices"
)

type Guest struct {
	CheckInDate  int
	CheckOutDate int
}

type Load struct {
	StartDate  int
	GuestCount int
}

func ComputeLoad(guests []Guest) []Load {
	if len(guests) == 0 {
		return []Load{}
	}

	lastDate := slices.MaxFunc(guests, func(a, b Guest) int {
		return cmp.Compare(a.CheckOutDate, b.CheckOutDate)
	}).CheckOutDate

	days := make([]Load, lastDate+1)

	for _, g := range guests {
		for i := g.CheckInDate; i < g.CheckOutDate; i++ {
			days[i].GuestCount++
		}
	}

	var load []Load
	lastLoad := 0

	for i, day := range days {
		if lastLoad != day.GuestCount {
			load = append(load, Load{StartDate: i, GuestCount: day.GuestCount})
			lastLoad = day.GuestCount
		}
	}

	return load
}
