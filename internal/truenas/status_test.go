package truenas

import "testing"

func TestGetAppStatus(t *testing.T) {
	apps := []App{
		{Name: "app1", UpgradeAvailable: false},
		{Name: "app2", UpgradeAvailable: true},
		{Name: "app3", UpgradeAvailable: false},
		{Name: "app4", UpgradeAvailable: true},
		{Name: "app5", UpgradeAvailable: true},
	}

	upToDate, toUpgrade := GetAppStatus(apps)

	if upToDate != 2 {
		t.Errorf("expected 2 apps up to date, got %d", upToDate)
	}
	if toUpgrade != 3 {
		t.Errorf("expected 3 apps to upgrade, got %d", toUpgrade)
	}
}
