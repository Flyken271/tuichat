package main

import (
	"fmt"
	"os"

	"github.com/gempir/go-twitch-irc/v3"
	"github.com/joho/godotenv"
	"github.com/marcusolsson/tui-go"
	"github.com/marcusolsson/tui-go/wordwrap"
)

/*

!!!! OAuth Token Generator by Twitch !!!!

Use this site developed by Twitch (not me) to get your oauth token:

https://twitchapps.com/tmi/

*/

func connect(client *twitch.Client) {
	if err := client.Connect(); err != nil {
		panic(err)
	}
}

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Could not load environment variables: %e", err)
	}

	USERNAME := os.Getenv("USERNAME")
	OAUTH := os.Getenv("OAUTH")
	channel := os.Getenv("channel")

	// or client := twitch.NewAnonymousClient() for an anonymous user (no write capabilities)
	client := twitch.NewClient(USERNAME, OAUTH)

	client.Join(channel)

	go connect(client)

	chat := tui.NewList()
	chatbox := tui.NewScrollArea(chat)
	chatbox.SetAutoscrollToBottom(true)
	chatBorder := tui.NewVBox(chatbox)
	chatBorder.SetBorder(true)

	input := tui.NewEntry()
	inputbox := tui.NewHBox(input)
	inputbox.SetBorder(true)
	input.SetFocused(true)

	chatBorder.SetSizePolicy(tui.Expanding, tui.Expanding)
	inputbox.SetSizePolicy(tui.Maximum, tui.Minimum)

	chatcontainer := tui.NewVBox(chatBorder, inputbox)

	root := tui.NewVBox(chatcontainer)

	ui, err := tui.New(root)
	if err != nil {
		panic(err)
	}

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		chat.AddItems(message.User.DisplayName + ": " + wordwrap.WrapString(message.Message, 80))
		ui.Repaint()
	})

	input.OnSubmit(func(entry *tui.Entry) {
		client.Say(channel, entry.Text())
		chat.AddItems(USERNAME + ": " + wordwrap.WrapString(entry.Text(), 80))
		input.SetText("")
		ui.Repaint()
	})

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	if err := ui.Run(); err != nil {
		panic(err)
	}

}
