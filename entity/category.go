package entity

//type Category struct {
//	ID          string
//	Name        string
//	description string
//}

type Category string

const (
	FootballCategory Category = "football"
	HistoryCategory  Category = "History"
)

func (c Category) IsValid() bool {
	switch c {
	case FootballCategory:
		return true
	}
	return false
}

func CategoryList() []Category {
	return []Category{FootballCategory, HistoryCategory}
}
