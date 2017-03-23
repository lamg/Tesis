package tesis

func (s *DummyManager) Candidates() (a []AccMatch, e error) {
	if !s.synced {
		a = []AccMatch{
			AccMatch{DBId: "0", SrcID: "8901191122"},
			AccMatch{DBId: "1", ADId: "3", ADName: "LUIS", SrcName: "Luis"},
		}
	} else {
		a = make([]AccMatch, 0)
	}
	//iterate DB and filter comparing with AD
	return
}

func (s *DummyManager) Synchronize(user string, a []AccMatch) (e error) {
	//make this according the specification
	var cs []AccMatch
	cs, _ = s.Candidates()
	if len(a) == len(cs) {
		s.synced = true
	}
	return
}

type DummyAuth struct {
}

func (d *DummyAuth) Authenticate(u, p string) (b bool) {
	b = u == p
	return
}

type DummyManager struct {
	synced bool
}

func NewDummyManager() (m *DummyManager) {
	m = &DummyManager{synced: false}
	return
}

func (m *DummyManager) UserInfo(u string) (inf *Info, e error) {
	var cs []AccMatch
	var re []Change
	cs, e = m.Candidates()
	if e == nil {
		re = make([]Change, 0)
		inf = &Info{Name: u, Matches: cs, Record: re}
	}
	//TODO populate inf with more meaningful information
	return
}
