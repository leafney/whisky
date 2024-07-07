/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2024-07-06 00:34
 * @Description:
 */

package cmds

import _ "embed"

//go:embed running_time.sh
var ScriptRunningTime string

//go:embed mem_usage.sh
var ScriptMemUsage string

//go:embed boot_time.sh
var ScriptBootTime string

//go:embed disk_usage.sh
var ScriptDiskUsage string

//go:embed temp_cpu.sh
var ScriptTempCpu string

//go:embed time_now.sh
var ScriptTimeNow string

//go:embed crash_start.sh
var ScriptCrashStart string

//go:embed crash_stop.sh
var ScriptCrashStop string

//go:embed crash_exist.sh
var ScriptCrashExist string

//go:embed network_device.sh
var ScriptNetworkDevice string

//go:embed network_dhcp.sh
var ScriptNetWorkDhcp string

//go:embed hostname.sh
var ScriptHostName string

//go:embed reboot.sh
var ScriptReboot string

//go:embed test.sh
var ScriptTest string

// -----------------------------

//go:embed yacd_mode.sh.gz
var ScriptYacdMode []byte

//go:embed yacd_allow_lan.sh.gz
var ScriptYacdAllowLan []byte

//go:embed yacd_stats.sh.gz
var ScriptYacdStats []byte
