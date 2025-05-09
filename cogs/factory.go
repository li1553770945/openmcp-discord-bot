package cogs

import (
	"context"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/disgoorg/disgo"
	disgobot "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
	"github.com/li1553770945/openmcp-discord-bot/cogs/model"
	"github.com/li1553770945/openmcp-discord-bot/infra/config"
	"log"
	"sync"
)

const MESSAGE_SEND_BUFFER_SIZE = 100

var (
	botOnce         sync.Once
	bot             disgobot.Client
	messageSendChan chan *model.MessageSendReq
)

func startMessageSender(ctx context.Context, wg *sync.WaitGroup) {
	messageSendChan = make(chan *model.MessageSendReq, MESSAGE_SEND_BUFFER_SIZE)
	go func() {
		wg.Add(1)
		defer wg.Done()
		var messageSendReq *model.MessageSendReq
		select {
		case messageSendReq = <-messageSendChan:

			var channelId uint64
			if messageSendReq.Channel == 0 {
				channelId = config.GetConfig().Discord.DefaultChannel
			} else {
				channelId = messageSendReq.Channel
			}

			_, err := bot.Rest().CreateMessage(snowflake.ID(channelId), discord.NewMessageCreateBuilder().SetContent(messageSendReq.Content).Build())
			if err != nil {
				logger.Errorf("发送消息到discord失败：%v", err)
			}
		case <-ctx.Done():

			return
		}
	}()
}

func InitGlobalBot(token string, ctx context.Context, wg *sync.WaitGroup) {
	botOnce.Do(func() {
		client, err := disgo.New(token,
			// set gateway options
			disgobot.WithGatewayConfigOpts(
				// set enabled intents
				gateway.WithIntents(
					gateway.IntentGuilds,
					gateway.IntentGuildMessages,
					gateway.IntentDirectMessages,
				),
			),
			// add event listeners
			disgobot.WithEventListenerFunc(func(e *events.MessageCreate) {
				// event code here
			}),
		)
		if err != nil {
			log.Fatal("init disgo bot client: ", err)
		}
		bot = client
		startMessageSender(ctx, wg)
	})

}

func GetBot() disgobot.Client {
	return bot
}
func GetMessageSendReqChan() chan *model.MessageSendReq {
	return messageSendChan
}
