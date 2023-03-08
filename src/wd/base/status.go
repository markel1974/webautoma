package base

type Status struct {
	Java struct {
		Version string
	}
	Build struct {
		Version  string
		Revision string
		Time     string
	}
	OS struct {
		Arch    string
		Name    string
		Version string
	}
	Ready   bool
	Message string
}
