package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type PrometheusRulesResp struct {
	Data struct {
		Groups []struct {
			Rules []PrometheusRule `json:"rules"`
		} `json:"groups"`
	} `json:"data"`
}

type PrometheusRule struct {
	Name   string `json:"name"`
	Labels struct {
		Severity string `json:"severity"`
	} `json:"labels"`
	Annotations struct {
		Summary string `json:"summary"`
	} `json:"annotations"`
	Type string `json:"type"`
}

// PrometheusClient represents a PrometheusClient client instance
type PrometheusClient struct {
	Server   string
	User     string
	Password string
}

// GetAlerts returns list of Prometheus alerts
func (pc PrometheusClient) GetAlerts() ([]Alert, error) {
	client := http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/rules", pc.Server), nil)
	if err != nil {
		fmt.Println(err)
		return []Alert{}, err
	}
	req.SetBasicAuth(pc.User, pc.Password)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return []Alert{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var rulesResp PrometheusRulesResp
	err = json.Unmarshal(body, &rulesResp)
	if err != nil {
		return []Alert{}, err
	}

	var alerts []Alert
	for _, group := range rulesResp.Data.Groups {

		for _, rule := range group.Rules {
			if rule.Type != "alerting" {
				continue
			}

			alert := Alert{
				Name:        rule.Name,
				Description: rule.Annotations.Summary,
				Severity:    rule.Labels.Severity,
			}

			alerts = append(alerts, alert)
		}
	}

	return alerts, nil
}
