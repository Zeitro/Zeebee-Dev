package auth

// fakeEnforcer is an in-memory policyEnforcer for rebuild diff tests.
// It stores policies in a Go map keyed by comma-joined policy tuples,
// mirroring the set semantics of casbin (duplicate adds are idempotent).
type fakeEnforcer struct {
	policies map[string][]string // keyed by BuildPolicy key
}

func newFakeEnforcer() *fakeEnforcer {
	return &fakeEnforcer{policies: map[string][]string{}}
}

func keyOf(p []string) string {
	out := ""
	for i, v := range p {
		if i > 0 {
			out += ","
		}
		out += v
	}
	return out
}

func (f *fakeEnforcer) seed(p [][]string) {
	for _, pp := range p {
		f.policies[keyOf(pp)] = pp
	}
}

func (f *fakeEnforcer) AllPolicies() ([][]string, error) {
	out := make([][]string, 0, len(f.policies))
	for _, p := range f.policies {
		out = append(out, p)
	}
	return out, nil
}

func (f *fakeEnforcer) AddPolicies(ps [][]string) error {
	for _, p := range ps {
		f.policies[keyOf(p)] = p
	}
	return nil
}

func (f *fakeEnforcer) RemovePolicies(ps [][]string) error {
	for _, p := range ps {
		delete(f.policies, keyOf(p))
	}
	return nil
}

func (f *fakeEnforcer) Enforce(rvals ...interface{}) (bool, error) {
	parts := make([]string, 0, len(rvals))
	for _, r := range rvals {
		s, ok := r.(string)
		if !ok {
			return false, nil
		}
		parts = append(parts, s)
	}
	_, ok := f.policies[keyOf(parts)]
	return ok, nil
}

func (f *fakeEnforcer) SavePolicy() error { return nil }

func (f *fakeEnforcer) LoadDefaultFilteredPolicy() error { return nil }
