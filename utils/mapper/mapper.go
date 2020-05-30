package mapper

import (
	//"strconv"

	"github.com/ucladevx/BConnect-backend/models"
	"strconv"
)

func UsersToMap(users []*models.User) map[string]interface{} {
	numUsers := len(users)
	userList := make([]map[string]string, numUsers)

	for _, val := range users {
		userList = append(userList, map[string]string{
			"userid":      val.UserID,
			"firstname":   val.FirstName,
			"lastname":    val.LastName,
			//"age":         strconv.Atoi(val.Age),
			"gradyear":    val.GradYear,
			//"currentjob":  val.CurrentJob,
			//"gender":      val.Gender,
			"email":       val.Email,
			"major":       val.Major,
			"phonenumber": val.PhoneNumber,
			"profilepic":  val.ProfilePic,
			"bio":         val.Bio,
		})
	}

	return map[string]interface{}{"num_users": numUsers, "users": userList}
}

func InterestsToMap(interests []*models.Interest) map[string]interface{} {
	numInterests := len(interests)
	interestList := make([]string, numInterests)

	for _, val := range interests {
		interestList = append(interestList, val.Interest)
	}

	return map[string]interface{}{"num_interests": numInterests, "interests": interestList}
}

func ClubsToMap(clubs []*models.Club) map[string]interface{} {
	numClubs := len(clubs)
	clubList := make([]string, numClubs)

	for _, val := range clubs {
		clubList = append(clubList, val.Club)
	}

	return map[string]interface{}{"num_clubs": numClubs, "clubs": clubList}
}
