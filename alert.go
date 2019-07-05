package main

// Alert represents a generic alerting rule
type Alert struct {
	Name        string
	Description string
	Severity    string
}

// DiffAlerts calculates diff between two Alert slices
func DiffAlerts(a, b []Alert) (add, remove []Alert) {

	for _, aElem := range a {
		present := AlertInAlerts(aElem, b)

		if !present {
			add = append(add, aElem)
			continue
		}
	}

	for _, bElem := range b {
		present := AlertInAlerts(bElem, a)

		if !present {
			remove = append(remove, bElem)
		}
	}

	return add, remove
}

// AlertInAlerts test if an Alert exists in a slice of Alerts
func AlertInAlerts(search Alert, alerts []Alert) bool {

	for _, alert := range alerts {
		if search == alert {
			return true
		}
	}

	return false
}
