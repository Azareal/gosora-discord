package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Azareal/Gosora/common"
)

func init() {
	common.Plugins.Add(&common.Plugin{UName: "discord", Name: "Discord", Author: "Azareal", Init: discordInit, Activate: discordActivate, Deactivate: discordDeactivate})
}

func discordValidate() error {
	webhook, ok := common.PluginConfig["DiscordWebhook"]
	if !ok || webhook == "" {
		return errors.New("You need to set a webhook to push to in config.json")
	}

	ev, ok := common.PluginConfig["DiscordEvents"]
	if ev != "" && ev != "threads" && ev != "replies" {
		return errors.New("Invalid value for DiscordEvents. Can only be blank, 'threads' or 'replies'")
	}

	fidsRaw := common.PluginConfig["DiscordForums"]
	if fidsRaw != "" {
		for _, fidRaw := range strings.Split(fidsRaw, ",") {
			_, err := strconv.Atoi(fidRaw)
			if err != nil {
				return errors.New("Invalid integer found in DiscordForums")
			}
		}
	}

	return nil
}

func discordInit(plugin *common.Plugin) error {
	err := discordValidate()
	if err != nil {
		return err
	}
	plugin.AddHook("action_end_create_topic", discordEventTopic)
	plugin.AddHook("action_end_create_reply", discordEventReply)
	return nil
}

// A bit of validation to make sure the admin isn't forgetting something or telling Plugin Discord to do something absurd
func discordActivate(plugin *common.Plugin) error {
	return discordValidate()
}

func discordDeactivate(plugin *common.Plugin) {
	plugin.RemoveHook("action_end_create_topic", discordEventTopic)
	plugin.RemoveHook("action_end_create_reply", discordEventReply)
}

func discordEventTopic(args ...interface{}) (skip bool, rerr common.RouteError) {
	discordEvent(0, args[0].(int))
	return false, nil
}
func discordEventReply(args ...interface{}) (skip bool, rerr common.RouteError) {
	discordEvent(1, args[0].(int))
	return false, nil
}

type DiscordData struct {
	Content  string `json:"content"`
	Username string `json:"username"`
}

func discordEvent(typ int, id int) {
	ev, ok := common.PluginConfig["DiscordEvents"]
	if !ok {
		return
	}
	if (ev == "threads" && typ != 0) || (ev == "replies" && typ != 1) {
		return
	}

	var content string
	var topic *common.Topic
	var err error
	if typ == 0 {
		topic, err = common.Topics.Get(id)
		if err != nil {
			return
		}
		content = topic.Content
	} else {
		reply, err := common.Rstore.Get(id)
		if err != nil {
			return
		}
		content = reply.Content

		topic, err = reply.Topic()
		if err != nil {
			return
		}
	}

	fidsRaw := common.PluginConfig["DiscordForums"]
	if fidsRaw != "" {
		var hasForum = false
		for _, fidRaw := range strings.Split(fidsRaw, ",") {
			fid, err := strconv.Atoi(fidRaw)
			if err != nil {
				return
			}
			if fid == topic.ParentID {
				hasForum = true
			}
		}
		if !hasForum {
			return
		}
	}
	if len(content) > 100 {
		content = content[:97] + "..."
	}

	var client = &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			Dial:                (&net.Dialer{Timeout: 5 * time.Second}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}
	dat := DiscordData{Content: content, Username: common.Site.Name}
	data, err := json.Marshal(dat)
	if err != nil {
		common.LogWarning(err)
		return
	}
	_, err = client.Post(common.PluginConfig["DiscordWebhook"], "application/json", bytes.NewBuffer(data))
	if err != nil {
		common.LogWarning(err)
	}
}

// TODO: Add a settings page or something?
