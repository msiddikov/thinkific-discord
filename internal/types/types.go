package types

import "time"

type (
	User struct {
		Id        int
		FirstName string
		LastName  string
		Email     string
		DiscordId string
		DiscordUN string
	}
	RolesResp []struct {
		Id   string
		Name string
	}

	CurrentRole struct {
		RoleId string
		Expire time.Time
	}

	RolesWithIds struct {
		Id    int
		Roles []CurrentRole
	}

	WebhookOrder struct {
		Payload struct {
			Course struct {
				Id   int
				Name string
			}
			User struct {
				Email      string
				First_name string
				Id         int
				Last_name  string
			}
			Expiry_date time.Time
		}
	}

	Courses struct {
		Items []struct {
			Id       int `json:"id"`
			Name     string
			Duration string
		}
	}
)
