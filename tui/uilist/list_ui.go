package uilist

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-redis/redis/v8"
	"github.com/mxmaxime/tedis/myredis"
)

type ListModel struct {
	RedisRepo *myredis.RedisRepo
	list      list.Model
}

// Init run any intial IO on program start
func (m ListModel) Init() tea.Cmd {
	return nil
}

func New(redis_cli *redis.Client) *ListModel {
	ctx := context.TODO()
	repo := myredis.RedisRepo{Cli: redis_cli}
	// init list
	// todo: move it in message stuff
	var items []list.Item

	keys, err := repo.GetKeys(ctx, 0, "", -1)
	if err != nil {
		panic(err)
	}

	for _, key := range keys {
		keyType := redis_cli.Type(ctx, key).Val()
		//fmt.Printf("key: %s kt %s\n", key, kt)
		items = append(items, ListItem{
			Key:     key,
			KeyType: keyType,
		})
	}

	fmt.Printf("found %s keys\n", len(keys))

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Keys"
	//l.SetShowFilter(false)
	//l.SetShowHelp(false)
	//l.SetFilteringEnabled(false)

	model := ListModel{
		RedisRepo: &repo,
		list:      l,
	}

	return &model
}
