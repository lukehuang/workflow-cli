package cmd

import (
	"fmt"
	"github.com/deis/controller-sdk-go/api"
	"github.com/deis/controller-sdk-go/appsettings"
	"strings"
)

// LabelsList list app's labels
func (d *DeisCmd) LabelsList(appID string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	appSettings, err := appsettings.List(s.Client, appID)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Printf("=== %s Label\n", appID)

	if appSettings.Label == nil || len(appSettings.Label) == 0 {
		d.Println("No labels found.")
	} else {
		d.Println(appSettings.Label)
	}

	return nil
}

// LabelsSet sets labels for app
func (d *DeisCmd) LabelsSet(appID string, labels []string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	labelsMap, err := parseLabels(labels)
	if err != nil {
		return err
	}

	d.Printf("Applying labels on %s... ", appID)

	quit := progress(d.WOut)

	_, err = appsettings.Set(s.Client, appID, api.AppSettings{Label: labelsMap})

	quit <- true
	<-quit

	if err != nil {
		return err
	}

	d.Println("done")
	return nil
}

// LabelsUnset removes labels for the app.
func (d *DeisCmd) LabelsUnset(appID string, labels []string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	labelsMap := make(map[string]interface{})

	for _, label := range labels {
		labelsMap[label] = nil
	}

	d.Printf("Removing labels on %s... ", appID)

	quit := progress(d.WOut)

	_, err = appsettings.Set(s.Client, appID, api.AppSettings{Label: labelsMap})

	quit <- true
	<-quit

	if err != nil {
		return err
	}

	d.Println("done")
	return nil
}

func parseLabels(labels []string) (map[string]interface{}, error) {
	labelsMap := make(map[string]interface{})

	for _, label := range labels {
		key, value, err := parseLabel(label)

		if err != nil {
			return nil, err
		}

		labelsMap[key] = value
	}

	return labelsMap, nil
}

func parseLabel(label string) (string, string, error) {
	parts := strings.Split(label, "=")

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf(`%s is invalid, Must be in format key=value
Examples: git_repo=https://github.com/deis/workflow team=frontend`, label)
	}

	return parts[0], parts[1], nil
}
