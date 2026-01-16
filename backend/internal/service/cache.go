package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ems/backend/pkg/redis"
	redisV9 "github.com/redis/go-redis/v9"
)

// CacheService provides high-level caching operations
type CacheService struct{}

// NewCacheService creates a new cache service
func NewCacheService() *CacheService {
	return &CacheService{}
}

// GetJSON retrieves and unmarshals JSON data from cache
func (s *CacheService) GetJSON(key string, dest interface{}) error {
	val, err := redis.Client.Get(redis.Ctx, key).Result()
	if err != nil {
		if err == redisV9.Nil {
			return ErrNotFound
		}
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

// SetJSON marshals and stores data in cache
func (s *CacheService) SetJSON(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return redis.Client.Set(redis.Ctx, key, data, expiration).Err()
}

// Delete deletes multiple keys from cache
func (s *CacheService) Delete(keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return redis.Client.Del(redis.Ctx, keys...).Err()
}

// InvalidateByPattern invalidates cache keys by pattern
func (s *CacheService) InvalidateByPattern(pattern string) error {
	iter := redis.Client.Scan(redis.Ctx, 0, pattern, 0).Iterator()
	for iter.Next(redis.Ctx) {
		if err := redis.Client.Del(redis.Ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

// Cache configuration constants
const (
	// Cache expiration times
	CacheTTLShort  = 5 * time.Minute   // Frequently changing data
	CacheTTLMedium = 15 * time.Minute  // Moderately changing data
	CacheTTLLong   = 1 * time.Hour     // Rarely changing data
	CacheTTLDaily  = 24 * time.Hour    // Static data
)

// Cache key builders for composite keys
func BuildUserListKey(filter interface{}) string {
	return fmt.Sprintf("ems:users:list:%v", filter)
}

func BuildEquipmentListKey(workshopID uint, status string) string {
	if workshopID == 0 && status == "" {
		return "ems:equipment:list:all"
	}
	return fmt.Sprintf("ems:equipment:list:workshop:%d:status:%s", workshopID, status)
}

func BuildRepairOrderListKey(status string, page int) string {
	return fmt.Sprintf("ems:repair:list:status:%s:page:%d", status, page)
}

func BuildInspectionTaskListKey(assignedTo uint, status string, date string) string {
	return fmt.Sprintf("ems:inspection:tasks:user:%d:status:%s:date:%s", assignedTo, status, date)
}

func BuildMaintenanceTaskListKey(assignedTo uint, status string) string {
	return fmt.Sprintf("ems:maintenance:tasks:user:%d:status:%s", assignedTo, status)
}

func BuildSparePartInventoryKey(partID uint) string {
	return fmt.Sprintf("ems:sparepart:inventory:%d", partID)
}

func BuildStatisticsKey(statType string, dateRange string) string {
	return fmt.Sprintf("ems:statistics:%s:%s", statType, dateRange)
}

// Equipment cache invalidation helpers
func InvalidateEquipmentCache(equipmentID uint) error {
	keys := []string{
		redis.BuildEquipmentKey(equipmentID),
	}
	return redis.Client.Del(redis.Ctx, keys...).Err()
}

// Workshop cache invalidation helpers
func InvalidateWorkshopCache(workshopID uint) error {
	keys := []string{
		redis.BuildWorkshopKey(workshopID),
	}
	return redis.Client.Del(redis.Ctx, keys...).Err()
}

// Factory cache invalidation helpers
func InvalidateFactoryCache(factoryID uint) error {
	keys := []string{
		redis.BuildFactoryKey(factoryID),
	}
	return redis.Client.Del(redis.Ctx, keys...).Err()
}

// Equipment type cache invalidation helpers
func InvalidateEquipmentTypeCache(typeID uint) error {
	keys := []string{
		redis.BuildEquipmentTypeKey(typeID),
	}
	return redis.Client.Del(redis.Ctx, keys...).Err()
}

// User cache invalidation helpers
func InvalidateUserCache(userID uint) error {
	keys := []string{
		redis.BuildUserKey(userID),
	}
	// Invalidate user-related lists
	iter := redis.Client.Scan(redis.Ctx, 0, "ems:users:list:*", 0).Iterator()
	for iter.Next(redis.Ctx) {
		redis.Client.Del(redis.Ctx, iter.Val())
	}
	return redis.Client.Del(redis.Ctx, keys...).Err()
}

// Inspection template cache invalidation
func InvalidateInspectionTemplateCache(templateID uint) error {
	keys := []string{
		redis.BuildInspectionTemplateKey(templateID),
	}
	return redis.Client.Del(redis.Ctx, keys...).Err()
}

// Statistics cache invalidation
func InvalidateStatisticsCache(statType string) error {
	pattern := fmt.Sprintf("ems:statistics:%s:*", statType)
	iter := redis.Client.Scan(redis.Ctx, 0, pattern, 0).Iterator()
	for iter.Next(redis.Ctx) {
		redis.Client.Del(redis.Ctx, iter.Val())
	}
	return iter.Err()
}

// GetOrSetJSON is a helper that gets from cache or sets from function
func (s *CacheService) GetOrSetJSON(key string, dest interface{}, expiration time.Duration, fn func() (interface{}, error)) error {
	// Try to get from cache
	err := s.GetJSON(key, dest)
	if err == nil {
		return nil // Cache hit
	}
	if err != ErrNotFound {
		return err // Cache error
	}

	// Cache miss, call function
	data, err := fn()
	if err != nil {
		return err
	}

	// Set to cache
	if err := s.SetJSON(key, data, expiration); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: failed to set cache: %v\n", err)
	}

	// Unmarshal to destination
	dataBytes, _ := json.Marshal(data)
	return json.Unmarshal(dataBytes, dest)
}
