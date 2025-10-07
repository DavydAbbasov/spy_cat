package service

import (
	"context"
	"errors"
	"testing"

	"github.com/DavydAbbasov/spy-cat/internal/domain"
	serviceserrors "github.com/DavydAbbasov/spy-cat/internal/servies_errors"
)

type mockBreedValidator struct {
	ok    bool
	err   error
	calls int
}

func (m *mockBreedValidator) IsValid(ctx context.Context, breed string) (bool, error) {
	m.calls++
	return m.ok, m.err
}

type mockRepo struct {
	createCalled bool
	retID        int64
	retErr       error
}

func (r *mockRepo) CreateCat(ctx context.Context, cat *domain.Cat) (int64, error) {
	r.createCalled = true
	return r.retID, r.retErr
}
func (r *mockRepo) ListCats(ctx context.Context, p domain.ListCatsParams) ([]domain.Cat, error) {
	return nil, nil
}
func (r *mockRepo) GetCat(ctx context.Context, id int64) (domain.Cat, error) {
	return domain.Cat{}, nil
}
func (r *mockRepo) DeleteCat(ctx context.Context, id int64) (int64, error) {
	return 0, nil
}
func (r *mockRepo) UpdateSalary(ctx context.Context, id int64, salary float64) (domain.Cat, error) {
	return domain.Cat{}, nil
}

func TestCreateCat_BreedValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		validatorOK      bool
		validatorErr     error
		repoID           int64
		wantErr          error
		wantCreateCalled bool
	}{
		{
			name:             "valid breed -> repo.CreateCat called",
			validatorOK:      true,
			repoID:           101,
			wantErr:          nil,
			wantCreateCalled: true,
		},
		{
			name:             "invalid breed -> ErrBreedInvalid",
			validatorOK:      false,
			wantErr:          serviceserrors.ErrBreedInvalid,
			wantCreateCalled: false,
		},
		{
			name:             "external error -> ErrExternalService",
			validatorErr:     errors.New("timeout"),
			wantErr:          serviceserrors.ErrExternalService,
			wantCreateCalled: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			val := &mockBreedValidator{ok: tc.validatorOK, err: tc.validatorErr}
			repo := &mockRepo{retID: tc.repoID}

			svc := NewCatService(repo, val)

			cat := &domain.Cat{
				Name:            "bro this is rapchik",
				YearsExperience: 33,
				Breed:           "siamese",
				Salary:          11111,
			}

			id, err := svc.CreateCat(context.Background(), cat)

			if tc.wantErr == nil {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if id != tc.repoID {
					t.Fatalf("want id=%d, got %d", tc.repoID, id)
				}
			} else {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("want error %v, got %v", tc.wantErr, err)
				}
			}

			if repo.createCalled != tc.wantCreateCalled {
				t.Fatalf("repo.CreateCat called=%v, want %v", repo.createCalled, tc.wantCreateCalled)
			}
			if val.calls != 1 {
				t.Fatalf("breed validator calls=%d, want 1", val.calls)
			}
		})
	}
}
