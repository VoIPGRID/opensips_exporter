package opensips

// Statistic holds the module, name and value of a statistic
// as returned by OpenSIPS.
type Statistic struct {
	Module, Name string
	Value        float64
}
