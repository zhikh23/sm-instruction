package mocks

import (
	"context"
	"sync"
	"time"

	"github.com/zhikh23/sm-instruction/internal/domain/sm"
	"github.com/zhikh23/sm-instruction/pkg/funcs"
)

type mockActivitiesRepository struct {
	m map[string]sm.Activity
	sync.RWMutex
}

func NewMockActivitiesRepository() sm.ActivitiesRepository {
	return &mockActivitiesRepository{
		m: make(map[string]sm.Activity),
	}
}

func (r *mockActivitiesRepository) Save(_ context.Context, activity *sm.Activity) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.m[activity.Name]; ok {
		return sm.ErrActivityAlreadyExists
	}

	r.m[activity.Name] = *activity

	return nil
}

func (r *mockActivitiesRepository) Activity(_ context.Context, uuid string) (*sm.Activity, error) {
	r.RLock()
	defer r.RUnlock()

	act, ok := r.m[uuid]
	if !ok {
		return nil, sm.ErrActivityNotFound
	}

	return &act, nil
}

func (r *mockActivitiesRepository) ActivityByAdmin(_ context.Context, adminUsername string) (*sm.Activity, error) {
	r.RLock()
	defer r.RUnlock()

	for _, act := range r.m {
		for _, admin := range act.Admins {
			if admin.Username == adminUsername {
				return &act, nil
			}
		}
	}

	return nil, sm.ErrActivityNotFound
}

func (r *mockActivitiesRepository) Activities(_ context.Context) ([]*sm.Activity, error) {
	r.RLock()
	defer r.RUnlock()

	return r.activities(), nil
}

func (r *mockActivitiesRepository) AvailableActivities(_ context.Context) ([]*sm.Activity, error) {
	r.RLock()
	defer r.RUnlock()

	return funcs.Filter(
		r.activities(),
		func(a *sm.Activity) bool {
			return a.Location != nil
		},
	), nil
}

func (r *mockActivitiesRepository) AdditionalActivities(_ context.Context) ([]*sm.Activity, error) {
	r.RLock()
	defer r.RUnlock()

	return funcs.Filter(
		r.activities(),
		func(a *sm.Activity) bool {
			return a.Location == nil
		},
	), nil
}

func (r *mockActivitiesRepository) UpdateSlots(
	ctx context.Context,
	activityUUID string,
	updateFn func(context.Context, *sm.Activity) error,
) error {
	r.Lock()
	defer r.Unlock()

	act, ok := r.m[activityUUID]
	if !ok {
		return sm.ErrActivityNotFound
	}

	err := updateFn(ctx, &act)
	if err != nil {
		return err
	}

	r.m[activityUUID] = act

	return nil
}

func (r *mockActivitiesRepository) activities() []*sm.Activity {
	acts := make([]*sm.Activity, 0, len(r.m))
	for _, act := range r.m {
		acts = append(acts, &act)
	}
	return acts
}

func timeWithMinutes(minutes int) time.Time {
	return time.Now().Round(time.Minute).Add(time.Duration(minutes) * time.Minute)
}
