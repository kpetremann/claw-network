package backends

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"

	. "github.com/kpetremann/claw-network/configs"
	"github.com/kpetremann/claw-network/pkg/topology"
)

type RedisRepository struct {
	Topologies []string
}

// Getter for Topologies
func (r *RedisRepository) GetTopologies() []string {
	return r.Topologies
}

// Setter for Topologies
func (r *RedisRepository) SetTopologies(topologies []string) {
	r.Topologies = topologies
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

func (r *RedisRepository) RefreshRepository() error {
	redisClient, _ := connect()
	defer disconnect(redisClient)

	r.Topologies = redisClient.Keys(context.Background(), "*").Val()

	return nil
}

func (r *RedisRepository) SaveTopology(name string, graph *topology.Graph) error {
	redisClient, redisJSON := connect()
	defer disconnect(redisClient)

	res, err := redisJSON.JSONSet(name, ".", graph)
	if err != nil {
		return err
	}

	if res != "OK" {
		return errors.New("Failed to save topology")
	}

	redisClient.BgSave(context.Background())

	return nil
}

func (r *RedisRepository) LoadTopology(topologyName string) (*topology.Graph, error) {
	redisClient, redisJSON := connect()
	defer disconnect(redisClient)

	// query Redis
	result, err := redisJSON.JSONGet(topologyName, ".")
	if err != nil {
		return nil, err
	}

	// cast result
	var graph topology.Graph
	if err := json.Unmarshal(result.([]byte), &graph); err != nil {
		return nil, err
	}

	return &graph, nil
}

func (r *RedisRepository) DeleteTopology(topologyName string) error {
	redisClient, redisJSON := connect()
	defer disconnect(redisClient)

	res, err := redisJSON.JSONDel(topologyName, ".")
	if err != nil {
		return err
	}

	if res.(int64) != 1 {
		return errors.New("Failed to delete topology")
	}

	redisClient.BgSave(context.Background())

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

// Get topologies with details such as number of nodes and links up/down
func (r *RedisRepository) ListTopologiesDetail() (map[string]map[string]int, error) {
	redisClient, redisJSON := connect()
	defer disconnect(redisClient)

	// query Redis
	linksResult, err := redisJSON.JSONMGet("$.links..status", r.Topologies...)
	if err != nil {
		return nil, err
	}

	nodesResult, err := redisJSON.JSONMGet("$.nodes..status", r.Topologies...)
	if err != nil {
		return nil, err
	}

	linksStatus := linksResult.([]interface{})
	nodesStatus := nodesResult.([]interface{})

	if linksStatus == nil || nodesStatus == nil {
		return nil, errors.New("No data found")
	}

	// compute results
	results := make(map[string]map[string]int)
	for i, topology := range r.Topologies {
		// cast results
		var links []bool
		var nodes []bool

		if err := json.Unmarshal(linksStatus[i].([]byte), &links); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(nodesStatus[i].([]byte), &nodes); err != nil {
			return nil, err
		}

		// get statistics
		linksUp, linksDown := countTrueFalse(links)
		nodesUp, nodesDown := countTrueFalse(nodes)

		results[topology] = map[string]int{
			"links_up":    linksUp,
			"links_down":  linksDown,
			"links_total": linksUp + linksDown,
			"nodes_up":    nodesUp,
			"nodes_down":  nodesDown,
			"nodes_total": nodesUp + nodesDown,
		}
	}

	return results, nil
}
