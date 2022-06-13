package backends

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"

	. "github.com/kpetremann/claw-network/configs"
	"github.com/kpetremann/claw-network/pkg/topology"
)

type RedisRepository struct {
	Topologies []string
	lock       sync.RWMutex
}

func connect() (*redis.Client, *rejson.Handler) {
	addr := fmt.Sprintf("%s:%d", Config.Backends.Redis.Host, Config.Backends.Redis.Port)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: Config.Backends.Redis.Password,
		DB:       Config.Backends.Redis.DB,
	})

	redisJSON := rejson.NewReJSONHandler()
	redisJSON.SetGoRedisClient(redisClient)

	return redisClient, redisJSON
}

func disconnect(redisClient *redis.Client) {
	if err := redisClient.Close(); err != nil {
		fmt.Println("Failed to close Redis connection:", err)
	}
}

// Getter for Topologies
func (r *RedisRepository) GetTopologies() []string {
	r.lock.RLock()
	defer r.lock.RUnlock()

	topologiesCopy := make([]string, len(r.Topologies))
	copy(topologiesCopy, r.Topologies)
	return topologiesCopy
}

func (r *RedisRepository) RefreshRepository() error {
	redisClient, _ := connect()
	defer disconnect(redisClient)

	r.lock.RLock()
	defer r.lock.RUnlock()
	r.Topologies = redisClient.Keys(context.Background(), "*").Val()

	return nil
}

func (r *RedisRepository) LoadTopology(topologyName string) (*topology.Graph, error) {
	redisClient, redisJSON := connect()
	defer disconnect(redisClient)

	// query Redis
	r.lock.RLock()
	result, err := redisJSON.JSONGet(topologyName, ".")
	if err != nil {
		r.lock.RUnlock()
		return nil, err
	}
	r.lock.RUnlock()

	// cast result
	var graph topology.Graph
	if err := json.Unmarshal(result.([]byte), &graph); err != nil {
		return nil, err
	}

	return &graph, nil
}

func (r *RedisRepository) SaveTopology(name string, graph *topology.Graph) error {
	redisClient, redisJSON := connect()
	defer disconnect(redisClient)

	r.lock.Lock()
	res, err := redisJSON.JSONSet(name, ".", graph)
	if err != nil {
		r.lock.Unlock()
		return err
	}
	r.lock.Unlock()

	if res != "OK" {
		return errors.New("Failed to save topology")
	}

	redisClient.BgSave(context.Background())

	// Refresh topology list
	if err := r.RefreshRepository(); err != nil {
		return err
	}

	return nil
}

func (r *RedisRepository) DeleteTopology(topologyName string) error {
	redisClient, redisJSON := connect()
	defer disconnect(redisClient)

	r.lock.Lock()
	res, err := redisJSON.JSONDel(topologyName, ".")
	if err != nil {
		r.lock.Unlock()
		return err
	}
	r.lock.Unlock()

	if res.(int64) != 1 {
		return errors.New("Failed to delete topology")
	}

	redisClient.BgSave(context.Background())

	// Refresh topology list
	if err := r.RefreshRepository(); err != nil {
		return err
	}

	return nil
}

func countTrueFalse(status []bool) (int, int) {
	var trueCount, falseCount int

	for _, status := range status {
		if status {
			trueCount++
		} else {
			falseCount++
		}
	}

	return trueCount, falseCount
}

func extractTopologyStatus(nodesStatus, linksStatus []byte) (map[string]int, error) {
	// cast results
	var nodes []bool
	var links []bool

	if err := json.Unmarshal(nodesStatus, &nodes); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(linksStatus, &links); err != nil {
		return nil, err
	}

	// get statistics
	linksUp, linksDown := countTrueFalse(links)
	nodesUp, nodesDown := countTrueFalse(nodes)

	results := map[string]int{
		"links_up":    linksUp,
		"links_down":  linksDown,
		"links_total": linksUp + linksDown,
		"nodes_up":    nodesUp,
		"nodes_down":  nodesDown,
		"nodes_total": nodesUp + nodesDown,
	}

	return results, nil
}

func (r *RedisRepository) GetTopologyDetails(topologyName string) (map[string]int, error) {
	redisClient, redisJSON := connect()
	defer disconnect(redisClient)

	// query Redis
	r.lock.RLock()
	linksResult, err := redisJSON.JSONGet(topologyName, "$.links..status")
	if err != nil {
		r.lock.Unlock()
		return nil, err
	}

	nodesResult, err := redisJSON.JSONGet(topologyName, "$.nodes..status")
	if err != nil {
		r.lock.Unlock()
		return nil, err
	}
	r.lock.Unlock()

	return extractTopologyStatus(nodesResult.([]byte), linksResult.([]byte))
}

// Get topologies with details such as number of nodes and links up/down
func (r *RedisRepository) ListTopologiesDetails() (map[string]map[string]int, error) {
	redisClient, redisJSON := connect()
	defer disconnect(redisClient)

	// query Redis
	r.lock.RLock()
	linksResult, err := redisJSON.JSONMGet("$.links..status", r.Topologies...)
	if err != nil {
		r.lock.Unlock()
		return nil, err
	}

	nodesResult, err := redisJSON.JSONMGet("$.nodes..status", r.Topologies...)
	if err != nil {
		r.lock.Unlock()
		return nil, err
	}
	r.lock.Unlock()

	linksStatus := linksResult.([]interface{})
	nodesStatus := nodesResult.([]interface{})

	if linksStatus == nil || nodesStatus == nil {
		return nil, errors.New("No data found")
	}

	// compute results
	results := make(map[string]map[string]int)
	for i, topology := range r.Topologies {
		results[topology], err = extractTopologyStatus(nodesStatus[i].([]byte), linksStatus[i].([]byte))
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}
