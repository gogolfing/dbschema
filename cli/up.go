package cli

const UpgradeName = "upgrade"

var upgradeAliases = []string{"up"}

type upgrade struct {
	count int
}
