package numa

import (
	"io/ioutil"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/onmetal/inventory/pkg/mem"
)

const (
	CNodeCPUListPath  = "/cpulist"
	CNodeDistancePath = "/distance"
	CNodeMemInfo      = "/meminfo"

	CDistanceTrimPattern = "\\D+"
	CCPUListTrimPattern  = "[^0-9\\-,]"
)

var CDistanceTrimRegexp = regexp.MustCompile(CDistanceTrimPattern)
var CCPUListTrimRegexp = regexp.MustCompile(CCPUListTrimPattern)

type NodeSvc struct {
	memInfoSvc *mem.InfoSvc
	statSvc    *StatSvc
}

func NewNodeSvc(memInfoSvc *mem.InfoSvc, statSvc *StatSvc) *NodeSvc {
	return &NodeSvc{
		memInfoSvc: memInfoSvc,
		statSvc:    statSvc,
	}
}

func (s *NodeSvc) GetNode(thePath string, nodeId int) (*Node, error) {
	distancePath := path.Join(thePath, CNodeDistancePath)
	distanceData, err := ioutil.ReadFile(distancePath)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read distance file from path %s", distancePath)
	}

	distanceString := string(distanceData)
	distanceStrings := strings.Fields(distanceString)
	distances := make([]int, len(distanceStrings))
	for i := range distances {
		distance, err := strconv.Atoi(distanceStrings[i])
		if err != nil {
			return nil, errors.Wrapf(err, "unable to convert distance string %s (from %s) to int", distanceStrings[i], distanceString)
		}
		distances[i] = distance
	}

	cpuListPath := path.Join(thePath, CNodeCPUListPath)
	cpuListData, err := ioutil.ReadFile(cpuListPath)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read distance file from path %s", cpuListPath)
	}

	cpuListString := string(cpuListData)
	cpuListTrimmedString := CCPUListTrimRegexp.ReplaceAllString(cpuListString, "")
	cpuListElements := strings.Split(cpuListTrimmedString, ",")

	cpuList := make([]int, 0)

	// NUMA CPU list looks like 0,3,5-8,9,11-15
	for _, element := range cpuListElements {
		if cpuId, err := strconv.Atoi(element); err == nil {
			cpuList = append(cpuList, cpuId)
			continue
		}

		cpuRange := strings.Split(element, "-")

		if len(cpuRange) != 2 {
			return nil, errors.Errorf("expected to have a NUMA CPU range, but got %s", element)
		}

		first, err := strconv.Atoi(cpuRange[0])
		if err != nil {
			return nil, errors.Errorf("expected to have a number in NUMA CPU range, but got %s", cpuRange[0])
		}
		last, err := strconv.Atoi(cpuRange[1])
		if err != nil {
			return nil, errors.Errorf("expected to have a number in NUMA CPU range, but got %s", cpuRange[1])
		}

		for i := first; i <= last; i++ {
			cpuList = append(cpuList, i)
		}
	}

	memPath := path.Join(thePath, CNodeMemInfo)
	memInfo, err := s.memInfoSvc.GetInfoFromFile(memPath)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to obtain meminfo for %s", thePath)
	}

	stat, err := s.statSvc.GetStat(thePath)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to obtain stat for %s", thePath)
	}

	return &Node{
		ID:        nodeId,
		Distances: distances,
		CPUs:      cpuList,
		Memory:    memInfo,
		Stat:      stat,
	}, nil
}
