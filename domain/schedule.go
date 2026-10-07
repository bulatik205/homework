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
		"informatics",
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
		"informatics",
		"english",
		"geometry",
		"geometry",
		"pe",
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
		"pe",
	},
	5: {
		"math",
		"math",
		"english",
		"projects",
	},
}

var DayNames = []string{"Пн", "Вт", "Ср", "Чт", "Пт", "Сб"}

func LessonsForDay(day int) []string {
	return Weekly[day]
}
