package gopifinder

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
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
	IsThrottled  bool      `json:"isThrottled"`  // If CPU is currently throttled
	DiskUsed     int64     `json:"freeDisk"`     // Disk Used Space in bytes
	DiskUsedPerc int       `json:"freeDiskPerc"` // Disk Used Space in percentage
	TotalMem     int64     `json:"totalMem"`     // Total Memory in bytes
	AvailMem     int64     `json:"availMem"`     // Available Memory in bytes
	Uptime       int       `json:"uptime"`       // CPU uptime in seconds
	Created      time.Time `json:"created"`      // The date and time the status was created
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
		if strings.Contains(err.Error(), "executable file not found") {
			d.OS = "WindowsNT"
		} else {
			return d, errors.New("Error getting device Operating System. " + err.Error())
		}
	} else {
		d.OS = strings.TrimSpace(string(out))
	}

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
	v, err := strconv.ParseFloat(strings.TrimSpace(txt), 64)
	if err != nil {
		return errors.New("Could not parse CPU Temperature. " + txt)
	}
	d.CPUTemp = v / 1000

	//get gpu temperature
	out, err := exec.Command("/opt/vc/bin/vcgencmd", "measure_temp").Output()
	if err != nil {
		return errors.New("Error getting GPU temperature. " + err.Error())
	}
	txt = string(out)
	txt = txt[5 : len(txt)-3]
	v, err = strconv.ParseFloat(txt, 64)
	if err != nil {
		return errors.New("Error parsing GPU temperature. " + err.Error())
	}
	d.GPUTemp = v

	//get throttled status
	out, err = exec.Command("/opt/vc/bin/vcgencmd", "get_throttled").Output()
	if err != nil {
		return errors.New("Error getting Throttled state. " + err.Error())
	}
	txt = string(out)
	txt = txt[12 : len(txt)-1]
	u, err := strconv.ParseUint(txt, 16, 64)
	if err != nil {
		return errors.New("Error parsing Throttled State. " + err.Error())
	}
	d.IsThrottled = (u&2 == 2)

	//get disk space
	re = regexp.MustCompile("/dev/root\\s*(\\d*)\\s*(\\d*)\\s*(\\d*)\\s*(\\d*)%")
	out, err = exec.Command("df").Output()
	if err != nil {
		return errors.New("Error getting Disk space. " + err.Error())
	}
	m = re.FindStringSubmatch(string(out))
	if len(m) >= 5 {
		v, err := strconv.ParseInt(m[3], 10, 64)
		if err != nil {
			return errors.New("Error parsing Disk space. " + err.Error())
		}
		d.DiskUsed = v

		n, err := strconv.Atoi(m[4])
		if err != nil {
			return errors.New("Error parsing Disk percentage. " + err.Error())
		}
		d.DiskUsedPerc = n
	}

	//get available memory
	re = regexp.MustCompile("MemAvailable:\\s*(\\d*)")
	re1 = regexp.MustCompile("MemTotal:\\s*(\\d*)")
	txt, err = ReadAllText("/proc/meminfo")
	if err != nil {
		return errors.New("Error getting available memory. " + err.Error())
	}
	m = re.FindStringSubmatch(txt)
	if len(m) >= 2 {
		v, err := strconv.ParseInt(m[1], 10, 64)
		if err != nil {
			return errors.New("Error parsing available memory. " + err.Error())
		}
		d.AvailMem = v
	}
	m = re1.FindStringSubmatch(txt)
	if len(m) >= 2 {
		v, err := strconv.ParseInt(m[1], 10, 64)
		if err != nil {
			return errors.New("Error parsing totoal memory. " + err.Error())
		}
		d.TotalMem = v
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
	out, err := exec.Command("wmic", "os", "get", "/value").Output()
	if err != nil {
		return errors.New("Error getting Operating System information. " + err.Error())
	}
	arr := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, i := range arr {
		if strings.HasPrefix(i, "Caption=") {
			d.OSName = strings.TrimSpace(i[8:])
		}
		if strings.HasPrefix(i, "FreePhysicalMemory=") {
			d.AvailMem = ConvToInt64(i[19:], 0)
		}
		if strings.HasPrefix(i, "LastBootUpTime=") {
			t := ConvToDate(i[15:])
			d.Uptime = int(time.Since(t).Seconds())
		}
		if strings.HasPrefix(i, "Version=") {
			d.OSVersion = strings.TrimSpace(i[8:])
		}
	}

	// Get the baseboard information
	out, err = exec.Command("wmic", "baseboard", "get", "/value").Output()
	if err != nil {
		return errors.New("Error getting Baseboard information. " + err.Error())
	}
	arr = strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, i := range arr {
		if strings.HasPrefix(i, "SerialNumber=") {
			d.HWSerialNo = strings.TrimSpace(i[13:])
			break
		}
		if strings.HasPrefix(i, "Manufacturer=") {
			d.HWType = strings.TrimSpace(i[13:])
		}
	}

	// Get the logical disk information
	out, err = exec.Command("wmic", "logicaldisk", "get", "/value").Output()
	if err != nil {
		return errors.New("Error getting Baseboard information. " + err.Error())
	}
	arr = strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, i := range arr {
		if strings.HasPrefix(i, "FreeSpace=") {
			d.DiskUsed = ConvToInt64(i[10:], 0)
		}
		if strings.HasPrefix(i, "Size=") {
			n := ConvToInt64(i[5:], 0)
			d.DiskUsedPerc = int(float64(d.DiskUsed) / float64(n) * float64(100))
			break
		}
	}

	return nil
}

// ReadFrom reads the string from the reader and deserializes it into the entity values
func (d *DeviceStatus) ReadFrom(r io.ReadCloser) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	if b != nil && len(b) != 0 {
		if err := json.Unmarshal(b, &d); err != nil {
			return err
		}
	}
	return nil
}

// WriteTo serializes the entity and writes it to the http response
func (d *DeviceStatus) WriteTo(w http.ResponseWriter) error {
	b, err := json.Marshal(d)
	if err != nil {
		return err
	}
	w.Header().Set("content-type", "application/json")
	w.Write(b)
	return nil
}

// Serialize serializes the entity and returns the serialized string
func (d *DeviceStatus) Serialize() (string, error) {
	b, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Deserialize deserializes the specified string into the entity values
func (d *DeviceStatus) Deserialize(v string) error {
	err := json.Unmarshal([]byte(v), &d)
	if err != nil {
		return err
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
	case "0010", "0013":
		return "Raspberry Pi B+"
	case "0011":
		return "Raspberry Pi Compute Module"
	case "0012", "0015":
		return "Raspberry Pi A+"
	case "a01041", "a21041":
		return "Raspberry Pi 2B"
	case "900092", "900093":
		return "Raspberry Pi Zero"
	case "a02082", "a22082", "a32082", "a52082", "a22083":
		return "Raspberry Pi 3B"
	case "a020d3":
		return "Raspberry Pi 3B+"
	case "a03111":
		return "Raspberry Pi 4 1Gb"
	case "b03111":
		return "Raspberry Pi 4 2Gb"
	case "c03111":
		return "Raspberry Pi 4 4Gb"
	case "9000c1":
		return "Raspberry Pi Zero W"
	default:
		return "Unknown Model"
	}
}
