package path

// This package defines some file paths, such as nginx executable file
// and its conf file. This will make modifying codes in dns and nginx
// package much more easier.

const (
	NginxDir 				= "/usr/local/nginx"
	NginxExecutableFile 	= NginxDir + "/sbin/nginx"
	NginxConfFile 			= NginxDir + "/conf/nginx.conf"
	NginxActionFlag 		= "-s"
	NginxReloadAction 		= "reload"
	NginxQuitAction 		= "quit"
	NginxStopAction 		= "stop"
	NginxIP  				= "127.0.0.1"
	NginxListenPort			= "80"

	HostsFile				= "/etc/hosts"

	ENS3					= "ens3"
	Flannel					= "flannel.1"
)