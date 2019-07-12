package cuttly

import (
	"encoding/json"
)

//ForceDevicesToRightType uses a json conversion as intermediary for filling the Stats.Devices
// struct with map[string]interface{} values
func ForceDevicesToRightType(dev interface{}) (Devices, error) {
	temp, err := json.Marshal(dev)
	if err != nil {
		return Devices{}, err
	}
	// Use a temporary variable of the right type
	var devices Devices
	err = json.Unmarshal(temp, &devices)
	if err != nil {
		return Devices{}, err
	}

	return devices, nil
}

// ForceRefsToRightType uses a json conversion as intermediary for filling the Stats.Refs
// struct with map[string]interface{} values
func ForceRefsToRightType(refs interface{}) (Refs, error) {
	temp, err := json.Marshal(refs)
	if err != nil {
		return Refs{}, err
	}
	// Use a temporary variable of the right type
	var references Refs
	err = json.Unmarshal(temp, &references)
	if err != nil {
		return Refs{}, err
	}

	return references, nil
}
