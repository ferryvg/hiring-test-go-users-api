package db

import (
	"strconv"
	"strings"
	"sync"

	"github.com/ferryvg/hiring-test-go-users-api/internal/config"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// ConnList represents list of opened database connections.
type ConnList interface {

	// SetNodes sets list of server nodes.
	SetNodes(nodes []string)

	// Connections returns list opened connections.
	Connections() []*sqlx.DB

	// Connection returns connection for specified node.
	Connection(node string) (*sqlx.DB, bool)

	// Close closes all opened connections.
	Close()
}

// ConnFactory represents database connection factory.
type ConnFactory interface {

	// Create creates connection for specified node.
	Create(node string) (*sqlx.DB, error)
}

type connListImpl struct {
	factory ConnFactory
	target  *config.ClusterNodeList
	logger  logrus.FieldLogger
	nodes   []string
	connMap map[string]*sqlx.DB
	mu      *sync.RWMutex
}

// NewConnList creates new database connections list.
func NewConnList(factory ConnFactory, target *config.ClusterNodeList, logger logrus.FieldLogger) ConnList {
	return &connListImpl{
		factory: factory,
		target:  target,
		logger:  logger,
		connMap: make(map[string]*sqlx.DB),
		mu:      new(sync.RWMutex),
	}
}

func (cl *connListImpl) SetNodes(nodes []string) {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	nodesMap, nodesList := cl.createConns(nodes)

	cl.closeUnusedConns(nodesMap)

	cl.updateTarget(nodesList)

	cl.nodes = nodesList
}

func (cl *connListImpl) Connections() []*sqlx.DB {
	cl.mu.RLock()
	defer cl.mu.RUnlock()

	cList := make([]*sqlx.DB, 0, len(cl.nodes))

	for _, node := range cl.nodes {
		if conn, exist := cl.connMap[node]; exist {
			cList = append(cList, conn)
		}
	}

	return cList
}

func (cl *connListImpl) Connection(node string) (*sqlx.DB, bool) {
	cl.mu.RLock()
	defer cl.mu.RUnlock()

	conn, exists := cl.connMap[node]

	return conn, exists
}

func (cl *connListImpl) Close() {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	for node, conn := range cl.connMap {
		cl.closeConn(node, conn)
	}

	cl.connMap = make(map[string]*sqlx.DB)
	cl.nodes = nil
}

func (cl *connListImpl) closeUnusedConns(nodesMap map[string]struct{}) {
	for node, conn := range cl.connMap {
		if _, exist := nodesMap[node]; !exist {
			cl.closeConn(node, conn)
			delete(cl.connMap, node)
		}
	}
}

func (cl *connListImpl) createConns(nodes []string) (map[string]struct{}, []string) {
	nodesMap := make(map[string]struct{})
	nodesList := make([]string, 0, len(nodes))

	for _, node := range nodes {
		nodesMap[node] = struct{}{}

		if _, exist := cl.connMap[node]; !exist {
			conn, err := cl.createConn(node)
			if err != nil {
				continue
			}

			cl.connMap[node] = conn
		}

		nodesList = append(nodesList, node)
	}

	return nodesMap, nodesList
}

// Creates connection for specified node.
func (cl *connListImpl) createConn(node string) (*sqlx.DB, error) {
	conn, err := cl.factory.Create(node)
	if err != nil {
		cl.logger.WithError(err).WithField("node", node).Error("Failed to create database connection")
	}

	return conn, err
}

func (cl *connListImpl) closeConn(node string, conn *sqlx.DB) {
	if err := conn.Close(); err != nil {
		cl.logger.WithField("node", node).WithError(err).Error("Failed to close database connection")
	}
}

func (cl *connListImpl) updateTarget(nodesList []string) {
	nodes := make([]*config.ClusterNode, 0, len(nodesList))

	for _, node := range nodesList {
		details := strings.SplitN(node, ":", 2)

		if len(details) == 0 {
			continue
		}

		host := details[0]
		var port int

		if len(details) == 2 {
			portN, err := strconv.ParseInt(details[1], 10, 32)
			if err != nil {
				cl.logger.WithError(err).WithField("node", node).Warn("Failed to split host/port of node")
			} else {
				port = int(portN)
			}
		}

		nodes = append(nodes, config.NewClusterNode(host, port))
	}

	cl.target.Set(nodes)
}
