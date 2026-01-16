package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/ems/backend/pkg/config"
	"github.com/redis/go-redis/v9"
)

var (
	Client *redis.Client
	Ctx    = context.Background()
)

// Init initializes the Redis client
func Init() error {
	Client = redis.NewClient(&redis.Options{
		Addr:     config.Cfg.Redis.Addr(),
		Password: config.Cfg.Redis.Password,
		DB:       config.Cfg.Redis.DB,
		PoolSize: config.Cfg.Redis.PoolSize,
	})

	// Test connection
	if err := Client.Ping(Ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return nil
}

// Close closes the Redis connection
func Close() error {
	if Client != nil {
		return Client.Close()
	}
	return nil
}

// CacheService provides caching operations
type CacheService struct {
	client *redis.Client
}

// NewCacheService creates a new cache service
func NewCacheService() *CacheService {
	return &CacheService{client: Client}
}

// Set stores a value in cache with expiration
func (c *CacheService) Set(key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(Ctx, key, value, expiration).Err()
}

// Get retrieves a value from cache
func (c *CacheService) Get(key string) (string, error) {
	return c.client.Get(Ctx, key).Result()
}

// GetBytes retrieves a value from cache as bytes
func (c *CacheService) GetBytes(key string) ([]byte, error) {
	return c.client.Get(Ctx, key).Bytes()
}

// Del deletes keys from cache
func (c *CacheService) Del(keys ...string) error {
	return c.client.Del(Ctx, keys...).Err()
}

// Exists checks if keys exist
func (c *CacheService) Exists(keys ...string) (int64, error) {
	return c.client.Exists(Ctx, keys...).Result()
}

// Expire sets expiration time for a key
func (c *CacheService) Expire(key string, expiration time.Duration) error {
	return c.client.Expire(Ctx, key, expiration).Err()
}

// TTL returns the remaining time to live of a key
func (c *CacheService) TTL(key string) (time.Duration, error) {
	return c.client.TTL(Ctx, key).Result()
}

// Incr increments the value of a key
func (c *CacheService) Incr(key string) (int64, error) {
	return c.client.Incr(Ctx, key).Result()
}

// Decr decrements the value of a key
func (c *CacheService) Decr(key string) (int64, error) {
	return c.client.Decr(Ctx, key).Result()
}

// HSet sets a field in a hash
func (c *CacheService) HSet(key, field string, value interface{}) error {
	return c.client.HSet(Ctx, key, field, value).Err()
}

// HGet gets a field from a hash
func (c *CacheService) HGet(key, field string) (string, error) {
	return c.client.HGet(Ctx, key, field).Result()
}

// HGetAll gets all fields from a hash
func (c *CacheService) HGetAll(key string) (map[string]string, error) {
	return c.client.HGetAll(Ctx, key).Result()
}

// HDel deletes fields from a hash
func (c *CacheService) HDel(key string, fields ...string) error {
	return c.client.HDel(Ctx, key, fields...).Err()
}

// ZAdd adds a member to a sorted set
func (c *CacheService) ZAdd(key string, score float64, member string) error {
	return c.client.ZAdd(Ctx, key, redis.Z{Score: score, Member: member}).Err()
}

// ZRange returns a range of members from a sorted set
func (c *CacheService) ZRange(key string, start, stop int64) ([]string, error) {
	return c.client.ZRange(Ctx, key, start, stop).Result()
}

// ZRemRangeByRank removes members from a sorted set by rank range
func (c *CacheService) ZRemRangeByRank(key string, start, stop int64) error {
	return c.client.ZRemRangeByRank(Ctx, key, start, stop).Err()
}

// FlushAll clears all keys in the current database
func FlushAll() error {
	return Client.FlushDB(Ctx).Err()
}

// Cache key prefixes for different data types
const (
	KeyPrefixUser       = "ems:user:"
	KeyPrefixEquipment  = "ems:equipment:"
	KeyPrefixEquipmentType = "ems:equipment_type:"
	KeyPrefixWorkshop   = "ems:workshop:"
	KeyPrefixFactory    = "ems:factory:"
	KeyPrefixInspectionTemplate = "ems:inspection_template:"
	KeyPrefixStatistics = "ems:stats:"
)

// BuildUserKey builds a cache key for user data
func BuildUserKey(userID uint) string {
	return fmt.Sprintf("%s%d", KeyPrefixUser, userID)
}

// BuildEquipmentKey builds a cache key for equipment data
func BuildEquipmentKey(equipmentID uint) string {
	return fmt.Sprintf("%s%d", KeyPrefixEquipment, equipmentID)
}

// BuildEquipmentQRKey builds a cache key for equipment by QR code
func BuildEquipmentQRKey(qrCode string) string {
	return fmt.Sprintf("%sqr:%s", KeyPrefixEquipment, qrCode)
}

// BuildEquipmentTypeKey builds a cache key for equipment type data
func BuildEquipmentTypeKey(typeID uint) string {
	return fmt.Sprintf("%s%d", KeyPrefixEquipmentType, typeID)
}

// BuildWorkshopKey builds a cache key for workshop data
func BuildWorkshopKey(workshopID uint) string {
	return fmt.Sprintf("%s%d", KeyPrefixWorkshop, workshopID)
}

// BuildFactoryKey builds a cache key for factory data
func BuildFactoryKey(factoryID uint) string {
	return fmt.Sprintf("%s%d", KeyPrefixFactory, factoryID)
}

// BuildInspectionTemplateKey builds a cache key for inspection template data
func BuildInspectionTemplateKey(templateID uint) string {
	return fmt.Sprintf("%s%d", KeyPrefixInspectionTemplate, templateID)
}

// BuildStatisticsKey builds a cache key for statistics data
func BuildStatisticsKey(statType string) string {
	return fmt.Sprintf("%s%s", KeyPrefixStatistics, statType)
}
