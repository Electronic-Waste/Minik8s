package iptables

import (
	"fmt"
	"strings"
	"net"

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
	DestinationFlag		string = "-d"
	JumpFlag        	string = "-j"
	ProtocolFlag    	string = "-p"
	MatchFlag       	string = "-m"
	MatchParamMark		string = "mark"
	MatchParamComment   string = "comment"
	MatchParamStatistic string = "statistic"
	ModeFlag			string = "--mode"
	ModeParamNth		string = "nth"
	EveryFlag			string = "--every"
	PacketFlag			string = "--packet"
	PacketFlagParam1	string = "0"
	MarkFlag			string = "--mark"
	NotFlag				string = "!"
	CommentFlag     	string = "--comment"
	DestPortFlag    	string = "--dport"
	SourcePortFlag		string = "--sport"
	DNATdestFlag		string = "--to-destination"
	SNATdestFlag		string = "--to-source"
	SetXMarkFlag    	string = "--set-xmark"
	KubeMarkParamPod	string = "0x4000/0x4000"
	KubeMarkParamHost	string = "0x6000/0x6000"
	KubeMarkParamDrop	string = "0x8000/0x8000"

	KubeServicesChainName        string = "KUBE-SERVICES"
	KubePostroutingChainName     string = "KUBE-POSTROUTING"
	KubeHostPostroutingChainName string = "KUBE-HOST-POSTROUTING"
	KubeMarkChainName            string = "KUBE-MARK-MASQ"
	KubeHostMarkChainName        string = "KUBE-HOST-MARK-MASQ"
	KubeServiceChainPrefix       string = "KUBE-SVC-"
	KubePodChainPrefix           string = "KUBE-SEP-"

	AcceptTarget		string = "ACCEPT"
	DropTarget			string = "DROP"
	SnatTarget			string = "SNAT"
	DnatTarget			string = "DNAT"
	MasqTarget			string = "MASQUERADE"
	MarkTarget			string = "MARK"
	ReturnTarget		string = "RETURN"

	ProtocolIPv4		string = "IPv4"
	ProtocolTCP			string = "tcp"
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
	// Creates iptables chains for all minik8s service
	InitServiceIPTables() error
	// Deinit iptabels chains for all minik8s service
	DeinitServiceIPTables() error
	// Creates an iptables chain for one service
	CreateServiceChain() string
	// Adds a rule to KUBE-SERVICES chain, jump to a service chain KUBE_SVC-<serviceChainID>
	ApplyServiceChain(serviceName string, clusterIP string, serviceChainName string, port uint16) error
	// Clears and Deletes an iptables chain for service and rules related to it
	DeleteServiceChain(serviceName string, clusterIP string, serviceChainName string, port uint16) error
	// Clears an iptables chain for service
	ClearServiceChain(serviceName string, serviceChainName string) error
	// Creates an iptables chain for a pod in a service chain
	CreatePodChain() string
	// Adds a jump-to-mark rule and a DNAT rule to a pod chain
	ApplyPodChainRules(podChainName string, podIP string, targetPort uint16) error
	// Inserts a jump-to-pod-chain rule to KUBE-SEP-<podChainID> chain
	// num is the sequence number of this pod in the service, used for round robin.
	ApplyPodChain(serviceName string, serviceChainName string, podName string, podChainName string, num int) error
	// Clears and deletes an iptables chain for pod and rules related to it
	DeletePodChain(podName string, podChainName string) error
}

type IPTablesClient struct {
	HostIP		string
	flannelIP	string
	iptables	*iptables.IPTables
}

func NewIPTablesClient(hostIP, flannelIP string) (*IPTablesClient, error) {
	iptables, err := iptables.New(iptables.IPFamily(iptables.ProtocolIPv4))
	if err != nil {
		return nil, err
	}
	if net.ParseIP(hostIP) == nil {
		return nil, fmt.Errorf("hostIP is invalid\n")
	}
	if net.ParseIP(flannelIP) == nil {
		return nil, fmt.Errorf("flannelIP is invalid\n")
	}
	return &IPTablesClient{
		HostIP: hostIP,
		flannelIP: flannelIP,
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
				"error %v in inserting %s chain to PREROUTING", 
				err, KubeServicesChainName,
			)
		}
	}

	// Check if already exists in OUTPUT
	exist, err = cli.iptables.Exists(
		NATTable,
		OutputChain,
		JumpFlag,
		KubeServicesChainName,
	)
	if err != nil {
		return err
	}

	// If the rule does not exist in OUTPUT chain, insert it
	if !exist {
		err = cli.iptables.Insert(
			NATTable,
			OutputChain,
			1,
			JumpFlag,
			KubeServicesChainName,
		)
		if err != nil {
			return fmt.Errorf(
				"error %v in inserting %s chain to OUTPUT",
				err, KubeServicesChainName,
			)
		}
	}

	// 2. Create KUBE-POSTROUTING chain in nat table
	cli.iptables.NewChain(NATTable, KubePostroutingChainName)

	// Check if already exists in POSTROUTING
	exist, err = cli.iptables.Exists(
		NATTable,
		PostroutingChain,
		JumpFlag,
		KubePostroutingChainName,
	)
	if err != nil {
		return err
	}

	if !exist {
		// If the rule does no exist in POSTROUTING chain, insert it
		err = cli.iptables.Insert(
			NATTable,
			PostroutingChain,
			1,
			JumpFlag,
			KubePostroutingChainName,
		)
		if err != nil {
			return fmt.Errorf(
				"error %v in inserting %s chain to POSTROUTING",
				err, KubePostroutingChainName,
			)
		}

		// -A KUBE-POSTROUTING -m mark --mark 0x6000/0x6000 -j SNAT --to-source <flannelIP>
		err = cli.iptables.AppendUnique(
			NATTable,
			KubePostroutingChainName,
			MatchFlag,
			MatchParamMark,
			MarkFlag,
			KubeMarkParamHost,
			JumpFlag,
			SnatTarget,
			SNATdestFlag,
			cli.flannelIP,
		)
	}

	// 3. Create KUBE-MARK-MASQ
	cli.iptables.NewChain(NATTable, KubeMarkChainName)

	// -A KUBE-MARK-MASQ -j MARK --set-xmark 0x4000/0x4000
	err = cli.iptables.AppendUnique(
		NATTable,
		KubeMarkChainName,
		JumpFlag,
		MarkTarget,
		SetXMarkFlag,
		KubeMarkParamPod,
	)
	if err != nil {
		return fmt.Errorf(
			"error %v in assigning MARK rule to chain %s",
			err, KubeMarkChainName,
		)
	}

	// 4. Create KUBE-HOST-MARK-MASQ
	cli.iptables.NewChain(NATTable, KubeHostMarkChainName)

	// -A KUBE-HOST-MARK-MASQ -j MARK --set-xmark 0x6000/0x6000
	err = cli.iptables.AppendUnique(
		NATTable,
		KubeHostMarkChainName,
		JumpFlag,
		MarkTarget,
		SetXMarkFlag,
		KubeMarkParamHost,
	)
	if err != nil {
		return fmt.Errorf(
			"error %v in assigning MARK rule to chain %s",
			err, KubeHostMarkChainName,
		)
	}
	return nil
}

func (cli *IPTablesClient) DeinitServiceIPTables() error {
	// 1. Delete the rule that jumps to KUBE-SERVICES chain in PREROUTING chain.
	err := cli.iptables.DeleteIfExists(
		NATTable,
		PreroutingChain,
		JumpFlag,
		KubeServicesChainName,
	)
	if err != nil {
		return fmt.Errorf(
			"error %v in deleting the rule `-A %s -j %s`",
			err, PreroutingChain, KubeServicesChainName,
		)
	}

	// 2. Delete the rule that jumps to KUBE-SERVICES chain in OUTPUT chain
	err = cli.iptables.DeleteIfExists(
		NATTable,
		OutputChain,
		JumpFlag,
		KubeServicesChainName,
	)
	if err != nil {
		return fmt.Errorf(
			"error %v in deleting the rule `-A %s -j %s`",
			err, OutputChain, KubeServicesChainName,
		)
	}

	// 3. Delete the rule that jumps to KUBE-POSTROUTING chain in POSTROUTING chain
	err = cli.iptables.DeleteIfExists(
		NATTable,
		PostroutingChain,
		JumpFlag,
		KubePostroutingChainName,
	)
	if err != nil {
		return fmt.Errorf(
			"error %v in deleting the rule `-A %s -j %s`",
			err, PostroutingChain, KubePostroutingChainName,
		)
	}

	// 4. Delete and clear the KUBE-MARK-MASQ chain
	err = cli.iptables.ClearAndDeleteChain(NATTable, KubeMarkChainName)
	if err != nil {
		return fmt.Errorf(
			"error %v in deleting the %s chain",
			err, KubeMarkChainName,
		)
	}

	// 5. Delete and clear the KUBE-HOST-MASQ chain
	err = cli.iptables.ClearAndDeleteChain(NATTable, KubeHostMarkChainName)
	if err != nil {
		return fmt.Errorf(
			"error %v in deleting the %s chain",
			err, KubeHostMarkChainName,
		)
	}

	// 6. Delete and clear the KUBE-POSTROUTING chain
	err = cli.iptables.ClearAndDeleteChain(NATTable, KubePostroutingChainName)
	if err != nil {
		return fmt.Errorf(
			"error %v in deleting the %s chain",
			err, KubePostroutingChainName,
		)
	}

	// 7. Delete and clear the KUBE-SERVICES chain
	err = cli.iptables.ClearAndDeleteChain(NATTable, KubeServicesChainName)
	if err != nil {
		return fmt.Errorf(
			"error %v in deleting the %s chain",
			err, KubeServicesChainName,
		)
	}
	
	return nil
}

func (cli *IPTablesClient) CreateServiceChain() string {
	// Create a chain "KUBE-SVC-<serviceChainID>" in nat table.
	serviceChainID := strings.ToUpper(uuid.New().String()[:8])
	newChainName := KubeServiceChainPrefix + serviceChainID
	cli.iptables.NewChain(NATTable, newChainName)
	return newChainName
}

func (cli *IPTablesClient) ApplyServiceChain(serviceName string, clusterIP string, serviceChainName string, port uint16)  error{
	// Check whether clusterIP is valid or not
	if net.ParseIP(clusterIP) == nil {
		return fmt.Errorf("cluster IP %s is invalid", clusterIP)
	}

	// Adds a rule to KUBE-SERVICES chain, jump to a service chain
	// -A KUBE-SERVICES --p tcp -d <clusterIP> -m comment --comment <serviceName> -dport <port> -j <serviceChainName>
	err := cli.iptables.Insert(
		NATTable,
		KubeServicesChainName,
		1,
		ProtocolFlag,
		ProtocolTCP,
		DestinationFlag,
		clusterIP,
		MatchFlag,
		ProtocolTCP,
		DestPortFlag,
		fmt.Sprint(port),
		MatchFlag,
		MatchParamComment,
		CommentFlag,
		serviceName,
		JumpFlag,
		serviceChainName,
	)
	if err != nil {
		return fmt.Errorf(
			"error %v in adding rule to %s chain",
			err, KubeServicesChainName,
		)
	}

	return nil
}

func (cli *IPTablesClient) CreatePodChain() string {
	// Create a chain "KUBE-SEP-<podChainID>" in nat table
	podChainID := strings.ToUpper(uuid.New().String()[:8])
	newChainName := KubePodChainPrefix + podChainID
	cli.iptables.NewChain(NATTable, newChainName)
	return newChainName
}

func (cli *IPTablesClient) ApplyPodChainRules(podChainName string, podIP string, targetPort uint16, isSameHost bool) error {
	// Check whether podIP is valid or not
	if net.ParseIP(podIP) == nil {
		return fmt.Errorf("pod IP %s in not valid", podIP)
	}

	// Add a rule that jumps to KUBE-MARK-MASQ when the source IP is pod IP
	// -A <podChainName> -s <podIP> -j KUBE-MARK-MASQ
	err := cli.iptables.AppendUnique(
		NATTable,
		podChainName,
		SourceFlag,
		podIP,
		JumpFlag,
		KubeMarkChainName,
	)
	if err != nil {
		return fmt.Errorf(
			"error %v in adding a rule that jumps to %s when the source IP is pod IP",
			err, podIP,
		)
	}

	// Add a rule that jumps to KUBE-HOST-MARK-MASQ when the target pod is on another host
	// -A <podChainName> -s <hostIP> -j KUBE-HOST-MARK-MASQ
	if !isSameHost {
		err = cli.iptables.AppendUnique(
			NATTable,
			podChainName,
			SourceFlag,
			cli.HostIP,
			JumpFlag,
			KubeHostMarkChainName,
		)
		if err != nil {
			return fmt.Errorf(
				"error %v in adding a rule that jumps to %s when the pod in on another host and the source IP is hostIP",
				err, podIP,
			)
		}
	}

	// Add a rule that do DNAT conversion
	// -A <podChainName> -p tcp -m tcp -j DNAT --to-destination <podIP>:<targetPort>
	destination := fmt.Sprintf("%s:%d", podIP, targetPort)
	err = cli.iptables.AppendUnique(
		NATTable,
		podChainName,
		ProtocolFlag,
		ProtocolTCP,
		JumpFlag,
		DnatTarget,
		DNATdestFlag,
		destination,
	)
	if err != nil {
		return fmt.Errorf(
			"error %v in adding DNAT rule for pod IP %s",
			err, destination,
		)
	}

	// Add a SNAT rule to KUBE-POSTROUTING if the pod is mounted on another 

	return nil
}

func (cli *IPTablesClient) ApplyPodChain(serviceName string, serviceChainName string, podName string, podChainName string, num int) error {
	// Inserts a chain jump-to-pod-chain rule to KUBE-SVC-<serviceChainID> chain
	// -A <serviceChainName> -m comment --comment <podName> -m statistic --mode nth --every <num> -j <podChainName>
	err := cli.iptables.Insert(
		NATTable,
		serviceChainName,
		1,
		MatchFlag,
		MatchParamComment,
		CommentFlag,
		podName,
		MatchFlag,
		MatchParamStatistic,
		ModeFlag,
		ModeParamNth,
		EveryFlag,
		fmt.Sprint(num),
		PacketFlag,
		PacketFlagParam1,
		JumpFlag,
		podChainName,
	)
	if err != nil {
		return fmt.Errorf(
			"error %v in adding jump-to-podchain rule for pod %s to service %s",
			err, podName, serviceName,
		)
	}

	return nil
}


// Clears and Deletes an iptables chain for service and rules related to it
func (cli *IPTablesClient) DeleteServiceChain(serviceName string, clusterIP string, serviceChainName string, port uint16) error {
	// Delete the rule that jumps to KUBE-SVC-<serviceChainID> chain in KUBE-SERVICES chain.
	err := cli.iptables.DeleteIfExists(
		NATTable,
		KubeServicesChainName,
		ProtocolFlag,
		ProtocolTCP,
		DestinationFlag,
		clusterIP,
		MatchFlag,
		ProtocolTCP,
		DestPortFlag,
		fmt.Sprint(port),
		MatchFlag,
		MatchParamComment,
		CommentFlag,
		serviceName,
		JumpFlag,
		serviceChainName,
	)
	if err != nil {
		return fmt.Errorf(
			"error %v in deleting iptables rule for service %s",
			err, serviceName,
		)
	}

	// Clear and delete KUBE-SVC-<serviceChainID> chain
	err = cli.iptables.ClearAndDeleteChain(NATTable, serviceChainName)
	if err != nil {
		return fmt.Errorf(
			"error %v in deleting iptables chain for service %s",
			err, serviceName,
		)
	}
	return nil
}

// Clears an iptables chain for service
func (cli *IPTablesClient) ClearServiceChain(serviceName string, serviceChainName string) error {
	// Clear KUBE-SVC-<serviceChainID> chain
	err := cli.iptables.ClearChain(NATTable, serviceChainName)
	if err != nil {
		return fmt.Errorf(
			"error %v in clearing iptables chain for service %s",
			err, serviceName,
		)
	}
	return nil
}

func (cli *IPTablesClient) DeletePodChain(podName string, podChainName string) error {
	// Clear and delete KUBERBOAT-SEP-<podChainID> chain.
	err := cli.iptables.ClearAndDeleteChain(NATTable, podChainName)
	if err != nil {
		return fmt.Errorf(
			"error %v in deleting iptables for pod %s",
			err, podName,
		)
	}
	return nil
}