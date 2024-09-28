package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/zhikh23/sm-instruction/internal/adapters"
	"github.com/zhikh23/sm-instruction/internal/domain/sm"
)

var Groups = []string{
	"СМ1-11", "СМ1-11Б", "СМ1-12", "СМ10-11", "СМ10-11Б", "СМ10-12", "СМ11-11Б", "СМ12-11", "СМ13-11Б", "СМ13-12Б",
	"СМ13-13", "СМ2-11", "СМ3-11", "СМ3-12", "СМ4-11", "СМ4-12", "СМ4-13", "СМ5-11", "СМ5-11Б", "СМ5-12Б", "СМ6-11",
	"СМ6-12", "СМ6-19", "СМ7-11Б", "СМ7-12Б", "СМ7-13Б", "СМ7-14Б", "СМ8-11", "СМ9-11", "СМ9-12", "СМ9-13",
}

var Admins = []string{
	"cmr_admin", "brt_admin", "mkc_admin", "hyd_admin", "nti_admin", "sso_admin", "ssf_admin", "prf_admin", "trj_admin",
}

const slotDuration = 20 * time.Minute

func main() {
	importUsers()
	importCharacters()
	importActivities()
}

func importUsers() {
	repos, closeFn := adapters.NewPGUsersRepository()
	defer func() {
		err := closeFn()
		if err != nil {
			log.Fatal(err)
		}
	}()

	users := make([]sm.User, len(Groups)+len(Admins))
	for i, group := range Groups {
		users[i] = sm.MustNewUser(strings.ToLower(group), sm.Participant)
	}
	users[0].Username = "zhikhkirill"
	for i, admin := range Admins {
		users[len(Groups)+i] = sm.MustNewUser(strings.ToLower(admin), sm.Participant)
	}

	ctx := context.Background()
	for _, user := range users {
		if err := repos.Save(ctx, user); err != nil {
			log.Println(err.Error())
		}
	}
}

func importCharacters() {
	repos, closeFn := adapters.NewPGCharactersRepository()
	defer func() {
		err := closeFn()
		if err != nil {
			log.Fatal(err)
		}
	}()

	table := createTable()

	charSlots := func(groupName string) []*sm.Slot {
		slots := make([]*sm.Slot, 0)
		for _, row := range table.Rows {
			for x, group := range row.Data {
				if group == groupName {
					slot := sm.MustNewSlot(table.Head[x], table.Head[x].Add(slotDuration))
					_ = slot.Take(row.Caption) // activityName
					slots = append(slots, slot)
				}
			}
		}
		return slots
	}

	chars := make([]*sm.Character, len(Groups))
	for i, group := range Groups {
		chars[i] = sm.MustNewCharacter(group, strings.ToLower(group), charSlots(group))
	}
	chars[0].Username = "zhikhkirill"

	ctx := context.Background()
	for _, char := range chars {
		if err := repos.Save(ctx, char); err != nil {
			log.Println(err.Error())
		}
	}
}

func importActivities() {
	repos, closeFn := adapters.NewPGActivitiesRepository()
	defer func() {
		err := closeFn()
		if err != nil {
			log.Fatal(err)
		}
	}()

	table := createTable()

	activities := make([]*sm.Activity, 0)
	slotsFor := func(i int) []*sm.Slot {
		slots := make([]*sm.Slot, len(table.Head))
		for j, group := range table.Rows[i].Data {
			slots[j] = sm.MustNewSlot(table.Head[j], table.Head[j].Add(slotDuration))
			if group != "" {
				err := slots[j].Take(group)
				if err != nil {
					log.Println(err.Error())
				}
			}
		}
		return slots
	}

	desc1 := "Центр Молодежной Робототехники - это инновационное пространство, предназначенное для обучения и развития молодых талантов в области  робототехники, БПЛА, искусственного интеллекта и программирования. Их миссия - предоставить молодежи технологическую и информационную базу для развития в областях соревновательной робототехники и не только. Здесь вы можете: обрести команду единомышленников для дальнейшей реализации проектов, востребованных на рынке; создавать свои проекты в области робототехники под руководством опытных специалистов; работать на современном оборудовании: 3D принтеры, лазерные станки, станки с программным управлением и многое другое; улучшить навыки в организации, управлении проектом: начиная с задумки и заканчивая созданием рабочего прототипа; представлять Университет на внешних площадках и инженерных соревнованиях; организовывать собственные инженерные мероприятия и соревнования."
	loc1 := "ИЦАР"
	activities = append(activities, sm.MustNewActivity(
		"ЦМР", &desc1, &loc1,
		[]sm.User{sm.MustNewUser("cmr_admin", sm.Administrator)},
		[]sm.SkillType{sm.Engineering, sm.Researching}, 6,
		slotsFor(0),
	))

	desc2 := "Центр Молодежной Робототехники - это инновационное пространство, предназначенное для обучения и развития молодых талантов в области  робототехники, БПЛА, искусственного интеллекта и программирования. Их миссия - предоставить молодежи технологическую и информационную базу для развития в областях соревновательной робототехники и не только. Здесь вы можете: обрести команду единомышленников для дальнейшей реализации проектов, востребованных на рынке; создавать свои проекты в области робототехники под руководством опытных специалистов; работать на современном оборудовании: 3D принтеры, лазерные станки, станки с программным управлением и многое другое; улучшить навыки в организации, управлении проектом: начиная с задумки и заканчивая созданием рабочего прототипа; представлять Университет на внешних площадках и инженерных соревнованиях; организовывать собственные инженерные мероприятия и соревнования."
	loc2 := "ИЦАР"
	activities = append(activities, sm.MustNewActivity(
		"BRT", &desc2, &loc2,
		[]sm.User{sm.MustNewUser("brt_admin", sm.Administrator)},
		[]sm.SkillType{sm.Engineering, sm.Researching}, 6,
		slotsFor(1),
	))

	desc3 := "Учебно-научный молодежный космический центр (УН МКЦ) основан в 1989 с целью поиска творчески одаренных школьников для  привлечения их на ракетно-космические специальности Университета. В задачи МКЦ входят: разработка и реализация научно-образовательных программ для учащихся, студентов и аспирантов, организация научно-технического творчества молодежи, популяризация достижений космонавтики, развитие и укрепление связей с российскими и международными молодежными организациями. На базе лаборатории перспективных космических технологий в МКЦ выполняются проекты создания реальных образцов космической техники: студенты разрабатывают микроспутники серии «Бауманец», нано- и пикоспутники серии «Парус — МГТУ», полезные нагрузки для выполнения космических экспериментов, специальное программно-математическое обеспечение."
	loc3 := "ХЗ"
	activities = append(activities, sm.MustNewActivity(
		"УН МКЦ", &desc3, &loc3,
		[]sm.User{sm.MustNewUser("mkc_admin", sm.Administrator)},
		[]sm.SkillType{sm.Engineering, sm.Researching}, 6,
		slotsFor(2),
	))

	desc4 := "Учебно-научный молодежный центр «Гидронавтика» (УНМЦ «Гидронавтика») был создан в декабре 2010 г. профессором кафедры «Подводные роботы и аппараты» Станиславом Павловичем Северовым, для внедрения в учебный процесс проектно-конкурентного подхода к инженерному образованию.Основной деятельностью УНМЦ «Гидронавтика» является разработка студенческими инженерными коллективами телеуправляемых необитаемых подводных аппаратов (ТНПА).Основная задача УНМЦ «Гидронавтика» состоит в том, чтобы в течение одного учебного года сформировать команду студентов, которая должна разработать подводный телеуправляемый аппарат, предназначенный для участия в международных соревнованиях Marine Advanced Technology Education (MATE). Каждый член команды получает неоценимый опыт работы с реальными проектами, связь с предыдущими участниками центра и выпускниками МГТУ, а также возможность дальнейшего трудоустройства в рамках проектов «Гидронавтики»."
	loc4 := "ХЗ"
	activities = append(activities, sm.MustNewActivity(
		"Гидронавтика", &desc4, &loc4,
		[]sm.User{sm.MustNewUser("hyd_admin", sm.Administrator)},
		[]sm.SkillType{sm.Engineering, sm.Researching}, 6,
		slotsFor(3),
	))

	desc5 := "Центр НТИ «Цифровое материаловедение: новые материалы и вещества» - структурное подразделение МГТУ им. Н.Э. Баумана, созданное 28 декабря 2020 года для реализации цифрового подхода к «быстрому» и «сквозному» проектированию, разработке, испытанию и применению новых материалов и веществ. Центр НТИ формирует национальный банк данных и знаний по материалам и их «цифровым двойникам», обеспечивающий получение «цифровых паспортов» и ускоренную сертификацию новых материалов. Основная задача Центра НТИ это разработка «Киберполигона цифрового материаловедения» - программно-аппаратного комплекса, обеспечивающего хранение данных о материалах и технологиях, их переработки, компьютерное моделирование материалов и их испытаний, а также системы прогнозирования свойств новых материалов."
	loc5 := "ХЗ"
	activities = append(activities, sm.MustNewActivity(
		"НТИ", &desc5, &loc5,
		[]sm.User{sm.MustNewUser("nti_admin", sm.Administrator)},
		[]sm.SkillType{sm.Engineering, sm.Researching}, 6,
		slotsFor(4),
	))

	desc6 := "Студенческий совет общежития №11 существует уже много лет, и на данный момент в его состав входят 50 активистов, работающих на благо общежития и факультета. Помимо поселения первокурсников и периодических санитарных обходов, благодаря которым поддерживаются уют и порядок, внутри общежития проводится множество мероприятий, направленных на улучшение досуга студентов, их развитие и сплочение. Это всеми любимый костюмированный фестиваль Хэллоуин, кинопоказы, предновогодняя вечеринка и новогодняя ночь, для тех, кто остался в общежитии, Масленица, конкурс на лучшую комнату, кулинарный конкурс. Также раз в семестр проводятся оздоровительные и обучающие выезды в УЦ \"Бауманец\". Студенческий совет общежития №11 часто является объектом внимания СМИ, таких федеральных каналов как Россия-1 и МИР-24."
	loc6 := "ХЗ"
	activities = append(activities, sm.MustNewActivity(
		"ССО №11", &desc6, &loc6,
		[]sm.User{sm.MustNewUser("sso_admin", sm.Administrator)},
		[]sm.SkillType{sm.Social}, 6,
		slotsFor(5),
	))

	desc7 := "Студсовет СМ – это не просто факультетская организация, состоящая из четырех отделов. Это объединение людей, их идей, и, что самое главное, стремлений идти в ногу со временем. Вместе с ребятами ты сможешь учиться и тусоваться, работать и путешествовать, творить и развиваться! В такой компании даже учёба идёт легче: рядом всегда есть товарищи, которые вдохновляют идти только вперёд! Тебе помогут начать свой путь в общественную деятельность и обучим необходимым навыкам! К активистам Студсовета СМ ты всегда можешь обратиться за помощью и поддержкой! Ребята открыты для всех студентов, стараемся поддерживать идеи и инициативы каждого. На протяжении всего времени существования организация уверенно развивается и объединяет активистов разных курсов и кафедр. За 5 лет студенты провели большое количество уникальных мероприятий, и их число продолжает расти. А точкой зарождения новых идей неизменно становится аудитория 509м."
	loc7 := "509м"
	activities = append(activities, sm.MustNewActivity(
		"ССФСМ 1", &desc7, &loc7,
		[]sm.User{sm.MustNewUser("zhikhkirill", sm.Administrator)},
		[]sm.SkillType{sm.Social}, 6,
		slotsFor(6),
	))

	activities = append(activities, sm.MustNewActivity(
		"ССФСМ 2", &desc7, &loc7,
		[]sm.User{sm.MustNewUser("ssf_admin", sm.Administrator)},
		[]sm.SkillType{sm.Social}, 6,
		slotsFor(7),
	))

	desc8 := "Траектория description"
	loc8 := "ХЗ"
	activities = append(activities, sm.MustNewActivity(
		"Траектория", &desc8, &loc8,
		[]sm.User{sm.MustNewUser("trj_admin", sm.Administrator)},
		[]sm.SkillType{sm.Social}, 6,
		slotsFor(8),
	))

	desc9 := "Целью деятельности Профсоюза студентов МГТУ им. Н.Э. Баумана является выражение и защита социальных, экономических и иных законных прав и интересов студентов, а также продвижение студенческих инициатив в стенах университета. В настоящее время Профсоюз студентов МГТУ им. Н.Э. Баумана - это более 26300 членов, более 1500 профоргов, 19 первичных профсоюзных организаций (Профбюро факультетов), 11 комиссий, 9 клубов и 5 добровольческих движения. В рамках оздоровления студентов Профком предлагает путевки в летние лагеря на Черноморское побережье, а также в Подмосковье, где ежегодно отдыхает более 1500 наших студентов. Важным аспектом профсоюза студентов является ведение его кадровой политики, систематическое обучение профоргов. На основе этого ежегодно реализуется конкурс «Лучший профорг». Ежегодно в этом конкурсе принимают участие более 500 профоргов групп нашего университета."
	loc9 := "ХЗ"
	activities = append(activities, sm.MustNewActivity(
		"Профсоюз", &desc9, &loc9,
		[]sm.User{sm.MustNewUser("prf_admin", sm.Administrator)},
		[]sm.SkillType{sm.Social}, 6,
		activitiesSlots(),
	))

	ctx := context.Background()
	for _, a := range activities {
		err := repos.Save(ctx, a)
		if err != nil {
			log.Println(err.Error())
		}
	}
}

type THead []time.Time

type TRow struct {
	Caption string
	Data    []string
}
type Table struct {
	Head THead
	Rows []TRow
}

func createTable() Table {
	return Table{
		Head: THead{
			todayTime(11, 20),
			todayTime(11, 40),
			todayTime(12, 0),
			todayTime(12, 20),
			todayTime(12, 40),
			todayTime(13, 0),
			todayTime(13, 20),
			todayTime(13, 40),
			todayTime(14, 0),
			todayTime(14, 20),
			todayTime(14, 40),
			todayTime(15, 0),
			todayTime(15, 20),
			todayTime(15, 40),
			todayTime(16, 0),
			todayTime(16, 20),
			todayTime(16, 40),
			todayTime(17, 0),
		},
		Rows: []TRow{
			{
				Caption: "ЦМР",
				Data: []string{
					"СМ7-11Б", "СМ7-12Б", "СМ7-13Б", "СМ7-14Б", "СМ11-11Б", "СМ9-13", "СМ5-11Б", "СМ5-12Б",
					"СМ5-11", "СМ13-11Б", "СМ1-12", "СМ3-12", "СМ1-11", "", "", "", "", "",
				},
			},
			{
				Caption: "BRT",
				Data: []string{
					"СМ9-11", "СМ9-12", "СМ10-11Б", "СМ10-11", "СМ10-12", "СМ4-11", "СМ4-12", "СМ4-13",
					"СМ6-19", "СМ2-11", "СМ1-11Б", "СМ13-12Б", "СМ9-13", "СМ8-11", "", "", "", "",
				},
			},
			{
				Caption: "УН МЦК",
				Data: []string{
					"СМ3-11", "СМ1-12", "СМ1-11Б", "СМ2-11", "СМ1-11", "СМ3-12", "СМ8-11", "СМ9-12",
					"СМ6-11", "СМ10-11Б", "СМ9-11", "СМ5-11", "СМ6-12", "СМ6-19", "СМ11-11Б", "СМ4-12", "", "",
				},
			},
			{
				Caption: "Гидронавтика",
				Data: []string{
					"", "", "СМ9-11", "СМ7-12Б", "СМ12-11", "СМ7-14Б", "СМ7-11Б", "СМ11-11Б",
					"СМ10-12", "СМ7-13Б", "СМ10-11", "СМ7-12Б", "СМ5-11Б", "СМ10-12", "СМ13-12Б", "СМ13-13", "", "",
				},
			},
			{
				Caption: "НТИ",
				Data: []string{
					"СМ12-11", "", "СМ13-11Б", "СМ1-11Б", "СМ13-12Б", "СМ13-13", "СМ1-12", "СМ6-12",
					"СМ9-13", "СМ3-11", "СМ4-13", "СМ6-11", "СМ8-11", "СМ5-12Б", "СМ3-12", "СМ4-11", "СМ5-11Б", "",
				},
			},
			{
				Caption: "ССО №11",
				Data: []string{
					"", "", "", "СМ7-11Б", "СМ3-11", "СМ10-11Б", "СМ10-11", "СМ13-12Б",
					"СМ7-14Б", "СМ12-11", "СМ5-12Б", "СМ9-12", "СМ7-13Б", "СМ2-11", "СМ1-11", "СМ5-11", "СМ13-13", "",
				},
			},
			{
				Caption: "ССФСМ 1",
				Data: []string{
					"", "СМ12-11", "СМ1-12", "СМ9-12", "СМ6-11", "СМ1-11", "СМ10-12", "СМ9-11",
					"СМ4-12", "СМ11-11Б", "СМ7-11Б", "", "СМ10-11Б", "СМ7-14Б", "СМ4-11", "СМ8-11", "СМ5-12Б", "СМ5-11",
				},
			},
			{
				Caption: "ССФСМ 2",
				Data: []string{
					"", "СМ3-11", "", "СМ13-11Б", "СМ2-11", "СМ7-13Б", "СМ6-12", "СМ7-12Б",
					"СМ13-13", "СМ5-11Б", "", "СМ13-11Б", "СМ1-11Б", "СМ10-11", "СМ9-13", "СМ3-12", "СМ4-13", "СМ6-19",
				},
			},
			{
				Caption: "Траектория",
				Data: []string{
					"", "СМ6-11", "", "СМ6-12", "", "", "", "СМ4-11",
					"СМ4-13", "СМ4-12", "СМ6-19", "", "", "", "", "", "", "",
				},
			},
		},
	}
}

func activitiesSlots() []*sm.Slot {
	slots := make([]*sm.Slot, 0)
	for start := todayTime(11, 0); start.Before(todayTime(17, 20)); start = start.Add(slotDuration) {
		slot := sm.MustNewSlot(start, start.Add(slotDuration))
		slots = append(slots, slot)
	}
	return slots
}

func characterSlots(start time.Time, end time.Time) []*sm.Slot {
	slots := make([]*sm.Slot, 0)
	for ; start.Before(end); start = start.Add(slotDuration) {
		slot := sm.MustNewSlot(start, start.Add(slotDuration))
		slots = append(slots, slot)
	}
	return slots
}

func todayTime(hours int, minutes int) time.Time {
	return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hours, minutes, 0, 0, time.UTC)
}
