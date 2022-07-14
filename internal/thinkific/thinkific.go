package thinkific

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"thinkific-discord/internal/types"
	"time"
)

func GetCourses() types.Courses {
	courses := types.Courses{}

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, "https://api.thinkific.com/api/public/v1/courses", nil)

	req.Header.Add("X-Auth-API-Key", os.Getenv("THINKIFIC_API_KEY"))
	req.Header.Add("X-Auth-Subdomain", os.Getenv("THINKIFIC_SUBDOMAIN"))
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Panic(err)
	}

	responseBody, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(responseBody, &courses)
	return courses
}

func GetExpiryDate(userId, courseId int) (time.Time, error) {
	enrollments := struct {
		Items []struct {
			Expired bool
		}
	}{}
	client := &http.Client{}
	endpoint := fmt.Sprintf("https://api.thinkific.com/api/public/v1/enrollments?query[user_id]=%v&query[course_id]=%v", userId, courseId)
	req, _ := http.NewRequest(http.MethodGet, endpoint, nil)

	req.Header.Add("X-Auth-API-Key", os.Getenv("THINKIFIC_API_KEY"))
	req.Header.Add("X-Auth-Subdomain", os.Getenv("THINKIFIC_SUBDOMAIN"))
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Panic(err)
	}

	responseBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(responseBody))
	json.Unmarshal(responseBody, &enrollments)

	if len(enrollments.Items) > 0 {
		if !enrollments.Items[0].Expired {
			return time.Now().Add(24 * time.Hour), nil
		}
		return time.Now().Add(-1 * time.Hour), nil
	} else {
		return time.Now().Add(-1 * time.Hour), fmt.Errorf("No enrollmant found")
	}
}
