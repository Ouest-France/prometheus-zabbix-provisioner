package main

import (
	"fmt"
	"os"
)

func main() {

	prometheusAddr := os.Getenv("PZP_PROMETHEUS_ADDR")
	prometheusUser := os.Getenv("PZP_PROMETHEUS_USER")
	prometheusPassword := os.Getenv("PZP_PROMETHEUS_PASSWORD")
	zabbixAddr := os.Getenv("PZP_ZABBIX_ADDR")
	zabbixUser := os.Getenv("PZP_ZABBIX_USER")
	zabbixPassword := os.Getenv("PZP_ZABBIX_PASSWORD")
	zabbixHost := os.Getenv("PZP_ZABBIX_HOST")

	// Create prometheus connection
	pc := PrometheusClient{
		Server:   prometheusAddr,
		User:     prometheusUser,
		Password: prometheusPassword,
	}

	// Create zabbix connection
	zc, err := NewZabbixClient(zabbixAddr, zabbixUser, zabbixPassword)
	if err != nil {
		fmt.Println(err)
		return
	}
	zc.Host = zabbixHost

	// Get prometheus alerts
	prometheusAlerts, err := pc.GetAlerts()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Get zabbix alerts
	zabbixAlerts, err := zc.GetAlerts()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Calculate diff between prometheus and zabbix alerts
	add, remove := DiffAlerts(prometheusAlerts, zabbixAlerts)

	fmt.Println("add", len(add), "update", 0, "remove", len(remove))

	// Remove Zabbix alerts no longer present in Prometheus
	for _, alert := range remove {
		fmt.Println("Remove", alert.Name, alert.Severity)

		err = zc.RemoveAlert(alert.Name)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// Add missing alerts in Zabbix
	for _, alert := range add {
		fmt.Println("Add", alert.Name, alert.Severity)

		err = zc.AddAlert(alert.Name, alert.Description, alert.Severity)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
