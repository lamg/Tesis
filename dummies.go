package tesis

type DummySync struct {
	synced bool
}

func NewDummySync() (s *DummySync) {
	s = &DummySync{synced: false}
	return
}

func (s *DummySync) Candidates() (a []AccMatch, e error) {
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

func (s *DummySync) Synchronize(a []AccMatch) (e error) {
	s.synced = true
	return
}

type DummyAuth struct {
}

func (d *DummyAuth) Authenticate(u, p string) (b bool) {
	b = u == p
	return
}

type DummyManager struct {
	syncer Synchronizer
}

func NewDummyManager() (m *DummyManager) {
	m = &DummyManager{syncer: NewDummySync()}
	return
}

func (m *DummyManager) UserInfo(u string) (inf *Info, e error) {
	var cs []AccMatch
	var re []Change
	cs, e = m.syncer.Candidates()
	if e == nil {
		re = make([]Change, 0)
		inf = &Info{Name: u, Matches: cs, Record: re}
	}
	//TODO populate inf with more meaningful information
	return
}
