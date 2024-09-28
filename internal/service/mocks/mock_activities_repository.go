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
	r := &mockActivitiesRepository{
		m: make(map[string]sm.Activity),
	}

	desc1 := "Студсовет СМ – это не просто факультетская организация, состоящая из четырех отделов. Это объединение людей, их идей, и, что самое главное, стремлений идти в ногу со временем. Вместе с ребятами ты сможешь учиться и тусоваться, работать и путешествовать, творить и развиваться! В такой компании даже учёба идёт легче: рядом всегда есть товарищи, которые вдохновляют идти только вперёд! Тебе помогут начать свой путь в общественную деятельность и обучим необходимым навыкам! К активистам Студсовета СМ ты всегда можешь обратиться за помощью и поддержкой! Ребята открыты для всех студентов, стараемся поддерживать идеи и инициативы каждого. На протяжении всего времени существования организация уверенно развивается и объединяет активистов разных курсов и кафедр. За 5 лет студенты провели большое количество уникальных мероприятий, и их число продолжает расти. А точкой зарождения новых идей неизменно становится аудитория 509м."
	loc1 := "509м"
	err := r.Save(nil, sm.MustNewActivity(
		"Студенческий совет СМ",
		&desc1,
		&loc1,
		[]sm.User{
			sm.MustNewUser("zhikhkirill", sm.Administrator),
		},
		[]sm.SkillType{sm.Social},
		8,
		[]*sm.Slot{
			sm.MustNewSlot(timeWithMinutes(0), timeWithMinutes(15)),
			sm.MustNewSlot(timeWithMinutes(15), timeWithMinutes(30)),
			sm.MustNewSlot(timeWithMinutes(30), timeWithMinutes(45)),
			sm.MustNewSlot(timeWithMinutes(45), timeWithMinutes(60)),
			sm.MustNewSlot(timeWithMinutes(60), timeWithMinutes(75)),
			sm.MustNewSlot(timeWithMinutes(75), timeWithMinutes(90)),
			sm.MustNewSlot(timeWithMinutes(90), timeWithMinutes(105)),
			sm.MustNewSlot(timeWithMinutes(105), timeWithMinutes(120)),
		},
	))
	if err != nil {
		panic(err)
	}

	desc2 := "Bauman Racing Team основана в 2012 году и за свою историю постоянно меняющийся коллектив студентов смог успешно реализовать 8 проектов гоночных болидов. В том числе, первый в России беспилотный гоночный электроболид. В данный момент команда занимается разработкой второго беспилотного гоночного болида в своей истории. Организация ставит перед собой масштабную цель: оказаться в числе первых в мире студенческих команд и проектов, воспитать новое поколение инженеров и оставить свой след в истории. Гоночную технику собирают студенты, охватывая все стадии создания гоночного болида, организовав производство как настоящий бизнес-проект. Перед командой стоит задача не просто спроектировать машину, но и успешно выступить в гоночных соревнованиях, а также продумать бизнес-проект команды до мельчайших деталей и \"продать\" его жюри."
	loc2 := "УЛК"
	err = r.Save(nil, sm.MustNewActivity(
		"Bauman Racing Team",
		&desc2,
		&loc2,
		[]sm.User{
			sm.MustNewUser("zhikhkirill", sm.Administrator),
		},
		[]sm.SkillType{sm.Engineering, sm.Social},
		4,
		[]*sm.Slot{
			sm.MustNewSlot(timeWithMinutes(0), timeWithMinutes(15)),
			sm.MustNewSlot(timeWithMinutes(15), timeWithMinutes(30)),
			sm.MustNewSlot(timeWithMinutes(30), timeWithMinutes(45)),
			sm.MustNewSlot(timeWithMinutes(45), timeWithMinutes(60)),
		},
	))
	if err != nil {
		panic(err)
	}

	return r
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
