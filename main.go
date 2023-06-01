package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.New()
)

func main() {
	setupLogging()
	m := setupMetrics()
	c := setupConfig()

	// start metrics checks
	go check(m, c)

	// serve metrics endpoint
	serve(c)
}

func check(m *metrics, c *config) {
	for {
		// create a context with a timeout to stop hanging database timeouts.
		_, cancel := context.WithTimeout(context.Background(), c.timeout)

		// Get block height
		networkBlockHeight, err := getFSBlockHeight(c)
		if err != nil {
			// set measurement to 0 on error. Since we're not getting the block details
			// reset the sigs count and details.
			m.fsNetworkHeight.WithLabelValues(c.mcNetwork).Set(0)
			cancel()
			time.Sleep(c.checkWait)
			continue
		}

		// set the block_height metric.
		log.Debugf("Set Block Height Metric to: %.0f", networkBlockHeight)
		m.fsNetworkHeight.WithLabelValues(c.mcNetwork).Set(networkBlockHeight)

		// cancel context timeout and sleep until next round
		cancel()
		time.Sleep(c.checkWait)
	}
}

func getFSBlockHeight(c *config) (float64, error) {
	request := &FsRequest{
		ID:      1,
		JsonRPC: "2.0",
		Method:  "get_network_status",
	}
	result := &FsGetNetworkStatusReturn{}

	// query full-service
	err := getFS(c, request, result)
	if err != nil {
		return 0, err
	}

	// convert the string return to a float
	networkBlockHeight, err := strconv.ParseFloat(result.Result.NetworkStatus.NetworkBlockHeight, 64)
	if err != nil {
		log.Errorf("Unable to parse block height: %v", err)
		return 0, err
	}

	return networkBlockHeight, nil
}

// give it a request pointer and a result struct.
// The body returned will be unmarshaled into to the result struct.
func getFS(c *config, request *FsRequest, result interface{}) error {
	payload, err := json.Marshal(request)
	if err != nil {
		log.Fatalf("unable to marshal payload: %v", err)
	}
	log.Debugf("Request Payload: %v", string(payload))

	// send the request
	resp, err := http.Post(c.url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Errorf("Post failed: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		// Read the body response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("Unable to read response body: %v", err)
			return err
		}

		log.Debugf("Return: %s", string(body))
		// Unmarshal the json into result object
		err = json.Unmarshal(body, result)
		if err != nil {
			log.Errorf("Unable to unmarshal full-service response: %v", err)
			return err
		}
	} else {
		log.Errorf("HTTP Response Error %d", resp.StatusCode)
		return err
	}

	return nil
}

// Get configuration from env vars
func setupConfig() *config {
	listenHost := getEnv("LISTEN_HOST", "127.0.0.1")

	listenPort, err := strconv.Atoi(getEnv("LISTEN_PORT", "9090"))
	if err != nil {
		log.Fatalf("Unable to parse LISTEN_PORT: %v", err)
	}

	checkWait, err := strconv.Atoi(getEnv("CHECK_WAIT", "30"))
	if err != nil {
		log.Fatalf("Unable to parse CHECK_WAIT: %v", err)
	}

	timeout, err := strconv.Atoi(getEnv("TIMEOUT", "10"))
	if err != nil {
		log.Fatalf("Unable to parse TIMEOUT: %v", err)
	}

	// Number of nodes in the network
	numOfNodes, err := strconv.Atoi(getEnv("NUM_OF_NODES", "10"))
	if err != nil {
		log.Fatalf("Unable to parse NUM_OF_NODES: %v", err)
	}

	return &config{
		url:        getEnv("FULL_SERVICE_URL", "https://readonly-fs-mainnet.mobilecoin.com/wallet/v2"),
		mcNetwork:  getEnv("MC_NETWORK", "main"),
		checkWait:  time.Duration(checkWait) * time.Second,
		timeout:    time.Duration(timeout) * time.Second,
		listenHost: listenHost,
		listenPort: listenPort,
		numOfNodes: numOfNodes,
	}
}

func setupMetrics() *metrics {
	m := &metrics{}

	// this sets the metric name mc_network_block_height with a label network
	m.fsNetworkHeight = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "mc",
		Subsystem: "network",
		Name:      "block_height",
		Help:      "Latest block height",
	}, []string{"network"})

	return m
}

func serve(c *config) {
	log.Infof("Serving metrics at %s:%d/metrics", c.listenHost, c.listenPort)
	http.Handle("/metrics", promhttp.Handler())

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", c.listenHost, c.listenPort), nil)
	if err != nil {
		log.Fatalf("error serving http: %v", err)
	}
}
