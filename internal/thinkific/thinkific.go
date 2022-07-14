package thinkific

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
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
	enrollments := types.Enrollments{}
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
	//fmt.Println(string(responseBody))
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

func GetMembers(ids []int) ([]types.WebhookOrder, error) {
	res := []types.WebhookOrder{}
	enrollments := types.Enrollments{}
	client := &http.Client{}

	for _, id := range ids {
		curPage := 1
		totalPages := 1
		for totalPages >= curPage {
			endpoint := fmt.Sprintf("https://api.thinkific.com/api/public/v1/enrollments?page=%v&limit=%v&query[course_id]=%v", curPage, 100, id)
			req, _ := http.NewRequest(http.MethodGet, endpoint, nil)

			req.Header.Add("X-Auth-API-Key", os.Getenv("THINKIFIC_API_KEY"))
			req.Header.Add("X-Auth-Subdomain", os.Getenv("THINKIFIC_SUBDOMAIN"))
			req.Header.Add("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				return res, err
			}

			responseBody, _ := ioutil.ReadAll(resp.Body)
			//fmt.Println(string(responseBody))
			json.Unmarshal(responseBody, &enrollments)

			for _, v := range enrollments.Items {
				if v.Expired || strings.Contains(v.User_name, "Test") {
					continue
				}

				l := types.WebhookOrder{}
				l.Payload.Course.Id = v.Course_id
				l.Payload.Course.Name = v.Course_name

				l.Payload.Expiry_date = v.Expiry_date

				l.Payload.User.Email = v.User_email
				l.Payload.User.Id = v.User_id
				l.Payload.User.First_name = v.User_name[:strings.Index(v.User_name, " ")]
				l.Payload.User.Last_name = v.User_name[strings.Index(v.User_name, " ")+1:]
				res = append(res, l)
			}

			curPage++
			totalPages = enrollments.Meta.Pagination.Total_pages

		}
	}

	return res, nil
}
