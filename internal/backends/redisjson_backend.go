package backends

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"

	"github.com/kpetremann/claw-network/pkg/topology"
)

type RedisRepository struct {
	Topologies  []string
	redisClient *redis.Client
	redisJSON   *rejson.Handler
}

func (r *RedisRepository) Connect() {
	r.redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	r.redisJSON = rejson.NewReJSONHandler()
	r.redisJSON.SetGoRedisClient(r.redisClient)
}

func (r *RedisRepository) Disconnect() {
	if err := r.redisClient.Close(); err != nil {
		fmt.Println("Failed to close Redis connection:", err)
	}
}

func (r *RedisRepository) RefreshRepository() error {
	r.Connect()
	defer r.Disconnect()

	r.Topologies = r.redisClient.Keys(context.Background(), "*").Val()

	return nil
}

func (r *RedisRepository) SaveTopology(name string, graph *topology.Graph) error {
	r.Connect()
	defer r.Disconnect()

	res, err := r.redisJSON.JSONSet(name, ".", graph)
	if err != nil {
		return err
	}

	if res != "OK" {
		return errors.New("Failed to save topology")
	}

	return nil
}

func (r *RedisRepository) LoadTopology(topologyName string) (*topology.Graph, error) {
	r.Connect()
	defer r.Disconnect()

	// query Redis
	result, err := r.redisJSON.JSONGet(topologyName, ".")
	if err != nil {
		return nil, err
	}

	// cast result
	var graph topology.Graph
	json.Unmarshal(result.([]byte), &graph)

	return &graph, nil
}

func (r *RedisRepository) DeleteTopology(topologyName string) error {
	r.Connect()
	defer r.Disconnect()

	res, err := r.redisJSON.JSONDel(topologyName, ".")
	if err != nil {
		return err
	}

	if res.(int64) != 1 {
		return errors.New("Failed to delete topology")
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

// Get topologies with details such as number of nodes and links up/down
func (r *RedisRepository) ListTopologiesDetail() (map[string]map[string]int, error) {
	r.Connect()
	defer r.Disconnect()

	// query Redis
	linksResult, err := r.redisJSON.JSONMGet("$.links..status", r.Topologies...)
	if err != nil {
		return nil, err
	}

	nodesResult, err := r.redisJSON.JSONMGet("$.nodes..status", r.Topologies...)
	if err != nil {
		return nil, err
	}

	linksStatus := linksResult.([]interface{})
	nodesStatus := nodesResult.([]interface{})

	if linksStatus == nil || nodesStatus == nil {
		return nil, fmt.Errorf("No data found")
	}

	// compute results
	results := make(map[string]map[string]int)
	for i, topology := range r.Topologies {
		// cast results
		var links []bool
		var nodes []bool

		json.Unmarshal(linksStatus[i].([]byte), &links)
		json.Unmarshal(nodesStatus[i].([]byte), &nodes)

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
