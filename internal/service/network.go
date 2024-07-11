/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2024-07-06 11:30
 * @Description:
 */

package service

import (
	"github.com/leafney/rose"
	"github.com/leafney/whisky/global"
	"github.com/leafney/whisky/internal/vmodel"
	"github.com/leafney/whisky/pkgs/cmds"
	"github.com/leafney/whisky/utils"
	"regexp"
	"strings"
)

func NetWorkInfo() (res *vmodel.NetWork) {
	// 初始化
	res = &vmodel.NetWork{
		Wan:     &vmodel.Extranet{},
		Lan:     make([]*vmodel.Device, 0),
		Devices: make([]*vmodel.Dhcp, 0),
	}

	res.Wan, _ = NetWorkExtranetIP()
	res.HostName = NetWorkHostName()
	res.Lan, _ = NetWorkIntranetIP()

	//
	res.Devices, _ = NetWorkDhcp()

	return
}

// NetWorkIntranetIP 查询内网 IP
func NetWorkIntranetIP() (res []*vmodel.Device, err error) {
	res = make([]*vmodel.Device, 0)
	deviceInfo, err := utils.RunBash(cmds.ScriptNetworkDevice)
	if err != nil {
		deviceInfo = ""
	}

	// 对返回的结果去除前后空行
	deviceInfo = rose.StrTrim(deviceInfo)
	global.GXLog.Debugf("获取到的设备 [%v]", deviceInfo)

	devices := strings.Split(deviceInfo, "\n")
	for _, v := range devices {
		global.GXLog.Debugf("获取到的设备222 [%v]", v)
		vv := strings.Split(v, "#")
		if len(vv) == 2 {
			res = append(res, &vmodel.Device{
				Device: vv[0],
				IP:     vv[1],
			})
		}
	}

	return
}

// NetWorkExtranetIP 查询外网 IP
func NetWorkExtranetIP() (res *vmodel.Extranet, err error) {
	//res = new(vmodel.Extranet)
	res = &vmodel.Extranet{}

	ipInfo, err := utils.RunBash(cmds.ScriptMyIP)
	if err != nil {
		ipInfo = ""
	}

	if !rose.StrIsEmpty(ipInfo) {
		global.GXLog.Debugf("ip 结果[%v]", ipInfo)

		//	当前 IP：123.121.56.153  来自于：中国 北京 北京  联通
		// 使用正则表达式提取IP地址
		ipRegex := regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)
		tIP := ipRegex.FindString(ipInfo)
		res.IP = rose.StrTrim(tIP)

		// 使用正则表达式提取地理位置
		locationRegex := regexp.MustCompile(`来自于：(.*)`)
		temps := locationRegex.FindStringSubmatch(ipInfo)
		tLoc := ""
		if len(temps) >= 2 {
			tLoc = temps[1]
		}
		res.Location = rose.StrTrim(tLoc)

		global.GXLog.Debugf("提取后结果 [%v]-[%v]", tIP, tLoc)
	}

	return
}

func NetWorkHostName() string {
	hostInfo, err := utils.RunBash(cmds.ScriptHostName)
	if err != nil {
		hostInfo = ""
	}

	return hostInfo
}

func NetWorkDhcp() (res []*vmodel.Dhcp, err error) {
	res = make([]*vmodel.Dhcp, 0)
	dhcpInfo, err := utils.RunBash(cmds.ScriptNetWorkDhcp)
	if err != nil {
		dhcpInfo = ""
	}
	if !rose.StrIsEmpty(dhcpInfo) {
		global.GXLog.Debugf("返回的 dhcp [%v]", dhcpInfo)
		dhcps := strings.Split(dhcpInfo, "\n")
		for _, cp := range dhcps {
			cps := rose.StrAnySplit(cp, " ")
			if len(cps) == 5 {
				res = append(res, &vmodel.Dhcp{
					Mac:      rose.StrToUpper(cps[4]),
					IP:       cps[2],
					HostName: cps[3],
				})
			}
		}
	}

	return
}
