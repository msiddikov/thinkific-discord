package thinkific

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"thinkific-discord/internal/types"
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
