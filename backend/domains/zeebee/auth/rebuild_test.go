package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
)

// Test_newAuthWithEnforcer_seam verifies the test seam compiles + a fakeEnforcer
// satisfies the policyEnforcer interface and behaves correctly.
func Test_newAuthWithEnforcer_seam(t *testing.T) {
	fake := newFakeEnforcer()
	a := newAuthWithEnforcer(nil, nil, dto.JwtConfig{}, fake, "")
	require.NotNil(t, a)

	// seed a policy
	fake.seed([][]string{
		{"Role:1", "Lender", "Resource:5", "allow"},
	})

	ok, err := a.enforcer.Enforce("Role:1", "Lender", "Resource:5", "allow")
	require.NoError(t, err)
	assert.True(t, ok)

	ok, err = a.enforcer.Enforce("Role:2", "Lender", "Resource:5", "allow")
	require.NoError(t, err)
	assert.False(t, ok)
}

// Test_fakeEnforcer_addRemove verifies idempotency + AllPolicies correctness.
func Test_fakeEnforcer_addRemove(t *testing.T) {
	f := newFakeEnforcer()

	p1 := BuildPolicy(1, 10, constants.UserTypeLender)
	p2 := BuildPolicy(1, 20, constants.UserTypeLender)

	// add
	require.NoError(t, f.AddPolicies([][]string{p1, p2}))
	all, err := f.AllPolicies()
	require.NoError(t, err)
	assert.Len(t, all, 2)

	// idempotent add
	require.NoError(t, f.AddPolicies([][]string{p1}))
	all2, _ := f.AllPolicies()
	assert.Len(t, all2, 2)

	// remove
	require.NoError(t, f.RemovePolicies([][]string{p1}))
	all3, _ := f.AllPolicies()
	assert.Len(t, all3, 1)

	// idempotent remove
	require.NoError(t, f.RemovePolicies([][]string{p1}))
	all4, _ := f.AllPolicies()
	assert.Len(t, all4, 1)
}

// Test_fakeEnforcer_subjectIsolation verifies the subject-scoped filtering logic
// that rebuild_test.go's diff tests will exercise: policies whose subject is
// outside the scope's role-id set are left untouched.
func Test_fakeEnforcer_subjectIsolation(t *testing.T) {
	f := newFakeEnforcer()

	org10Role1 := BuildPolicy(100, 1, constants.UserTypeLender)  // Role:100
	org20Role2 := BuildPolicy(200, 2, constants.UserTypeLender)  // Role:200

	f.seed([][]string{org10Role1, org20Role2})

	all, err := f.AllPolicies()
	require.NoError(t, err)
	assert.Len(t, all, 2)

	// verify org-scoped filter: policies for Role:100 only
	org10Set := map[string]int64{"Role:100": 100}
	count := 0
	for _, p := range all {
		if _, ok := org10Set[p[0]]; ok {
			count++
		}
	}
	assert.Equal(t, 1, count, "only one policy should belong to Role:100")
}
