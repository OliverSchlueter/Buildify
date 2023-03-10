package builds

var Builds []Build

type Build struct {
	id      int
	hash    string
	message string
}

func (receiver Build) GetId() int {
	return receiver.id
}

func Create(id int, hash string, message string) Build {
	b := Build{
		id:      id,
		hash:    hash,
		message: message,
	}

	Builds = append(Builds, b)
	return b
}
