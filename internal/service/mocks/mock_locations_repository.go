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
	l1 := *sm.MustNewLocation(
		"1",
		"Студенческий совет факультета СМ",
		"Студсовет СМ – это не просто факультетская организация, состоящая из четырех отделов. Это объединение людей, их идей, и, что самое главное, стремлений идти в ногу со временем. Вместе с ребятами ты сможешь учиться и тусоваться, работать и путешествовать, творить и развиваться! В такой компании даже учёба идёт легче: рядом всегда есть товарищи, которые вдохновляют идти только вперёд! Тебе помогут начать свой путь в общественную деятельность и обучим необходимым навыкам! К активистам Студсовета СМ ты всегда можешь обратиться за помощью и поддержкой! Ребята открыты для всех студентов, стараемся поддерживать идеи и инициативы каждого. На протяжении всего времени существования организация уверенно развивается и объединяет активистов разных курсов и кафедр. За 5 лет студенты провели большое количество уникальных мероприятий, и их число продолжает расти. А точкой зарождения новых идей неизменно становится аудитория 509м.\n",
		"509м",
		[]sm.SkillType{sm.Social, sm.Creative})
	l2 := *sm.MustNewLocation(
		"2",
		"BRT",
		"Bauman Racing Team основана в 2012 году. В данный момент команда занимается тестированием первого беспилотного гоночного болида в России. Организация ставит перед собой масштабную цель: оказаться в числе первых в мире студенческих команд, создавших беспилотный гоночный болид с электрической силовой установкой. Автомобили собирают студенты, организовав производство как настоящий бизнес-проект. Перед командой стоит задача не просто спроектировать машину, но и сделать презентацию мелкосерийного производства и “продать” ее жюри.\n",
		"ИЦАР",
		[]sm.SkillType{sm.Engineering, sm.Sportive})
	l3 := *sm.MustNewLocation(
		"3",
		"Центр молодёжной робототехники",
		"Центр Молодежной Робототехники - это инновационное пространство, предназначенное для обучения и развития молодых талантов в области робототехники, искусственного интеллекта и программирования. Их миссия - предоставить молодежи возможность исследовать и создавать будущее с помощью передовых технологий и творчества. Здесь вы можете: обрести команду единомышленников для дальнейшей реализации проектов, востребованных на рынке; создавать свои проекты в области робототехники под руководством опытных специалистов; работать на современном оборудовании: 3D принтеры, лазерные станки, станки с программным управлением и многое другое; улучшить навыки в организации, управлении проектом: начиная с задумки и заканчивая созданием рабочего прототипа; представлять Университет на внешних площадках и инженерных соревнованиях.\n",
		"ИЦАР",
		[]sm.SkillType{sm.Engineering, sm.Creative})
	_ = l1.AddAdministrator(sm.MustNewUser("zhikhkirill", sm.Administrator))

	r := &mockLocationsRepository{
		m: map[string]sm.Location{l1.UUID: l1, l2.UUID: l2, l3.UUID: l3},
	}
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

func (r *mockLocationsRepository) LocationByAdmin(ctx context.Context, username string) (*sm.Location, error) {
	r.RLock()
	defer r.RUnlock()

	for _, loc := range r.m {
		if loc.HasAdministrator(username) {
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
