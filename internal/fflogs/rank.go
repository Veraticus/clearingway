package fflogs

type Rank struct {
	Percent float64 `json:"rankPercent"`
}

var Ranks = []string{"NA'S Comfiest", "Gray", "Green", "Blue", "Purple", "Orange", "Pink", "Gold"}

func (r *Rank) Color() string {
	if r.Percent == 100 {
		return "Gold"
	}
	if r.Percent >= 99 {
		return "Pink"
	}
	if r.Percent >= 95 {
		return "Orange"
	}
	if r.Percent >= 75 {
		return "Purple"
	}
	if r.Percent >= 50 {
		return "Blue"
	}
	if r.Percent >= 25 {
		return "Green"
	}
	if r.Percent > 0 {
		return "Gray"
	}

	return "NA'S Comfiest"
}
