package thinkific

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"thinkific-discord/internal/types"
)

func GetCourses() types.Courses {
	courses := types.Courses{}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, "https://api.thinkific.com/api/public/v1/courses", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("X-Auth-API-Key", os.Getenv("THINKIFIC_API_KEY"))
	req.Header.Add("X-Auth-Subdomain", os.Getenv("THINKIFIC_SUBDOMAIN"))
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Unable to get courses list:" + err.Error())
		return courses
	}

	responseBody, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(responseBody, &courses)
	return courses
}
