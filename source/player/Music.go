package player

type Music struct {
	title    string
	duration int
	id       int
}

func (music *Music) Title() string {
	return music.title
}

func (music *Music) Duration() int {
	return music.duration
}

func (music *Music) ID() int {
	return music.id
}
