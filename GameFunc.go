package main

import (
	"fmt"
	"strings"
	"time"

	"Loup.Garou/config"
	Permissions "Loup.Garou/perm"

	"github.com/bwmarrin/discordgo"
)

func GameStart() {
	config.CurrentGame = config.Game{
		Started:  true,
		Finished: false,
		Players:  config.Players,
	}
}

func sendComp() {
	Comp := map[string]int{}
	Cards := map[string]string{}
	for i := range config.CurrentGame.Players {
		if config.CurrentGame.Players[i].Role.Name != "MDJ" {
			Comp[config.CurrentGame.Players[i].Role.Name]++
			Cards[config.CurrentGame.Players[i].Role.Name] = config.CurrentGame.Players[i].Role.Image
		}
	}

	WillSend := mkimage(Comp, Cards)

	for i := range WillSend {
		Image := &discordgo.File{
			Reader: WillSend[i],
			Name:   "compo.png",
		}
		Params := &discordgo.MessageSend{
			Files: []*discordgo.File{Image},
			Embed: &discordgo.MessageEmbed{
				Title: "Composition",
				Image: &discordgo.MessageEmbedImage{
					URL: "attachment://compo.png",
				},
				Color: config.EmbedColor,
			},
		}
		_, err := dg.ChannelMessageSendComplex(config.CurrentGame.GameStats.ID, Params)
		if err != nil {
			fmt.Println(err)
			//
		}
		time.Sleep(config.SleepTime * time.Millisecond)
	}
}

func ChannelReload() {
	ChannelDelete()
	ChannelGenerator()
}

func ChannelDelete() {
	channels, err := dg.GuildChannels(config.GuildID)
	if err != nil {
		return
	}
	for k := range channels {
		for i := range config.Channels {
			if config.Channels[i] == channels[k].Name {
				_, err := dg.ChannelDelete(channels[k].ID)
				if err != nil {
					return
				}
				time.Sleep(config.SleepTime * time.Millisecond)
			}
		}
		for i := range config.TeamChannels {
			if config.TeamChannels[i] == channels[k].Name {
				_, err := dg.ChannelDelete(channels[k].ID)
				if err != nil {
					return
				}
				time.Sleep(config.SleepTime * time.Millisecond)
			}
		}
		for i := range config.SpecialChannels {
			if config.SpecialChannels[i] == channels[k].Name {
				_, err := dg.ChannelDelete(channels[k].ID)
				if err != nil {
					return
				}
				time.Sleep(config.SleepTime * time.Millisecond)
			}
		}
	}
	if &config.CurrentGame.Deads == nil {
		_, err := dg.ChannelDelete(config.CurrentGame.Deads.ID)
		if err != nil {
			return
		}
		time.Sleep(config.SleepTime * time.Millisecond)
	}
	if &config.CurrentGame.GameStats == nil {
		_, err := dg.ChannelDelete(config.CurrentGame.GameStats.ID)
		if err != nil {
			return
		}
		time.Sleep(config.SleepTime * time.Millisecond)
	}
	if &config.CurrentGame.Votes == nil {
		_, err := dg.ChannelDelete(config.CurrentGame.Votes.ID)
		if err != nil {
			return
		}
	}
}

func GenerateSpecials() {
	Perms := []*discordgo.PermissionOverwrite{{
		ID:    config.GuildID,
		Type:  "role",
		Allow: Permissions.VIEW_CHANNEL,
		Deny:  Permissions.SEND_MESSAGES,
	}, {
		ID:    dg.State.User.ID,
		Type:  "member",
		Allow: Permissions.VIEW_CHANNEL + Permissions.SEND_MESSAGES,
	}}
	data := discordgo.GuildChannelCreateData{
		Name:                 "lg-gamestats",
		Type:                 discordgo.ChannelTypeGuildText,
		ParentID:             config.CategoryID,
		PermissionOverwrites: Perms,
	}
	channel, err := dg.GuildChannelCreateComplex(config.GuildID, data)
	if err != nil {
		//
	}
	config.CurrentGame.GameStats = channel
	time.Sleep(config.SleepTime * time.Millisecond)
	Perms = []*discordgo.PermissionOverwrite{{
		ID:   config.GuildID,
		Type: "role",
		Deny: Permissions.VIEW_CHANNEL,
	}, {
		ID:    dg.State.User.ID,
		Type:  "member",
		Allow: Permissions.VIEW_CHANNEL + Permissions.SEND_MESSAGES,
	}}
	data = discordgo.GuildChannelCreateData{
		Name:                 "lg-vote",
		Type:                 discordgo.ChannelTypeGuildText,
		ParentID:             config.CategoryID,
		PermissionOverwrites: Perms,
	}
	channel, err = dg.GuildChannelCreateComplex(config.GuildID, data)
	if err != nil {
		//
	}
	config.CurrentGame.Votes = channel
	time.Sleep(config.SleepTime * time.Millisecond)
	data = discordgo.GuildChannelCreateData{
		Name:                 "lg-morts",
		Type:                 discordgo.ChannelTypeGuildText,
		ParentID:             config.CategoryID,
		PermissionOverwrites: Perms,
	}
	channel, err = dg.GuildChannelCreateComplex(config.GuildID, data)
	if err != nil {
		//
	}
	config.CurrentGame.Deads = channel
	time.Sleep(config.SleepTime * time.Millisecond)
}

func ChannelGenerator() {
	toCreate := []string{}
	for i := range config.CurrentGame.Players {
		existNormal := false
		existTeam := false
		for k := range toCreate {
			if toCreate[k] == config.CurrentGame.Players[i].Role.ChannelName {
				existNormal = true
			}
			if toCreate[k] == config.CurrentGame.Players[i].Role.Team.ChannelName {
				existTeam = true
				break
			} else {
				break
			}
		}
		if !existNormal && config.CurrentGame.Players[i].Role.ChannelName != "" {
			toCreate = append(toCreate, config.CurrentGame.Players[i].Role.ChannelName)
		}
		if !existTeam && config.CurrentGame.Players[i].Role.Team.ChannelName != "" {
			toCreate = append(toCreate, config.CurrentGame.Players[i].Role.Team.ChannelName)
		}
	}
	MDJid := ""
	for i := range config.CurrentGame.Players {
		if config.CurrentGame.Players[i].Role.Name == "MDJ" {
			MDJid = config.CurrentGame.Players[i].ID
		}
	}

	for i := range toCreate {
		Perms := []*discordgo.PermissionOverwrite{{
			ID:   config.GuildID,
			Type: "role",
			Deny: Permissions.VIEW_CHANNEL,
		}, {
			ID:    dg.State.User.ID,
			Type:  "member",
			Allow: Permissions.VIEW_CHANNEL,
		}, {
			ID:    MDJid,
			Type:  "member",
			Allow: Permissions.VIEW_CHANNEL,
		}}
		data := discordgo.GuildChannelCreateData{
			Name:                 toCreate[i],
			Type:                 discordgo.ChannelTypeGuildText,
			ParentID:             config.CategoryID,
			PermissionOverwrites: Perms,
		}
		channel, err := dg.GuildChannelCreateComplex(config.GuildID, data)
		if err != nil {
			fmt.Println(err)
		}
		config.CurrentGame.Channels = append(config.CurrentGame.Channels, channel)
		time.Sleep(config.SleepTime * time.Millisecond)
	}
	GenerateSpecials()
}

func DMSender() {
	for i := range config.CurrentGame.Players {
		if config.CurrentGame.Players[i].Role.Name != "MDJ" {
			dm, err := dg.UserChannelCreate(config.CurrentGame.Players[i].ID)
			if err != nil {
				fmt.Println(err)
				return
			}
			dat, err := box.FindString("static/" + config.CurrentGame.Players[i].Role.Image)
			check(err)
			Image := &discordgo.File{
				Reader: strings.NewReader(dat),
				Name:   "role.png",
			}
			Params := &discordgo.MessageSend{
				Files: []*discordgo.File{Image},
				Embed: &discordgo.MessageEmbed{
					Title:       config.CurrentGame.Players[i].Role.Name,
					Description: config.CurrentGame.Players[i].Role.Description,
					Image: &discordgo.MessageEmbedImage{
						URL: "attachment://role.png",
					},
					Color: config.EmbedColor,
				},
			}
			_, err = dg.ChannelMessageSendComplex(dm.ID, Params)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
	for i := range config.CurrentGame.Players {
		if config.CurrentGame.Players[i].Role.Name == "MDJ" {
			dm, err := dg.UserChannelCreate(config.CurrentGame.Players[i].ID)
			if err != nil {
				return
			}
			embed := &discordgo.MessageEmbed{
				Title:       "Loup-Garou",
				Description: "Les rôles ont étés distribués",
				Color:       config.EmbedColor,
			}
			_, err = dg.ChannelMessageSendEmbed(dm.ID, embed)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func GameBegin(Data Received) {
	for i := range Data.Game.Players {
		for k := range config.CurrentGame.Players {
			if config.CurrentGame.Players[k].ID == Data.Game.Players[i].ID {
				rl, error := config.GetRoleByName(Data.Game.Players[i].RoleName)
				if error {
					return
				}
				config.CurrentGame.Players[k].Role = rl
			}
		}
	}
	go ChannelGenerator()
	go DMSender()
}
