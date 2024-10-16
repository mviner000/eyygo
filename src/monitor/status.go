package monitor

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

var StartTime = time.Now()

type SystemMetrics struct {
	CPUUsage     float64
	RAMUsage     float64
	StorageUsage float64
}

func getSystemMetrics() (SystemMetrics, error) {
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return SystemMetrics{}, err
	}

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return SystemMetrics{}, err
	}

	diskStat, err := disk.Usage("/")
	if err != nil {
		return SystemMetrics{}, err
	}

	return SystemMetrics{
		CPUUsage:     cpuPercent[0],
		RAMUsage:     vmStat.UsedPercent,
		StorageUsage: diskStat.UsedPercent,
	}, nil
}

func HandleStatus(c *fiber.Ctx) error {
	metrics, err := getSystemMetrics()
	if err != nil {
		return c.Status(500).SendString("Error fetching system metrics")
	}

	status := struct {
		Time             string
		ConnectedClients int
		ActiveRooms      int
		Uptime           string
		RoomInfo         []map[string]interface{}
		CPUUsage         float64
		RAMUsage         float64
		StorageUsage     float64
	}{
		Time:         time.Now().Format(time.RFC1123),
		Uptime:       time.Since(StartTime).Round(time.Second).String(),
		CPUUsage:     metrics.CPUUsage,
		RAMUsage:     metrics.RAMUsage,
		StorageUsage: metrics.StorageUsage,
	}

	switch c.Path() {
	case "/status/server-info":
		return c.Render("eyygo/monitor/views/status_partial", status)
	case "/status/cpu":
		return c.SendString(fmt.Sprintf("%.2f", metrics.CPUUsage))
	case "/status/ram":
		return c.SendString(fmt.Sprintf("%.2f", metrics.RAMUsage))
	case "/status/storage":
		return c.SendString(fmt.Sprintf("%.2f", metrics.StorageUsage))
	case "/status/old":
		return c.Render("eyygo/monitor/views/status_partial", status)
	case "/status/room-console":
		return c.Render("eyygo/monitor/views/room_console", status)
	default:
		return c.Render("eyygo/monitor/views/status", status)
	}
}
