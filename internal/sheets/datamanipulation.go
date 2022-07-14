package sheets

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	"thinkific-discord/internal/tgbot"
	"thinkific-discord/internal/thinkific"
	"thinkific-discord/internal/types"
	"time"

	"google.golang.org/api/sheets/v4"
)

const (
	enumCourses       = "enums!B6:C"
	enumRoles         = "enums!F5:G"
	settingsBindings  = "settings!E4:G"
	dataRange         = "data!B4:H"
	dataRangeStart    = "data!B"
	dataRangeEnd      = "H"
	datarangeStartRow = 4
)

func UpdateCourses() {
	defer func() {
		err := recover()
		if err != nil {
			tgbot.SendString(fmt.Sprint(err))
		}
		return
	}()
	courses := thinkific.GetCourses()
	var vr sheets.ValueRange

	for _, v := range courses.Items {
		vr.Values = append(vr.Values, []interface{}{v.Name, v.Id})
	}
	// adding some whitespace for deleted courses
	for i := 0; i < 3; i++ {
		vr.Values = append(vr.Values, []interface{}{"", ""})
	}

	writeRange := enumCourses
	err := update(spreadsheetId, writeRange, &vr)
	if err != nil {
		log.Panicf("Unable to retrieve data from sheet. %v", err)
	}
}

func UpdateRoles(roles types.RolesResp) {
	var vr sheets.ValueRange
	for _, v := range roles {
		vr.Values = append(vr.Values, []interface{}{v.Name, v.Id})
	}
	// adding some whitespace for deleted roles
	for i := 0; i < 3; i++ {
		vr.Values = append(vr.Values, []interface{}{"", ""})
	}

	writeRange := enumRoles
	err := update(spreadsheetId, writeRange, &vr)
	if err != nil {
		log.Panicf("Unable to retrieve data from sheet. %v", err)
	}
}

func AddUser(user types.User) {
	readRange := dataRange
	resp, err := get(spreadsheetId, readRange)
	if err != nil {
		log.Panicf("Unable to retrieve data from sheet: %v", err)
	}

	var vr sheets.ValueRange
	found := false
	for _, row := range resp.Values {

		if len(row) == 0 {
			continue
		}
		if row[0].(string) == strconv.Itoa(user.Id) {
			if len(row) > 5 {
				user.DiscordId = row[4].(string)
				user.DiscordUN = row[5].(string)
			}
			vr.Values = append(vr.Values, []interface{}{
				strconv.Itoa(user.Id),
				user.FirstName,
				user.LastName,
				user.Email,
				user.DiscordId,
				user.DiscordUN,
			})
			found = true
			continue
		}
		vr.Values = append(vr.Values, row)
	}

	if !found {
		vr.Values = append(vr.Values, []interface{}{
			strconv.Itoa(user.Id),
			user.FirstName,
			user.LastName,
			user.Email,
			user.DiscordId,
			user.DiscordUN,
		})
	}
	writeRange := dataRange
	err = update(spreadsheetId, writeRange, &vr)
	if err != nil {
		log.Panicf("Unable to retrieve data from sheet. %v", err)
	}
}

func GetDiscordIdByUserId(id int) string {
	readRange := dataRange
	resp, err := get(spreadsheetId, readRange)
	if err != nil {
		log.Panicf("Unable to retrieve data from sheet: %v", err)
	}

	for _, row := range resp.Values {

		if len(row) == 0 {
			continue
		}
		if row[0].(string) == strconv.Itoa(id) {
			for len(row) < 5 {
				return ""
			}
			return row[4].(string)
		}
	}
	return ""

}

func SetDiscordIdByUserId(id int, discorsId, discordUN string) error {
	readRange := dataRange
	resp, err := get(spreadsheetId, readRange)
	if err != nil {
		return err
	}
	var vr sheets.ValueRange

	for _, row := range resp.Values {

		if len(row) == 0 {
			continue
		}
		if row[0].(string) == strconv.Itoa(id) {
			for len(row) < 6 {
				row = append(row, "")
			}
			row[4] = discorsId
			row[5] = discordUN

		}
		vr.Values = append(vr.Values, row)
	}
	writeRange := dataRange

	err = update(spreadsheetId, writeRange, &vr)
	if err != nil {
		return err
	}
	return nil
}

func GetUserRoles(userId int) []types.CurrentRole {
	currentRoles := []types.CurrentRole{}
	readRange := dataRange
	resp, err := get(spreadsheetId, readRange)
	if err != nil {
		log.Panicf("Unable to retrieve data from sheet: %v", err)
	}

	for _, row := range resp.Values {

		if len(row) == 0 {
			continue
		}
		if row[0].(string) == strconv.Itoa(userId) {

			for len(row) < 7 {
				row = append(row, "")
			}
			if row[6].(string) != "" {
				json.Unmarshal([]byte(row[6].(string)), &currentRoles)
			}
		}
	}
	return currentRoles

}

func SetUserRoles(userId int, currentRoles []types.CurrentRole) error {
	readRange := dataRange
	resp, err := get(spreadsheetId, readRange)
	if err != nil {
		return err
	}
	var vr sheets.ValueRange

	for _, row := range resp.Values {
		if len(row) == 0 {
			continue
		}
		if len(row) == 0 {
			continue
		}
		if row[0].(string) == strconv.Itoa(userId) {

			for len(row) < 7 {
				row = append(row, "")
			}

			bytes, _ := json.Marshal(currentRoles)
			row[6] = string(bytes)
		}
		vr.Values = append(vr.Values, row)
	}
	writeRange := dataRange

	err = update(spreadsheetId, writeRange, &vr)
	if err != nil {
		return err
	}
	return nil
}

func GetCourseRole(courseId int) (string, error) {
	roleId := ""
	readRange := settingsBindings
	resp, err := get(spreadsheetId, readRange)
	if err != nil {
		return "", err
	}
	for _, row := range resp.Values {

		if len(row) == 0 {
			continue
		}
		if row[0].(string) == strconv.Itoa(courseId) {
			if len(row) < 2 {
				return "", nil
			}
			roleId = row[1].(string)
			return roleId, nil
		}
	}
	return "", nil
}

func GetRoleCourses(roleId string) ([]int, error) {
	courses := []int{}
	readRange := settingsBindings
	resp, err := get(spreadsheetId, readRange)
	if err != nil {
		return courses, err
	}

	for _, row := range resp.Values {

		if len(row) < 2 {
			return courses, fmt.Errorf("Not found")
		}

		if row[1].(string) == roleId {
			courseId, err := strconv.Atoi(row[0].(string))
			if err != nil {
				return courses, err
			}
			courses = append(courses, courseId)
		}
	}
	return courses, nil
}

func GetUsersRoles() []types.RolesWithIds {
	res := []types.RolesWithIds{}
	readRange := dataRange
	resp, err := get(spreadsheetId, readRange)
	if err != nil {
		log.Panicf("Unable to retrieve data from sheet: %v", err)
	}
	for _, row := range resp.Values {

		if len(row) < 7 {
			continue
		}
		roles := []types.CurrentRole{}
		json.Unmarshal([]byte(row[6].(string)), &roles)
		if len(roles) == 0 {
			continue
		}
		userId, _ := strconv.Atoi(row[0].(string))
		res = append(res, types.RolesWithIds{Id: userId, Roles: roles})
	}
	return res
}

func GetUserRow(userId int) ([]interface{}, string) {
	res := []interface{}{}
	readRange := dataRange
	resp, err := get(spreadsheetId, readRange)
	if err != nil {
		log.Panicf("Unable to retrieve data from sheet: %v", err)
	}
	resRow := len(resp.Values)
	for k, row := range resp.Values {
		if len(row) == 0 {
			continue
		}
		if row[0].(string) == strconv.Itoa(userId) {
			res = row
			resRow = datarangeStartRow + k
		}
	}

	for len(res) < 7 {
		res = append(res, "")
	}
	resRange := dataRangeStart + strconv.Itoa(resRow) + ":" + dataRangeEnd + strconv.Itoa(resRow)
	return res, resRange
}

func SetUserRow(row []interface{}, writeRange string) error {
	var vr sheets.ValueRange
	vr.Values = append(vr.Values, row)
	err := update(spreadsheetId, writeRange, &vr)
	if err != nil {
		log.Panicf("Unable to write data to sheets: %v", err)
	}
	return nil
}

func update(ss, r string, vr *sheets.ValueRange) error {
	const maxRetryPeriod = 600
	const maxRetryNo = 10
	currentRetries := 0

	for {
		res, err := svc.Spreadsheets.Values.Update(ss, r, vr).ValueInputOption("RAW").Do()
		if res.ServerResponse.HTTPStatusCode == 429 {
			if currentRetries == maxRetryNo {
				return err
			}

			sleepTime := math.Max(float64(maxRetryPeriod), math.Pow(2, float64(currentRetries)))
			time.Sleep(time.Duration(sleepTime) * time.Second)
			currentRetries++
			continue
		}

		return err

	}
}

func get(ss, r string) (*sheets.ValueRange, error) {
	const maxRetryPeriod = 600
	const maxRetryNo = 10
	currentRetries := 0

	for {
		res, err := svc.Spreadsheets.Values.Get(ss, r).Do()
		if res.ServerResponse.HTTPStatusCode == 429 {
			if currentRetries == maxRetryNo {
				return nil, err
			}

			sleepTime := math.Max(float64(maxRetryPeriod), math.Pow(2, float64(currentRetries)))
			time.Sleep(time.Duration(sleepTime) * time.Second)
			currentRetries++
			continue
		}

		return res, err

	}
}
