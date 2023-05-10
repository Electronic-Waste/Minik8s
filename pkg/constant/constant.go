package constant

const Cri_uri = "unix:///var/run/containerd/containerd.sock"
const Cli_uri = "/run/containerd/containerd.sock"
const SandBox_Image = "registry.aliyuncs.com/google_containers/pause:3.9"

// for the reason we use nerdctl tools to set up the pause container, so we need to init this path
const NerdctlDataRoot = "/var/lib/nerdctl"
