package envoy

import (
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
)

// Envoy manages a running Envoy proxy instance via the
// ListenerDiscoveryService and RouteDiscoveryService gRPC APIs.
type Envoy struct {
	cmd     *exec.Cmd
	LogPath string
	ldsSock string
	lds     *LDSServer
	rdsSock string
	rds     *RDSServer
}

func createConfig(filePath string, adminAddress string) {
	config := string("{\n" +
		"  \"listeners\": [],\n" +
		"  \"admin\": { \"access_log_path\": \"/dev/null\",\n" +
		"             \"address\": \"tcp://" + adminAddress + "\" },\n" +
		"  \"cluster_manager\": {\n" +
		"    \"clusters\": []\n" +
		"  }\n" +
		"}\n")

	log.Print("Config: ", config)
	err := ioutil.WriteFile(filePath, []byte(config), 0644)
	if err != nil {
		panic(err)
	}
}

// StartEnvoy starts an Envoy proxy instance. If 'debug' is true, an
// debug version of the Envoy binary is started with the log level
// 'debug', otherwise a production version is started at the default
// log level.
func StartEnvoy(debug bool, adminPort int, stateDir string, logDir string) *Envoy {
	bootstrapPath := stateDir + "bootstrap.pb"
	configPath := stateDir + "envoy-config.json"
	logPath := logDir + "cilium-envoy.log"
	adminAddress := "127.0.0.1:" + strconv.Itoa(adminPort)
	ldsPath := stateDir + "lds.sock"
	rdsPath := stateDir + "rds.sock"
	e := Envoy{LogPath: logPath, ldsSock: ldsPath, rdsSock: rdsPath}

	// Create configuration
	createBootstrap(bootstrapPath, "envoy1", "cluster1", "version1",
		"ldsCluster", ldsPath, "rdsCluster", rdsPath, "cluster1")
	createConfig(configPath, adminAddress)

	if debug {
		e.cmd = exec.Command("sh", "-c", "cilium-envoy-debug >"+logPath+" 2>&1 -l debug -c "+configPath+" -b "+bootstrapPath)
	} else {
		e.cmd = exec.Command("sh", "-c", "cilium-envoy >"+logPath+" 2>&1 -c "+configPath+" -b "+bootstrapPath)
	}

	e.lds = createLDSServer(ldsPath)
	e.rds = createRDSServer(rdsPath, e.lds)
	e.rds.run()
	e.lds.run(e.rds)

	err := e.cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Envoy process started at pid ", e.cmd.Process.Pid)
	return &e
}

// StopEnvoy kills the Envoy process started with StartEnvoy. The gRPC API streams are terminated
// first.
func (e *Envoy) StopEnvoy() {
	log.Print("Stopping Envoy process ", e.cmd.Process.Pid)
	e.rds.stop()
	e.lds.stop()
	err := e.cmd.Process.Kill()
	if err != nil {
		log.Fatal(err)
	}
	e.cmd.Wait()
}

// AddListener adds a listener to a running Envoy proxy.
func (e *Envoy) AddListener(name string, port uint32, l7rules []AuxRule) {
	e.lds.addListener(name, port, l7rules)
}

// UpdateListener changes to the L7 rules of an existing Envoy Listener.
func (e *Envoy) UpdateListener(name string, l7rules []AuxRule) {
	e.lds.updateListener(name, l7rules)
}

// RemoveListener removes an existing Envoy Listener.
func (e *Envoy) RemoveListener(name string) {
	e.lds.removeListener(name)
}
