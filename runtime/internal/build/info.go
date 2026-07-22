package build

var (
	Version = "dev"
)

type Info struct {
	Version string
}

func Current() Info {
	return Info{
		Version: Version,
	}
}
