package domain

var Names = map[string]string{
	"math":            "алгебра",
	"geometry":        "геометрия",
	"russian":         "русский язык",
	"literature":      "литература",
	"history":         "история",
	"physics":         "физика",
	"chemistry":       "химия",
	"biology":         "биология",
	"geography":       "география",
	"english":         "английский",
	"social":          "общестознание",
	"random-and-stat": "вероятность и статистика",
	"projects":        "проектная деят.",
}

func DisplayName(code string) string {
	if name, ok := Names[code]; ok {
		return name
	}
	return code
}
