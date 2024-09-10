package mocks

import (
	"context"
	"sync"

	"sm-instruction/internal/domain/sm"
)

type mockLocationsRepository struct {
	m map[string]sm.Location
	sync.RWMutex
}

func NewMockLocationsRepository() sm.LocationsRepository {
	r := &mockLocationsRepository{
		m: make(map[string]sm.Location),
	}
	r.m["1"] = *sm.MustNewLocation("1", "ССФСМ", "339м", []sm.SkillType{sm.Social, sm.Creative})
	r.m["2"] = *sm.MustNewLocation("2", "BRT", "ICAR", []sm.SkillType{sm.Engineering, sm.Sportive})
	return r
}

func (r *mockLocationsRepository) Save(_ context.Context, l *sm.Location) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.m[l.UUID]; ok {
		return sm.ErrLocationAlreadyExists
	}

	r.m[l.UUID] = *l

	return nil
}

func (r *mockLocationsRepository) Location(_ context.Context, uuid string) (*sm.Location, error) {
	r.RLock()
	defer r.RUnlock()

	loc, ok := r.m[uuid]
	if !ok {
		return nil, sm.ErrLocationNotFound
	}

	return &loc, nil
}

func (r *mockLocationsRepository) LocationByName(_ context.Context, name string) (*sm.Location, error) {
	r.RLock()
	defer r.RUnlock()

	for _, loc := range r.m {
		if loc.Name == name {
			return &loc, nil
		}
	}

	return nil, sm.ErrLocationNotFound
}

func (r *mockLocationsRepository) Locations(_ context.Context) ([]*sm.Location, error) {
	r.RLock()
	defer r.RUnlock()

	ls := make([]*sm.Location, 0, len(r.m))
	for _, loc := range r.m {
		ls = append(ls, &loc)
	}
	return ls, nil
}

func (r *mockLocationsRepository) Update(
	ctx context.Context,
	locationUUID string,
	updateFn func(innerCtx context.Context, loc *sm.Location) error,
) error {
	r.Lock()
	defer r.Unlock()

	loc, ok := r.m[locationUUID]
	if !ok {
		return sm.ErrLocationNotFound
	}

	err := updateFn(ctx, &loc)
	if err != nil {
		return err
	}

	r.m[locationUUID] = loc

	return nil
}
