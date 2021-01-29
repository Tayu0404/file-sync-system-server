package handler

import (
	"github.com/bwmarrin/snowflake"
	cache "github.com/patrickmn/go-cache"

	"github.com/Tayu0404/file-sync-system-server/api/model"
)

type handler struct {
	Model model.Database
	Node  *snowflake.Node
	Cache *cache.Cache
}

func NewHandler(d model.Database, n *snowflake.Node, c *cache.Cache) *handler {
	return &handler {
		Model: d,
		Node: n,
		Cache: c,
	}
}