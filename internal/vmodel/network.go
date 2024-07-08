/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2024-07-08 16:56
 * @Description:
 */

package vmodel

type (
	Extranet struct {
		IP       string `json:"ip"`       // ip
		Location string `json:"location"` // 位置
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

	NetWork struct {
		HostName string    `json:"host_name"`
		Wan      *Extranet `json:"wan"`
		Lan      []*Device `json:"lan"`
		Devices  []*Dhcp   `json:"devices"`
	}
)
