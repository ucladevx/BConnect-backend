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
func (f Filters) AllFilter(curr *models.FilterReturn, currentUser *models.User, emptyArgs []string) *models.FilterReturn {
	return curr
}

// NameFilter filters
func (f Filters) NameFilter(curr *models.FilterReturn, currentUser *models.User, names []string) *models.FilterReturn {
	for i := 0; i < len(names); i++ {
		curr.Filter = curr.Filter.Where("FIRST_NAME = ? AND USER_ID != ?", names[i], currentUser.UserID)
	}
	return curr
}

// MajorFilter filters
func (f Filters) MajorFilter(curr *models.FilterReturn, currentUser *models.User, majors []string) *models.FilterReturn {
	for i := 0; i < len(majors); i++ {
		curr.Filter = curr.Filter.Where("MAJOR = ? AND USER_ID != ?", majors[i], currentUser.UserID)
	}
	return curr
}

// GradYearFilter filters
func (f Filters) GradYearFilter(curr *models.FilterReturn, currentUser *models.User, gradYear []string) *models.FilterReturn {
	for i := 0; i < len(gradYear); i++ {
		curr.Filter = curr.Filter.Where("GRAD_YEAR = ? AND USER_ID != ?", gradYear[i], currentUser.UserID)
	}
	return curr
}

// InterestsFilter filters
func (f Filters) InterestsFilter(curr *models.FilterReturn, currentUser *models.User, interests []string) *models.FilterReturn {
	for i := 0; i < len(interests); i++ {
		curr.Filter = curr.Filter.Where("INTERESTS = ? AND USER_ID != ?", interests[i], currentUser.UserID)
	}
	return curr
}

// LocationRadiusFilter filters
func (f Filters) LocationRadiusFilter(curr *models.FilterReturn, currentUser *models.User, radii []string) *models.FilterReturn {
	for i := 0; i < len(radii); i++ {
		radius, _ := strconv.ParseFloat(radii[i], 64)
		curr.Filter = curr.Filter.Where("(LAT - ?)^2 + (LON - ?)^2 <= ? AND USER_ID != ?", strconv.FormatFloat(currentUser.Lat, 'f', -1, 64), strconv.FormatFloat(currentUser.Lon, 'f', -1, 64), radius*radius, currentUser.UserID)
	}
	return curr
}

// FinalFilter final filter
func (f Filters) FinalFilter(filters map[string]models.Filterer, currentUser *models.User, args map[string][]string) []models.User {
	var user []models.User

	curr := &models.FilterReturn{
		Filter: f.client,
	}
	if _, ok := args["name"]; ok {
		curr = f.NameFilter(curr, currentUser, args["name"])
	}
	if _, ok := args["year"]; ok {
		curr = f.GradYearFilter(curr, currentUser, args["year"])
	}
	if _, ok := args["major"]; ok {
		curr = f.MajorFilter(curr, currentUser, args["major"])
	}
	if _, ok := args["interest"]; ok {
		curr = f.InterestsFilter(curr, currentUser, args["interest"])
	}
	if _, ok := args["radius"]; ok {
		curr = f.LocationRadiusFilter(curr, currentUser, args["radius"])
	}
	if err := curr.Filter.Find(&user); err != nil {
	}
	return user
}
