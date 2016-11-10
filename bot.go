package pinger

import (
	"os"
	"github.com/nlopes/slack"
	"fmt"
	"log"
)

type Bot struct {
	done            chan struct{}
	chSender        chan slack.OutgoingMessage
	chReceiver      chan slack.SlackEvent
	slackToken      string
	slackClient     *slack.Slack
	ws              *slack.SlackWS
	channelsByNames map[string]slack.Channel
}

var NewmanBot *Bot;

func init() {
	token := os.Getenv("slackToken")
	slackClient := slack.New(token)
	slackClient.SetDebug(true)
	NewmanBot = &Bot{make(chan struct{}), make(chan slack.OutgoingMessage), make(chan slack.SlackEvent), token, slackClient, nil, nil}
	NewmanBot.loadChannelsByNames()

}

func (bot *Bot) Stop() {
	close(bot.done)
}

func (bot *Bot) OnPingChange(transitions []Transition, total, fail int) {
	channel, ok := bot.channelsByNames["ops"]
	if !ok {
		log.Fatalf("Failed to find channel id for channel %q\n", "ops")
	}
	params := slack.NewPostMessageParameters()
	params.AsUser = true

	var msg string
	for _, transition := range transitions {
		name := transition.Address
		if transition.To == Success{
			name = fmt.Sprintf("~%s~", transition.Address)
		}else if transition.From ==  Success && transition.To == Fail{
			name = fmt.Sprintf("*%s*", transition.Address)
		}
		msg += fmt.Sprintf("%s\n", name)
	}
	log.Printf("msg is %q", msg)
	color := "#ff0000"
	if fail == 0{
		color = "#36a64f"
	}
	attachment := slack.Attachment{
		Text:    msg,
		Color: color,
		MarkdownIn: []string{"text", "title"},

	}
	params.Attachments = []slack.Attachment{attachment}

	channelID, timestamp, err := bot.slackClient.PostMessage(channel.Id, fmt.Sprintf("*Failed List Was Changed* *Total:* %d, *Fail:* %d", total, fail), params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)

}

func (bot *Bot) loadChannelsByNames() {
	channels, err := bot.slackClient.GetChannels(true)
	if err != nil {
		log.Fatalf("failed to get channles: %s\n", err.Error())
	}
	res := map[string]slack.Channel{}
	for _, channel := range channels {
		//log.Printf("got channel name: %s, id: %s\n", channel.Name, channel.Id)
		res[channel.Name] = channel
	}
	bot.channelsByNames = res
}

func (bot *Bot) postMessage(channelName, message string) {
	channel, ok := bot.channelsByNames[channelName]
	if !ok {
		log.Fatalf("Failed to find channel id for channel %q, message is %q\n", channelName, message)
	}
	params := slack.NewPostMessageParameters()
	params.AsUser = true
	channelID, timestamp, err := bot.slackClient.PostMessage(channel.Id, message, params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}

func (bot *Bot) Start() {
	//bot.postMessage("ops", "foo bar :+1:")
	/*
	wsAPI, err := bot.slackClient.StartRTM("", "http://example.com")
	if err != nil {
		fmt.Errorf("%s\n", err)
	}

	bot.ws = wsAPI
	go wsAPI.HandleIncomingEvents(bot.chReceiver)
	go wsAPI.Keepalive(20 * time.Second)
	go func(wsAPI *slack.SlackWS, chSender chan slack.OutgoingMessage) {
		for {
			select {
			case <-bot.done:
				return
			case msg := <-chSender:
				log.Printf(" -------- > sending message %#v\n", msg)
				wsAPI.SendMessage(&msg)
			}
		}
	}(wsAPI, bot.chSender)

	go func() {
		for {
			select {
			case <-bot.done:
				return
			case msg := <-bot.chReceiver:
				fmt.Printf("---> Event Received:  %#v\n", msg)
				switch msg.Data.(type) {
				case slack.HelloEvent:
					//bot.chSender <- *bot.ws.NewOutgoingMessage("Newman is in the building", channlesMap["newman"].Id)
					//bot.chSender <- *bot.ws.NewOutgoingMessage("Newman is in the building", bot.channelsByNames["ops"].Id)
					//bot.postMessage(channlesMap["newman"].Id, "Newman is in the building")
				case *slack.MessageEvent:
					evt := msg.Data.(*slack.MessageEvent)
					fmt.Printf("--- Message: %v\n", evt)
				case *slack.PresenceChangeEvent:
					evt := msg.Data.(*slack.PresenceChangeEvent)
					fmt.Printf("--- Presence Change: %v\n", evt)
					user, err := bot.slackClient.GetUserInfo(evt.UserId);
					if err != nil {
						fmt.Printf("err is: %#v\n", err)
					} else {
						fmt.Printf("user is: %#v\n", user.Name)
					}

				case slack.LatencyReport:
					a := msg.Data.(slack.LatencyReport)
					fmt.Printf("--- Current latency: %v\n", a.Value)
				case *slack.SlackWSError:
					err := msg.Data.(*slack.SlackWSError)
					fmt.Printf("--- Error: %d - %s\n", err.Code, err.Msg)
				default:
					fmt.Printf("-- Unexpected: %v#\n", msg.Data)
				}
			}
		}
	}()
	*/

}

