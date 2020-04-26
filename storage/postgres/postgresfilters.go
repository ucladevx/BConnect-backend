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
func (f Filters) AllFilter(curr *gorm.DB, emptyArgs []string) *gorm.DB {
	return curr
}

// NameFilter filters
func (f Filters) NameFilter(curr *gorm.DB, names []string) *gorm.DB {
	for i := 0; i < len(names); i++ {
		curr = curr.Where("FIRST_NAME = ?", names[i])
	}
	return curr
}

// MajorFilter filters
func (f Filters) MajorFilter(curr *gorm.DB, majors []string) *gorm.DB {
	for i := 0; i < len(majors); i++ {
		curr = curr.Where("MAJOR = ?", majors[i])
	}
	return curr
}

// GradYearFilter filters
func (f Filters) GradYearFilter(curr *gorm.DB, gradYear []string) *gorm.DB {
	for i := 0; i < len(gradYear); i++ {
		curr = curr.Where("GRAD_YEAR = ?", gradYear[i])
	}
	return curr
}

// InterestsFilter filters
func (f Filters) InterestsFilter(curr *gorm.DB, interests []string) *gorm.DB {
	for i := 0; i < len(interests); i++ {
		curr = curr.Where("INTERESTS = ?", interests[i])
	}
	return curr
}

// LocationRadiusFilter filters
func (f Filters) LocationRadiusFilter(curr *gorm.DB, radius []string) *gorm.DB {
	for i := 0; i < len(radius); i++ {
		read, _ := strconv.ParseFloat(radius[i], 64)
		curr = curr.Where("LAT <= ?", read)
	}
	return curr
}

// FinalFilter final filter
func (f Filters) FinalFilter(filters map[string]models.Filterer, args map[string][]string) map[string]interface{} {
	var user []models.User
	curr := f.client
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
	if err := curr.Find(&user); err != nil {
	}
	var resp = map[string]interface{}{"matches": user}
	return resp
}
