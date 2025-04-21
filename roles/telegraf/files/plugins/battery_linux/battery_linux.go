//	code from: https://dev.sigpipe.me/dashie/telegraf-plugins

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func batteryThings(battery string) (map[string]string, bool) {
	things := make(map[string]string)

	filename := fmt.Sprintf("/sys/class/power_supply/%s/uevent", battery)
	//fmt.Println(os.Stderr, filename)

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s error trying to open the file: %s", filename, err)
		return nil, true
	}

	things["BAT_NAME"] = battery // Add battery name "BATx" to the things

	for _, line := range strings.Split(string(b), "\n") {
		if strings.TrimSpace(line) == "" {
			continue // if empty
		}
		l := strings.Split(line, "=")
		things[l[0]] = l[1]
	}

	return things, false
}

func printBattery(things map[string]string) {
	idx := strings.Replace(things["BAT_NAME"], "BAT", "", 1)

	// tags
	fmt.Printf("battery_linux,battery=%s,idx=%s ", things["BAT_NAME"], idx)

	// values
	fmt.Printf("present=%s,", things["POWER_SUPPLY_PRESENT"])
	fmt.Printf("cycle_count=%s,", things["POWER_SUPPLY_CYCLE_COUNT"])
	fmt.Printf("voltage_min_design=%s,", things["POWER_SUPPLY_VOLTAGE_MIN_DESIGN"])
	fmt.Printf("voltage_now=%s,", things["POWER_SUPPLY_VOLTAGE_NOW"])
	fmt.Printf("power_now=%s,", things["POWER_SUPPLY_POWER_NOW"])
	fmt.Printf("energy_full_design=%s,", things["POWER_SUPPLY_ENERGY_FULL_DESIGN"])
	fmt.Printf("energy_full=%s,", things["POWER_SUPPLY_ENERGY_FULL"])
	fmt.Printf("energy_now=%s,", things["POWER_SUPPLY_ENERGY_NOW"])
	fmt.Printf("capacity=%s,", things["POWER_SUPPLY_CAPACITY"])
	fmt.Printf("battery_status=\"%s\",", things["POWER_SUPPLY_STATUS"])
	fmt.Printf("capacity_level=\"%s\"\n", things["POWER_SUPPLY_CAPACITY_LEVEL"])
	return
}

func main() {
	files, _ := ioutil.ReadDir("/sys/class/power_supply/")
	for _, f := range files {
		if strings.Contains(f.Name(), "BAT") {
			stats, err := batteryThings(f.Name())
			if err {
				fmt.Fprintln(os.Stderr, "BatteryThings error")
				return
			}
			printBattery(stats)
		}
	}
}
