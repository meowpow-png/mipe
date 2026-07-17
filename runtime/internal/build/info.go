package build

var (
	Version = "dev"
	Commit  = ""
	Date    = ""
)

type Info struct {
	Version string
	Commit  string
	Date    string
}

func Current() Info {
	return Info{
		Version: Version,
		Commit:  Commit,
		Date:    Date,
	}
}
