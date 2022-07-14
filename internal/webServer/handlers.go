package webServer

import (
	"fmt"
	"thinkific-discord/internal/discord"
	"thinkific-discord/internal/discordBot"
	"thinkific-discord/internal/email"
	"thinkific-discord/internal/sheets"
	"thinkific-discord/internal/thinkific"
	"thinkific-discord/internal/types"
	"time"
)

func handleNewOrder(order types.WebhookOrder, forceSendInvite bool) {
	t1 := time.Now()
	isNew := sheets.AddUser(types.User{
		Id:        order.Payload.User.Id,
		FirstName: order.Payload.User.First_name,
		LastName:  order.Payload.User.Last_name,
		Email:     order.Payload.User.Email,
	})

	roles, err := sheets.AddCourseToUser(order.Payload.User.Id, order.Payload.Course.Id, order.Payload.Expiry_date)
	if roles == nil && err == nil {
		return
	}
	fmt.Println(time.Now().Sub(t1))
	if err != nil {
		panic(err)
	}
	go discordBot.SetRoles(order.Payload.User.Id, roles)

	if isNew || forceSendInvite {
		link := discord.GenerateLink(fmt.Sprintf("%v", order.Payload.User.Id))
		err = email.SendInviteLink(order.Payload.User.Email, link, order.Payload.User.First_name)
		if err != nil {
			panic(err)
		}
	}
}

func updateAllMembers() {
	ids, err := sheets.GetManagedCoursesId()
	if err != nil {
		panic(err)
	}
	orders, err := thinkific.GetMembers(ids)

	if err != nil {
		panic(err)
	}

	for _, v := range orders {
		handleNewOrder(v, false)
	}

	//fmt.Println(len(orders))

}
