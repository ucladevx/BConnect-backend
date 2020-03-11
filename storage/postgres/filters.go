package postgres

import (
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/ucladevx/BConnect-backend/models"
)

// Filters filters
type Filters struct {
	client *gorm.DB
	friend *gorm.DB
}

// NewFilterers filterer
func NewFilterers(client *gorm.DB, friendClient *gorm.DB) *Filters {
	return &Filters{
		client: client,
		friend: friendClient,
	}
}

// AllFilter filters
func (f Filters) AllFilter(emptyArgs []string) map[string]interface{} {
	var user []models.User
	if err := f.client.Find(&user); err != nil {

	}

	var resp = map[string]interface{}{"matches": user}
	return resp
}

// NameFilter filters
func (f Filters) NameFilter(names []string) map[string]interface{} {
	var user []models.User
	var curr *gorm.DB
	for i := 0; i < len(names); i++ {
		curr = f.client.Where("FIRST_NAME = ?", names[i])
	}
	if err := curr.Find(&user); err != nil {

	}

	var resp = map[string]interface{}{"matches": user}
	return resp
}

// MajorFilter filters
func (f Filters) MajorFilter(majors []string) map[string]interface{} {
	var user []models.User
	var curr *gorm.DB
	for i := 0; i < len(majors); i++ {
		curr = f.client.Where("MAJOR = ?", majors[i])
	}
	if err := curr.Find(&user); err != nil {

	}

	var resp = map[string]interface{}{"matches": user}
	return resp
}

// GradYearFilter filters
func (f Filters) GradYearFilter(gradYear []string) map[string]interface{} {
	var user []models.User
	var curr *gorm.DB
	for i := 0; i < len(gradYear); i++ {
		curr = f.client.Where("GRAD_YEAR = ?", gradYear[i])
	}
	if err := curr.Find(&user); err != nil {

	}

	var resp = map[string]interface{}{"matches": user}
	return resp
}

// InterestsFilter filters
func (f Filters) InterestsFilter(interests []string) map[string]interface{} {
	var user []models.User
	var curr *gorm.DB
	for i := 0; i < len(interests); i++ {
		curr = f.client.Where("INTERESTS = ?", interests[i])
	}
	if err := curr.Find(&user); err != nil {

	}

	var resp = map[string]interface{}{"matches": user}
	return resp
}

// LocationRadiusFilter filters
func (f Filters) LocationRadiusFilter(radius []string) map[string]interface{} {
	var user []models.User
	var curr *gorm.DB
	for i := 0; i < len(radius); i++ {
		read, _ := strconv.ParseFloat(radius[i], 64)
		curr = f.client.Where("LAT <= ?", read)
	}
	if err := curr.Find(&user); err != nil {

	}

	var resp = map[string]interface{}{"matches": user}
	return resp
}
