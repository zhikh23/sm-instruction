package adapters

import (
	"context"
	"os"

	"golang.org/x/oauth2/google"
	ss "gopkg.in/Iwark/spreadsheet.v2"

	"github.com/zhikh23/sm-instruction/internal/domain/sm"
)

type gsCharactersProvider struct {
	s ss.Spreadsheet
}

func NewDefaultGSCharactersProvider() sm.CharactersProvider {
	credentialsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS_FILE")
	if credentialsFile == "" {
		panic("GOOGLE_APPLICATION_CREDENTIALS_FILE environment variable is not set")
	}

	spreadsheetID := os.Getenv("GOOGLE_SPREADSHEET_ID")
	if spreadsheetID == "" {
		panic("GOOGLE_SPREADSHEET_ID environment variable is not set")
	}

	return NewGSCharactersProvider(credentialsFile, spreadsheetID)
}

func NewGSCharactersProvider(credentialsFile string, spreadsheetID string) sm.CharactersProvider {
	data, err := os.ReadFile(credentialsFile)
	checkError(err)

	conf, err := google.JWTConfigFromJSON(data, ss.Scope)
	checkError(err)

	client := conf.Client(context.Background())
	service := ss.NewServiceWithClient(client)
	spreadsheet, err := service.FetchSpreadsheet(spreadsheetID)
	checkError(err)

	return &gsCharactersProvider{
		s: spreadsheet,
	}
}

func (p *gsCharactersProvider) Characters(ctx context.Context) ([]*sm.Character, error) {
	sheet, err := p.s.SheetByTitle("EXPORT CHARACTERS")
	if err != nil {
		return nil, err
	}

	chars := make([]*sm.Character, 0)
	for _, row := range sheet.Rows[1:] {
		group := row[0].Value
		username := row[1].Value[1:]
		char := sm.MustNewCharacter(group, username, nil)
		chars = append(chars, char)
	}

	return chars, nil
}
