package cache

import (
	"context"
	"errors"
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"sync"
	"time"
)

// 技术选型考虑的点
//  1. 功能性：功能是否能够完全覆盖你的需求。
//  2. 社区和支持度：社区是否活跃，文档是否齐全，
//     以及百度（搜索引擎）能不能搜索到你需要的各种信息，有没有帮你踩过坑
//  3. 非功能性：易用性（用户友好度，学习曲线要平滑），
//     扩展性（如果开源软件的某些功能需要定制，框架是否支持定制，以及定制的难度高不高）
//     性能（追求性能的公司，往往有能力自研）

// LocalCodeCache 本地缓存实现
type LocalCodeCache struct {
	cache *lru.Cache
	// 普通锁，或者说写锁
	lock sync.Mutex
	//读写锁
	expiration time.Duration
}
