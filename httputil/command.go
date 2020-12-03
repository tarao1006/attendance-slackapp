package httputil

import (
	"context"

	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

type contextKey string

const commandKey contextKey = "command"

func SetSlashCommandToContext(ctx context.Context, command *slack.SlashCommand) context.Context {
	return context.WithValue(ctx, commandKey, command)
}

func GetSlashCommandFromContext(ctx context.Context) (*slack.SlashCommand, error) {
	v := ctx.Value(commandKey)
	command, ok := v.(*slack.SlashCommand)
	if !ok {
		return nil, errors.New("SlashCommand not found from context value")
	}
	return command, nil
}
