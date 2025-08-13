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
	"github.com/leafney/whisky/internal/vmodel"
	"github.com/leafney/whisky/pkg/cmds"
	"github.com/leafney/whisky/pkg/utils"
	"github.com/leafney/whisky/pkg/xlogx"

	"regexp"
	"strings"
)

type NetWork struct {
	XLog *xlogx.XLogSvc
}

func (s *NetWork) NetWorkInfo() (res *vmodel.NetWork, err error) {
	// 初始化
	res = &vmodel.NetWork{
		Wan:     &vmodel.Extranet{},
		Lan:     make([]*vmodel.Device, 0),
		Devices: make([]*vmodel.Dhcp, 0),
	}

	res.Wan, err = s.NetWorkExtranetIP()
	if err != nil {
		s.XLog.Errorf("[NetWorkExtranetIP] 执行失败 [%v]", err)
	}

	res.HostName, err = s.NetWorkHostName()
	if err != nil {
		s.XLog.Errorf("[NetWorkHostName] 执行失败 [%v]", err)
	}

	res.Lan, err = s.NetWorkIntranetIP()
	if err != nil {
		s.XLog.Errorf("[NetWorkIntranetIP] 执行失败 [%v]", err)
	}

	//
	res.Devices, err = s.NetWorkDhcp()
	if err != nil {
		s.XLog.Errorf("[NetWorkDhcp] 执行失败 [%v]", err)
	}

	return
}

// NetWorkIntranetIP 查询内网 IP
func (s *NetWork) NetWorkIntranetIP() (res []*vmodel.Device, err error) {
	res = make([]*vmodel.Device, 0)
	deviceInfo, err := utils.RunBash(cmds.ScriptNetworkDevice)
	if err != nil {
		s.XLog.Errorf("shell 脚本 [ScriptNetworkDevice] 执行失败 [%v]", err)
		return nil, err
	}

	// 对返回的结果去除前后空行
	deviceInfo = rose.StrTrim(deviceInfo)
	//
	devices := strings.Split(deviceInfo, "\n")
	for _, v := range devices {
		s.XLog.Debugf("获取到的设备 [%v]", v)
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
func (s *NetWork) NetWorkExtranetIP() (res *vmodel.Extranet, err error) {
	//res = new(vmodel.Extranet)
	res = &vmodel.Extranet{}

	ipInfo, err := utils.RunBash(cmds.ScriptMyIP)
	if err != nil {
		s.XLog.Errorf("shell 脚本 [ScriptMyIP] 执行失败 [%v]", err)
		return nil, err
	}

	if !rose.StrIsEmpty(ipInfo) {
		s.XLog.Debugf("ip 请求结果[%v]", ipInfo)

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

		s.XLog.Debugf("ip 提取后结果 [%v]-[%v]", tIP, tLoc)
	}

	return
}

func (s *NetWork) NetWorkHostName() (string, error) {
	hostInfo, err := utils.RunBash(cmds.ScriptHostName)
	if err != nil {
		s.XLog.Errorf("shell 脚本 [ScriptHostName] 执行失败 [%v]", err)
		return "", err
	}

	return hostInfo, nil
}

func (s *NetWork) NetWorkDhcp() (res []*vmodel.Dhcp, err error) {
	res = make([]*vmodel.Dhcp, 0)
	dhcpInfo, err := utils.RunBash(cmds.ScriptNetWorkDhcp)
	if err != nil {
		s.XLog.Errorf("shell 脚本 [ScriptNetWorkDhcp] 执行失败 [%v]", err)
		return nil, err
	}

	s.XLog.Debugf("返回的 dhcp 列表 [%v]", dhcpInfo)

	if !rose.StrIsEmpty(dhcpInfo) {
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
