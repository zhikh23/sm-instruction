package adapters

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	ss "gopkg.in/Iwark/spreadsheet.v2"

	"github.com/zhikh23/sm-instruction/internal/domain/sm"
)

const slotDuration = 20 * time.Minute

type gsActivitiesProvider struct {
	s ss.Spreadsheet
}

func NewDefaultGSActivitiesProvider() sm.ActivitiesProvider {
	credentialsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS_FILE")
	if credentialsFile == "" {
		panic("GOOGLE_APPLICATION_CREDENTIALS_FILE environment variable is not set")
	}

	spreadsheetID := os.Getenv("GOOGLE_SPREADSHEET_ID")
	if spreadsheetID == "" {
		panic("GOOGLE_SPREADSHEET_ID environment variable is not set")
	}

	return NewGSActivitiesProvider(credentialsFile, spreadsheetID)
}

func NewGSActivitiesProvider(credentialsFile string, spreadsheetID string) sm.ActivitiesProvider {
	data, err := os.ReadFile(credentialsFile)
	checkError(err)

	conf, err := google.JWTConfigFromJSON(data, ss.Scope)
	checkError(err)

	client := conf.Client(context.Background())
	service := ss.NewServiceWithClient(client)
	spreadsheet, err := service.FetchSpreadsheet(spreadsheetID)
	checkError(err)

	return &gsActivitiesProvider{
		s: spreadsheet,
	}
}

func (p *gsActivitiesProvider) Activities(ctx context.Context) ([]*sm.Activity, error) {
	sheet, err := p.s.SheetByTitle("EXPORT ACTIVITIES")
	if err != nil {
		return nil, err
	}

	column := sheet.Columns[1]
	start := 11
	total := 19
	times := make([]time.Time, total)
	for i, cell := range column[start : start+total] {
		t, err := time.Parse(sm.TimeFormat, cell.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to parse time in column 0 row %d: %w", i+5, err)
		}
		times[i] = todayTime(t.Hour(), t.Minute())
	}

	activities := make([]*sm.Activity, 0)
	for _, column := range sheet.Columns[2:] {
		name := column[0].Value

		maxAdmins := 4
		adminsNames := make([]string, 0, maxAdmins)
		for j := 0; j < maxAdmins; j++ {
			adminName := column[2+j].Value
			if adminName == "" {
				break
			}
			adminsNames = append(adminsNames, adminName[1:]) // Remove '@'
		}
		admins := make([]sm.User, len(adminsNames))
		for j, adminName := range adminsNames {
			user, err := sm.NewUser(adminName, sm.Administrator)
			if err != nil {
				return nil, err
			}
			admins[j] = user
		}

		desc := column[6].Value
		location := column[7].Value

		maxSkills := 2
		skills := make([]sm.SkillType, 0, maxSkills)
		for j := 0; j < maxSkills; j++ {
			skillName := column[8+j].Value
			if skillName == "" {
				break
			}
			skill, err := sm.NewSkillTypeFromString(skillName)
			if err != nil {
				return nil, err
			}
			skills = append(skills, skill)
		}

		maxPoints := 0
		maxPointStr := column[10].Value
		if maxPointStr != "" {
			maxPoints, err = strconv.Atoi(maxPointStr)
			if err != nil {
				return nil, fmt.Errorf("failed to parse max points: %w", err)
			}
		}

		slots := make([]*sm.Slot, 0, total)
		for j := 0; j < total; j++ {
			startTime := times[j]
			groupName := column[start+j].Value
			slot, err := sm.NewSlot(startTime, startTime.Add(slotDuration))
			if err != nil {
				return nil, err
			}
			if groupName != "" {
				err = slot.Take(groupName)
				if err != nil {
					return nil, err
				}
			}
			slots = append(slots, slot)
		}

		activity, err := sm.NewActivity(
			name,
			pointerIfNotEmpty(desc),
			pointerIfNotEmpty(location),
			admins,
			skills,
			maxPoints,
			slots,
		)
		if err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}

	return activities, nil
}

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func todayTime(hours int, minutes int) time.Time {
	return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hours, minutes, 0, 0, time.Local)
}

func pointerIfNotEmpty(s string) *string {
	if s != "" {
		return &s
	}
	return nil
}
