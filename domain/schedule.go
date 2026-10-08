package domain

var Weekly = map[int][]string{
	0: {
		"russian",
		"english",
		"english",
		"random-and-stat",
		"history",
	},
	1: {
		"physics",
		"physics",
		"russian",
		"history",
		"chemistry",
	},
	2: {
		"russian",
		"literature",
		"social",
		"social",
		"history",
		"history",
		"projects",
	},
	3: {
		"english",
		"geometry",
		"geometry",
		"history",
		"biology",
		"history",
	},
	4: {
		"math",
		"math",
		"literature",
		"literature",
		"geography",
		"social",
		"social",
	},
	5: {
		"math",
		"math",
		"russian",
		"projects",
	},
}

var DayNames = []string{"Пн", "Вт", "Ср", "Чт", "Пт", "Сб"}

func LessonsForDay(day int) []string {
	return Weekly[day]
}
