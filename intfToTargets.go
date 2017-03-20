package tesis

type DBSync struct {
}

func NewDBSync() (d *DBSync, e error) {
	return
}

func (d *DBSync) Synchronize(a []AccMatch) (e error) {
	return
}

func (d *DBSync) Candidates() (a []AccMatch, e error) {
	return
}

//all Candidates not synchronized are stored for
//future synchronizations
