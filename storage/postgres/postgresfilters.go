package postgres

import (
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/ucladevx/BConnect-backend/models"
)

// Filters filters
type Filters struct {
	client *gorm.DB
}

// NewFilterers filterer
func NewFilterers(client *gorm.DB) *Filters {
	return &Filters{
		client: client,
	}
}

// AllFilter filters
func (f Filters) AllFilter(curr *models.FilterReturn, emptyArgs []string) *models.FilterReturn {
	return curr
}

// NameFilter filters
func (f Filters) NameFilter(curr *models.FilterReturn, names []string) *models.FilterReturn {
	for i := 0; i < len(names); i++ {
		curr.Filter = curr.Filter.Where("FIRST_NAME = ?", names[i])
	}
	return curr
}

// MajorFilter filters
func (f Filters) MajorFilter(curr *models.FilterReturn, majors []string) *models.FilterReturn {
	for i := 0; i < len(majors); i++ {
		curr.Filter = curr.Filter.Where("MAJOR = ?", majors[i])
	}
	return curr
}

// GradYearFilter filters
func (f Filters) GradYearFilter(curr *models.FilterReturn, gradYear []string) *models.FilterReturn {
	for i := 0; i < len(gradYear); i++ {
		curr.Filter = curr.Filter.Where("GRAD_YEAR = ?", gradYear[i])
	}
	return curr
}

// InterestsFilter filters
func (f Filters) InterestsFilter(curr *models.FilterReturn, interests []string) *models.FilterReturn {
	for i := 0; i < len(interests); i++ {
		curr.Filter = curr.Filter.Where("INTERESTS = ?", interests[i])
	}
	return curr
}

// LocationRadiusFilter filters
func (f Filters) LocationRadiusFilter(curr *models.FilterReturn, radius []string) *models.FilterReturn {
	for i := 0; i < len(radius); i++ {
		read, _ := strconv.ParseFloat(radius[i], 64)
		curr.Filter = curr.Filter.Where("LAT <= ?", read)
	}
	return curr
}

// FinalFilter final filter
func (f Filters) FinalFilter(filters map[string]models.Filterer, args map[string][]string) map[string]interface{} {
	var user []models.User

	curr := &models.FilterReturn{
		Filter: f.client,
	}
	if _, ok := args["name"]; ok {
		curr = f.NameFilter(curr, args["name"])
	}
	if _, ok := args["year"]; ok {
		curr = f.GradYearFilter(curr, args["year"])
	}
	if _, ok := args["major"]; ok {
		curr = f.MajorFilter(curr, args["major"])
	}
	if _, ok := args["interest"]; ok {
		curr = f.InterestsFilter(curr, args["interest"])
	}
	if _, ok := args["radius"]; ok {
		curr = f.LocationRadiusFilter(curr, args["radius"])
	}
	if err := curr.Filter.Find(&user); err != nil {
	}
	var resp = map[string]interface{}{"matches": user}
	return resp
}
