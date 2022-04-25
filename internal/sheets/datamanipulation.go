package sheets

import (
	"encoding/json"
	"log"
	"strconv"
	"thinkific-discord/internal/thinkific"
	"thinkific-discord/internal/types"

	"google.golang.org/api/sheets/v4"
)

const (
	enumCourses      = "enums!B5:C"
	enumRoles        = "enums!F5:G"
	settingsBindings = "settings!E4:G"
	dataRange        = "data!B4:H"
)

func UpdateCourses() {
	courses := thinkific.GetCourses()
	var vr sheets.ValueRange

	for _, v := range courses.Items {
		vr.Values = append(vr.Values, []interface{}{v.Name, v.Id})
	}
	writeRange := enumCourses
	_, err := svc.Spreadsheets.Values.Update(spreadsheetId, writeRange, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet. %v", err)
	}
}

func UpdateRoles(roles types.RolesResp) {
	var vr sheets.ValueRange
	for _, v := range roles {
		vr.Values = append(vr.Values, []interface{}{v.Name, v.Id})
	}
	writeRange := enumRoles
	_, err := svc.Spreadsheets.Values.Update(spreadsheetId, writeRange, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet. %v", err)
	}
}

func AddUser(user types.User) {
	readRange := dataRange
	resp, err := svc.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	var vr sheets.ValueRange
	found := false
	for _, row := range resp.Values {

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

	_, err = svc.Spreadsheets.Values.Update(spreadsheetId, writeRange, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet. %v", err)
	}
}

func GetDiscordIdByUserId(id int) string {
	readRange := dataRange
	resp, err := svc.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	for _, row := range resp.Values {

		if row[0].(string) == strconv.Itoa(id) {
			for len(row) < 5 {
				row = append(row, "")
			}
			return row[4].(string)
		}
	}
	return ""

}

func SetDiscordIdByUserId(id int, discorsId, discordUN string) {
	readRange := dataRange
	resp, err := svc.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}
	var vr sheets.ValueRange

	for _, row := range resp.Values {

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

	_, err = svc.Spreadsheets.Values.Update(spreadsheetId, writeRange, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet. %v", err)
	}

}

func GetUserRoles(userId int) []types.CurrentRole {
	currentRoles := []types.CurrentRole{}
	readRange := dataRange
	resp, err := svc.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	for _, row := range resp.Values {

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

func SetUserRoles(userId int, currentRoles []types.CurrentRole) {
	readRange := dataRange
	resp, err := svc.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}
	var vr sheets.ValueRange

	for _, row := range resp.Values {

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

	_, err = svc.Spreadsheets.Values.Update(spreadsheetId, writeRange, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet. %v", err)
	}

}

func getCourseRole(courseId int) string {
	roleId := ""
	readRange := settingsBindings
	resp, err := svc.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}
	for _, row := range resp.Values {

		if row[0].(string) == strconv.Itoa(courseId) {
			if len(row) < 2 {
				return ""
			}
			roleId = row[1].(string)
			return roleId
		}
	}
	return ""
}

func GetUsersRoles() []types.RolesWithIds {
	res := []types.RolesWithIds{}
	readRange := dataRange
	resp, err := svc.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
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
