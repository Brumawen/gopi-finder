package gopifinder

import (
	"errors"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// DeviceStatus holds current status information about the Device
type DeviceStatus struct {
	HostName     string    `json:"hostName"`     // Current Host Name
	OS           string    `json:"os"`           // OS Type
	OSName       string    `json:"osName"`       // Operating System Name
	OSVersion    string    `json:"osVersion"`    // Operating System version
	HWType       string    `json:"hwType"`       // Hardware type
	HWSerialNo   string    `json:"hwSerialNo"`   // Hardware SerialNo
	CPUTemp      float64   `json:"cpuTemp"`      // CPU temperature in Celcius
	GPUTemp      float64   `json:"gpuTemp"`      // GPU temperature in Celcius
	FreeDisk     int       `json:"freeDisk"`     // Free Disk Space in bytes
	FreeDiskPerc int       `json:"freeDiskPerc"` // Free Disk Space in percentage
	AvailMem     int       `json:"availMem"`     // Available Memory in bytes
	Uptime       int       `json:"uptime"`       // CPU uptime in seconds
	Created      time.Time `json:"created"`
}

// NewDeviceStatus creates a new DeviceStatus struct and populates it with the values
// for the current device
func NewDeviceStatus() (DeviceStatus, error) {
	d := DeviceStatus{Created: time.Now()}

	//get hostname
	out, err := exec.Command("hostname").Output()
	if err != nil {
		return d, errors.New("Error getting HostName. " + err.Error())
	}
	d.HostName = strings.TrimSpace(string(out))

	// Get the operating system
	out, err = exec.Command("uname").Output()
	if err != nil {
		return d, errors.New("Error getting device Operating System. " + err.Error())
	}
	d.OS = strings.TrimSpace(string(out))

	if d.OS == "Linux" {
		err = d.loadValuesForLinux()
	} else if d.OS == "WindowsNT" {
		err = d.loadValuesForWindows()
	}
	return d, err
}

func (d *DeviceStatus) loadValuesForLinux() error {
	//get operating system info
	re := regexp.MustCompile("PRETTY_NAME=\"([\\w\\d\\s/()]+)\"")
	txt, err := ReadAllText("/etc/os-release")
	if err != nil {
		return errors.New("Error getting OS Information. " + err.Error())
	}
	m := re.FindStringSubmatch(txt)
	if len(m) >= 2 {
		d.OSName = m[1]
	}

	txt, err = ReadAllText("/etc/debian_version")
	if err != nil {
		return errors.New("Error getting OS Version. " + err.Error())
	}
	d.OSVersion = strings.TrimSpace(txt)

	//get hardware type
	re = regexp.MustCompile("Revision\\s:\\s([a-e\\d]+)")
	re1 := regexp.MustCompile("Serial\\s\\s:\\s([a-f\\d]*)")
	txt, err = ReadAllText("/proc/cpuinfo")
	if err != nil {
		return errors.New("Error getting Hardware type. " + err.Error())
	}
	m = re.FindStringSubmatch(txt)
	if len(m) >= 2 {
		d.HWType = getHardwareType(m[1])
	}
	m = re1.FindStringSubmatch(txt)
	if len(m) >= 2 {
		d.HWSerialNo = m[1]
	}

	//get cpu temperature
	txt, err = ReadAllText("/sys/class/thermal/thermal_zone0/temp")
	if err != nil {
		return errors.New("Error getting CPU temperature. " + err.Error())
	}
	if v, err := strconv.ParseFloat(strings.TrimSpace(txt), 64); err != nil {
		log.Println("Could not parse CPU Temperature.", txt)
	} else {
		d.CPUTemp = v / 1000
	}

	//get gpu temperature
	out, err := exec.Command("/opt/vc/bin/vcgencmd", "measure_temp").Output()
	if err != nil {
		return errors.New("Error getting GPU temperature. " + err.Error())
	}
	txt = string(out)
	txt = txt[5 : len(txt)-3]
	v, err := strconv.ParseFloat(txt, 64)
	if err != nil {
		return errors.New("Error parsing GPU temperature. " + err.Error())
	}
	d.GPUTemp = v

	//get disk space
	re = regexp.MustCompile("/dev/root\\s*(\\d*)\\s*(\\d*)\\s*(\\d*)\\s*(\\d*)%")
	out, err = exec.Command("df").Output()
	if err != nil {
		return errors.New("Error getting Disk space. " + err.Error())
	}
	m = re.FindStringSubmatch(string(out))
	if len(m) >= 5 {
		v, err := strconv.Atoi(m[3])
		if err != nil {
			return errors.New("Error parsing Disk space. " + err.Error())
		}
		d.FreeDisk = v

		v, err = strconv.Atoi(m[4])
		if err != nil {
			return errors.New("Error parsing Disk percentage. " + err.Error())
		}
		d.FreeDiskPerc = v
	}

	//get available memory
	re = regexp.MustCompile("MemAvailable:\\s*(\\d*)")
	txt, err = ReadAllText("/proc/meminfo")
	if err != nil {
		return errors.New("Error getting available memory. " + err.Error())
	}
	m = re.FindStringSubmatch(txt)
	if len(m) >= 2 {
		v, err := strconv.Atoi(m[1])
		if err != nil {
			return errors.New("Error parsing available memory. " + err.Error())
		}
		d.AvailMem = v

	}

	//get uptime
	txt, err = ReadAllText("/proc/uptime")
	if err != nil {
		return errors.New("Error getting system uptime. " + err.Error())
	}
	i := strings.IndexRune(txt, '.')
	if i >= 0 {
		v, err := strconv.Atoi(txt[:i])
		if err != nil {
			return errors.New("Error parsing system uptime. " + err.Error())
		}
		d.Uptime = v
	}
	return nil
}

func (d *DeviceStatus) loadValuesForWindows() error {
	// Get the operating system information
	re := regexp.MustCompile("([\\w\\s]+)\\s\\[Version ([\\d\\.]+)\\]")
	out, err := exec.Command("ver").Output()
	if err != nil {
		return errors.New("Error getting Operating System information. " + err.Error())
	}
	txt := strings.TrimSpace(string(out))
	m := re.FindStringSubmatch(txt)
	if len(m) >= 3 {
		d.OSName = m[1]
		d.OSVersion = m[2]
	}
	return nil
}

func getHardwareType(code string) string {
	switch code {
	case "0002", "0003":
		return "Raspberry Pi B rev 1.0"
	case "0004", "0005", "0006", "000d", "000e", "000f":
		return "Raspberry Pi B rev 2.0"
	case "0007", "0008", "0009":
		return "Raspberry Pi A"
	case "0010":
		return "Raspberry Pi B+"
	case "0011":
		return "Raspberry Pi Compute Module"
	case "0012":
		return "Raspberry Pi A+"
	case "a01041", "a21041":
		return "Raspberry Pi 2B"
	case "900092", "900093":
		return "Raspberry Pi Zero"
	case "a02082", "a22082":
		return "Raspberry Pi 3B"
	case "9000c1":
		return "Raspberry Pi Zero W"
	default:
		return "Unknown Model"
	}
}
