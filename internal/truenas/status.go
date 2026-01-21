package truenas

// GetAppStatus counts the number of apps that are up-to-date vs. need an upgrade.
func GetAppStatus(apps []App) (upToDate int, toUpgrade int) {
	for _, app := range apps {
		if app.UpgradeAvailable {
			toUpgrade++
		} else {
			upToDate++
		}
	}
	return upToDate, toUpgrade
}
