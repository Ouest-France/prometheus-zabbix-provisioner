package main

import (
	"fmt"
	"strings"

	"github.com/Ouest-France/zabbix"
)

// ZabbixClient represents a Zabbix client instance
type ZabbixClient struct {
	Server   string
	User     string
	Password string
	Host     string
	API      *zabbix.API
}

// NewZabbixClient creates a ZabbixClient with given parameters
func NewZabbixClient(server, user, password string) (ZabbixClient, error) {

	client := ZabbixClient{
		Server:   server,
		User:     user,
		Password: password,
	}

	client.API = zabbix.NewAPI(server + "/api_jsonrpc.php")
	_, err := client.API.Login(user, password)

	return client, err
}

// GetAlerts returns Zabbix items/triggers as slice of alerts
func (zc ZabbixClient) GetAlerts() ([]Alert, error) {

	itemParams := zabbix.Params{
		"host": zc.Host,
	}

	items, err := zc.API.ItemsGet(itemParams)
	if err != nil {
		return []Alert{}, err
	}

	var alerts []Alert
	for _, item := range items {

		triggerParams := zabbix.Params{
			"host":    zc.Host,
			"itemids": []string{item.ItemId},
		}

		triggers, err := zc.API.TriggersGet(triggerParams)
		if err != nil {
			return alerts, err
		}

		if len(triggers) != 1 {
			return alerts, fmt.Errorf("Found %d trigger(s) for %q, expected 1", len(triggers), item.Name)
		}

		alert := Alert{
			Name:        item.Name,
			Description: item.Description,
			Severity:    getSeverityString(triggers[0].Priority),
		}

		alerts = append(alerts, alert)
	}

	return alerts, nil
}

// AddAlert creates a new item/trigger in Zabbix
func (zc ZabbixClient) AddAlert(name, description, severity string) error {

	host, err := zc.API.HostGetByHost(zc.Host)
	if err != nil {
		return err
	}

	fmt.Printf("Add zabbix item %q\n", name)

	item := zabbix.Item{
		HostId:      host.HostId,
		Name:        name,
		Key:         fmt.Sprintf("prometheus.%s", strings.ToLower(name)),
		Type:        zabbix.ZabbixTrapper,
		ValueType:   zabbix.Unsigned,
		Description: description,
	}

	err = zc.API.ItemsCreate(zabbix.Items{item})
	if err != nil {
		return err
	}

	fmt.Printf("Add zabbix trigger %q with severity %q\n", name, severity)

	trigger := zabbix.Trigger{
		Description: name,
		Expression:  fmt.Sprintf("{%s:prometheus.%s.last(#1)}>0", zc.Host, strings.ToLower(name)),
		Priority:    getZabbixSeverity(severity),
	}

	err = zc.API.TriggersCreate(zabbix.Triggers{trigger})

	return err
}

// RemoveAlert deletes an item/trigger from Zabbix
func (zc ZabbixClient) RemoveAlert(name string) error {

	itemParams := zabbix.Params{
		"host": zc.Host,
	}

	items, err := zc.API.ItemsGet(itemParams)
	if err != nil {
		return err
	}

	for _, item := range items {
		if item.Name != name {
			continue
		}

		triggerParams := zabbix.Params{
			"host":    zc.Host,
			"itemids": []string{item.ItemId},
		}

		triggers, err := zc.API.TriggersGet(triggerParams)
		if err != nil {
			return err
		}

		if len(triggers) != 1 {
			return fmt.Errorf("Found %d trigger(s) for %q, expected 1", len(triggers), item.Name)
		}

		fmt.Printf("Remove zabbix trigger %q\n", name)

		err = zc.API.TriggersDelete(zabbix.Triggers{triggers[0]})
		if err != nil {
			return err
		}

		fmt.Printf("Remove zabbix item %q\n", name)

		err = zc.API.ItemsDelete(zabbix.Items{item})

		return err
	}

	return fmt.Errorf("Item %q not found", name)
}

// getZabbixSeverity converts severity string to zabbix.SeverityType
func getZabbixSeverity(severity string) zabbix.SeverityType {

	switch strings.ToLower(severity) {
	case "information":
		return zabbix.Information
	case "warning":
		return zabbix.Warning
	case "average":
		return zabbix.Average
	case "high":
		return zabbix.High
	case "critical":
		return zabbix.Critical
	default:
		return zabbix.NotClassified
	}
}

// getSeverityString converts zabbix.SeverityType to severity string
func getSeverityString(severity zabbix.SeverityType) string {

	switch severity {
	case zabbix.Information:
		return "information"
	case zabbix.Warning:
		return "warning"
	case zabbix.Average:
		return "average"
	case zabbix.High:
		return "high"
	case zabbix.Critical:
		return "critical"
	default:
		return "none"
	}
}
