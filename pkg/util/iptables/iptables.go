package iptables

import (
	"fmt"

	"github.com/coreos/go-iptables/iptables"
	"github.com/google/uuid"
)

const (
	PrependFlag 		string = "-I"
	AppendFlag 			string = "-A"
	CreateChainFlag 	string = "-N"
	FlushChainFlag  	string = "-F"
	DeleteChainFlag 	string = "-X"
	ListChainFlag   	string = "-S"
	CheckRuleFlag   	string = "-C"
	DeleteRuleFlag  	string = "-D"
	SourceFlag      	string = "-s"
	JumpFlag        	string = "-j"
	ProtocolFlag    	string = "-p"
	MatchFlag       	string = "-m"
	MatchParamComment   string = "comment"
	CommentFlag     	string = "--comment"
	DestPortFlag    	string = "--dport"
	SourcePortFlag		string = "--sport"
	DNATdestFlag		string = "--to-destination"
	SNATdestFlag		string = "--to-source"
	RandomFullyFlag 	string = "--random-fully"
	SetXMarkFlag    	string = "--set-xmark"

	KubeServicesChainName        string = "KUBE-SERVICES"
	KubePostroutingChainName     string = "KUBE-POSTROUTING"
	KubeHostPostroutingChainName string = "KUBE-HOST-POSTROUTING"
	KubeMarkChainName            string = "KUBE-MARK-MASQ"
	KubeHostMarkChainName        string = "KUBE-HOST-MARK-MASQ"
	KubeServiceChainPrefix       string = "KUBE-SVC-"
	KubePodChainPrefix           string = "KUBE-SEP-"

	ProtocolIPv4		string = "IPv4"
	NATTable 			string = "nat"
	FilterTable			string = "filter"
	MangleTable			string = "mangle"
	RawTable			string = "raw"
	PostroutingChain 	string = "POSTROUTING"
	PreroutingChain 	string = "PREROUTING"
	OutputChain 		string = "OUTPUT"
	InputChain 			string = "INPUT"
	ForwardChain 		string = "FORWARD"
	IPTablesSaveCmd		string = "iptables-save"
	IPTablesRestoreCmd  string = "iptables-restore"
	IPTablesCmd         string = "iptables"

)

type Inteface interface {
	InitServiceIPTables() error
	CreateServiceChain() string
	ApplyServiceChain(serviceName string, clusterIP string, serviceChainName string, port uin16) error
	DeleteServiceChain(serviceName string, clusterIP string, serviceChainName string, port uint16) error
	FlushServiceChain(serviceName string, serviceChainName string) error
}

type IPTablesClient struct {
	hostIP		string
	iptables	*iptables.IPTables
}

func NewIPTablesClient(hostIP string) (*IPTablesClient, error) {
	iptables, err := iptables.New(iptables.IPFamily(iptables.ProtocolIPv4))
	if err != nil {
		return nil, err
	}
	return &IPTablesClient{
		hostIP: hostIP,
		iptables: iptables,
	}, nil
}

func (cli *IPTablesClient) InitServiceIPTables() error {
	// 1. Create KUBE-SERVICES chain in nat table
	
	// Create the chain
	cli.iptables.NewChain(NATTable, KubeServicesChainName)

	// Check if already exists in PREROUTING
	exist, err := cli.iptables.Exists(
		NATTable,
		PreroutingChain,
		JumpFlag,
		KubeServicesChainName,
	)
	if err != nil {
		return err
	}

	// If the rule does not exist in PREROUTING chain, insert it
	if !exist {
		err = cli.iptables.Insert(
			NATTable,
			PreroutingChain,
			1,
			JumpFlag,
			KubeServicesChainName,
		)
		if err != nil {
			return fmt.Errorf(
				"error %s in inserting %s chain", 
				err, KubeServicesChainName,
			)
		}
	}
}