package redis_article

import (
	"blogX_server/global"
	"blogX_server/utils/date"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type articleCacheType string

const (
	articleCacheLook    articleCacheType = "article_look_key"
	articleCacheFavor   articleCacheType = "article_favor_key"
	articleCacheCollect articleCacheType = "article_collect_key"
	articleCacheComment articleCacheType = "article_comment_key"
)

// Lookcount ----> id:3,num:1002
// 设置键值缓存
func set(key articleCacheType, articleID uint, n int) {
	num, _ := global.Redis.HGet(context.Background(), string(key), strconv.Itoa(int(articleID))).Int()
	num += n
	global.Redis.HSet(context.Background(), string(key), strconv.Itoa(int(articleID)), num)
}

// SetCacheLook 浏览量
func SetCacheLook(articleID uint, increase bool) {
	var n = 1
	if !increase {
		n = -1
	}
	set(articleCacheLook, articleID, n)
}

func SetCacheFavor(articleID uint, increase bool) {
	var n = 1
	if !increase {
		n = -1
	}
	set(articleCacheFavor, articleID, n)
}
func SetCacheCollect(articleID uint, increase bool) {
	var n = 1
	if !increase {
		n = -1
	}
	set(articleCacheCollect, articleID, n)
}

// SetCacheComment 设置评论数
func SetCacheComment(articleID uint, n int) {

	set(articleCacheComment, articleID, n)
}

// 从缓存取
func get(cacheType articleCacheType, articleID uint) int {
	num, _ := global.Redis.HGet(context.Background(), string(cacheType), strconv.Itoa(int(articleID))).Int()
	return num
}
func GetCacheLook(articleID uint) int {
	return get(articleCacheLook, articleID)
}
func GetCacheFavor(articleID uint) int {
	return get(articleCacheFavor, articleID)
}
func GetCacheCollect(articleID uint) int {
	return get(articleCacheCollect, articleID)
}
func GetCacheComment(articleID uint) int {
	return get(articleCacheComment, articleID)
}

// GetAll 取所有
func GetAll(cacheType articleCacheType) (mp map[uint]int) {
	res, err := global.Redis.HGetAll(context.Background(), string(cacheType)).Result()
	if err != nil {
		logrus.Errorf("redis取缓存错误:%s", err)
		return
	}
	mp = make(map[uint]int)
	for k, n := range res {
		key, err := strconv.Atoi(k)
		if err != nil {
			continue
		}
		num, err := strconv.Atoi(n)
		if err != nil {
			continue
		}
		mp[uint(key)] = num
	}
	return mp
}
func GetAllCacheLook() (mps map[uint]int) {
	return GetAll(articleCacheLook)
}
func GetAllCacheFavor() (mps map[uint]int) {
	return GetAll(articleCacheFavor)
}
func GetAllCacheCollect() (mps map[uint]int) {
	return GetAll(articleCacheCollect)
}
func GetAllCacheComment() (mps map[uint]int) {
	return GetAll(articleCacheComment)
}

func Clear() {
	err := global.Redis.Del(context.Background(), "article_look_key", "article_favor_key", "article_collect_key").Err()
	if err != nil {
		logrus.Error(err)
	}
}

// SetUserArticleHistoryCache 设置当日足迹缓存
func SetUserArticleHistoryCache(articleID, userID uint) {
	key := fmt.Sprintf("history_%d", userID)
	field := fmt.Sprintf("%d", articleID)

	endTime := date.GetNowAfter()
	ttl := time.Until(endTime)
	err := global.Redis.HSet(context.Background(), key, field, "").Err()
	if err != nil {
		logrus.Error(err)
		return
	}

	//设置过期时间0点
	//标记：ExpireAt用不了，为什么
	err = global.Redis.Expire(context.Background(), key, ttl).Err()
	if err != nil {
		logrus.Error(err)
		return
	}
}

// GetUserArticleHistoryCache 获取当日足迹缓存
func GetUserArticleHistoryCache(articleID, userID uint) (ok bool) {
	key := fmt.Sprintf("histroy_%d", userID)
	field := fmt.Sprintf("%d", articleID)
	err := global.Redis.HGet(context.Background(), key, field).Err()
	if err != nil {
		return false
	}
	return true
}
