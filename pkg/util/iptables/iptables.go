package iptables



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

}