/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2024-07-06 14:44
 * @Description:
 */

package vmodel

type (
	OpenWrt struct {
		HostName   string `json:"host_name"`
		ExtranetIP string `json:"extranet_ip"`
	}

	Stat struct {
		CpuTemp     string `json:"cpu_temp"`
		MemUsage    string `json:"mem_usage"`
		DiskUsage   string `json:"disk_usage"`
		RunningTime string `json:"running_time"`
		BootTime    string `json:"boot_time"`
		NowTime     string `json:"now_time"`
	}

	Clash struct {
		HttpPort  string `json:"http_port"`
		SocksPort string `json:"socks_port"`
		RedirPort string `json:"redir_port"`
		MixedPort string `json:"mixed_port"`
		AllowLan  bool   `json:"allow_lan"`
		Mode      string `json:"mode"`
	}

	// 局域网 dhcp 设备
	Dhcp struct {
		HostName string `json:"host_name"`
		IP       string `json:"ip"`
		Mac      string `json:"mac"`
	}

	// 网络设备
	Device struct {
		Device string `json:"device"`
		IP     string `json:"ip"`
	}
)
