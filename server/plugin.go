package main

import (
	"strings"
	"sync"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
}

const EMAIL = "@mattermost.com"

func (p *Plugin) MessageWillBePosted(c *plugin.Context, post *model.Post) (*model.Post, string) {
	channel, err := p.API.GetChannel(post.ChannelId)
	if err != nil {
		// Don't block posts in case of error
		return post, ""
	}

	if channel.Type != model.ChannelTypeDirect {
		return post, ""
	}

	members, err := p.API.GetChannelMembers(channel.Id, 0, 2)
	if err != nil || len(members) != 2 {
		return post, ""
	}

	user1, err := p.API.GetUser(members[0].UserId)
	if err != nil {
		return post, ""
	}

	user2, err := p.API.GetUser(members[1].UserId)
	if err != nil {
		return post, ""
	}

	if strings.HasSuffix(user1.Email, EMAIL) && strings.HasSuffix(user2.Email, EMAIL) {
		p.API.SendEphemeralPost(post.UserId, &model.Post{
			ChannelId: post.ChannelId,
			Message:   "Contact your colleagues in the hub server",
			RootId:    post.RootId,
		})
		return nil, "plugin.message_will_be_posted.dismiss_post"
	}

	return post, ""
}
