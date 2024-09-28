package adapters

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/zhikh23/pgutils"

	"github.com/zhikh23/sm-instruction/internal/domain/sm"
)

type pgCharactersRepository struct {
	db *sqlx.DB
}

func NewPGCharactersRepository() (sm.CharactersRepository, func() error) {
	uri := os.Getenv("DATABASE_URL")
	if uri == "" {
		panic("DATABASE_URL environment variable not set")
	}
	db := sqlx.MustConnect("postgres", uri)

	return &pgCharactersRepository{db: db}, db.Close
}

func (r *pgCharactersRepository) Save(
	ctx context.Context,
	character *sm.Character,
) error {
	if err := pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		return r.save(ctx, tx, character)
	}); pgutils.IsUniqueViolationError(err) {
		return sm.ErrCharacterAlreadyExists
	} else if err != nil {
		return err
	}
	return nil
}

func (r *pgCharactersRepository) Character(
	ctx context.Context,
	groupName string,
) (*sm.Character, error) {
	var char *sm.Character
	var err error
	if err = pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		char, err = r.character(ctx, tx, groupName)
		return err
	}); errors.Is(err, sql.ErrNoRows) {
		return nil, sm.ErrCharacterNotFound
	} else if err != nil {
		return nil, err
	}
	return char, nil
}

func (r *pgCharactersRepository) CharacterByUsername(
	ctx context.Context,
	username string,
) (*sm.Character, error) {
	var char *sm.Character
	var err error
	if err = pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		char, err = r.characterByUsername(ctx, tx, username)
		return err
	}); errors.Is(err, sql.ErrNoRows) {
		return nil, sm.ErrCharacterNotFound
	} else if err != nil {
		return nil, err
	}
	return char, nil
}

func (r *pgCharactersRepository) Update(
	ctx context.Context,
	groupName string,
	updateFn func(innerCtx context.Context, char *sm.Character) error,
) error {
	return pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		char, err := r.character(ctx, tx, groupName)
		if errors.Is(err, sql.ErrNoRows) {
			return sm.ErrCharacterNotFound
		} else if err != nil {
			return err
		}

		err = updateFn(ctx, char)
		if err != nil {
			return err
		}

		return r.update(ctx, tx, char)
	})
}

func (r *pgCharactersRepository) save(
	ctx context.Context,
	ex sqlx.ExtContext,
	character *sm.Character,
) error {
	var err error
	if err = r.requireExecResult(sqlx.NamedExecContext(ctx, ex,
		`INSERT INTO
			characters (group_name, username, started_at) 
		 VALUES (:group_name, :username, :started_at)`,
		marshallCharacterToRow(character),
	)); err != nil {
		return err
	}

	if err = r.requireExecResult(sqlx.NamedExecContext(ctx, ex,
		`INSERT INTO
			character_skills (group_name, skill_type, points)
		 VALUES (:group_name, :skill_type, :points)`,
		marshallCharacterSkillsToRows(character.GroupName, character.Skills),
	)); err != nil {
		return err
	}

	if err = r.requireExecResult(sqlx.NamedExecContext(ctx, ex,
		`INSERT INTO
			character_slots (group_name, start, end_, activity_name) 
		 VALUES (:group_name, :start, :end_, :activity_name)`,
		marshallCharacterSlotsToRows(character.GroupName, character.Slots),
	)); err != nil {
		return err
	}

	return nil
}

func (r *pgCharactersRepository) character(
	ctx context.Context,
	qx sqlx.QueryerContext,
	groupName string,
) (*sm.Character, error) {
	var err error

	var characterRow characterRow
	if err = sqlx.GetContext(ctx, qx, &characterRow,
		`SELECT group_name, username, started_at
   		 FROM characters
		 WHERE group_name = $1`, groupName,
	); err != nil {
		return nil, err
	}

	var characterSkillsRows []characterSkillRow
	if err = sqlx.SelectContext(ctx, qx, &characterSkillsRows,
		`SELECT group_name, skill_type, points
		 FROM   character_skills
		 WHERE  group_name = $1`, characterRow.GroupName,
	); err != nil {
		return nil, err
	}
	skills, err := unmarshallCharacterSkillsFromRows(characterSkillsRows)
	if err != nil {
		return nil, err
	}

	var characterSlotsRows []characterSlotRow
	if err = sqlx.SelectContext(ctx, qx, &characterSlotsRows,
		`SELECT   group_name, start, end_, activity_name
		 FROM     character_slots
		 WHERE    group_name = $1
		 ORDER BY start`, characterRow.GroupName,
	); err != nil {
		return nil, err
	}
	slots, err := unmarshallCharacterSlotsFromRows(characterSlotsRows)
	if err != nil {
		return nil, err
	}

	return sm.UnmarshallCharacterFromDB(
		characterRow.GroupName,
		characterRow.Username,
		skills,
		timeLocalOrNil(characterRow.StartedAt),
		slots,
	)
}

func (r *pgCharactersRepository) characterByUsername(
	ctx context.Context,
	qx sqlx.QueryerContext,
	username string,
) (*sm.Character, error) {
	var err error

	var characterRow characterRow
	if err = sqlx.GetContext(ctx, qx, &characterRow,
		`SELECT group_name, username, started_at
   		 FROM   characters
		 WHERE  username = $1`, username,
	); err != nil {
		return nil, err
	}

	var characterSkillsRows []characterSkillRow
	if err = sqlx.SelectContext(ctx, qx, &characterSkillsRows,
		`SELECT group_name, skill_type, points
		 FROM   character_skills
		 WHERE  group_name = $1`, characterRow.GroupName,
	); err != nil {
		return nil, err
	}
	skills, err := unmarshallCharacterSkillsFromRows(characterSkillsRows)
	if err != nil {
		return nil, err
	}

	var characterSlotsRows []characterSlotRow
	if err = sqlx.SelectContext(ctx, qx, &characterSlotsRows,
		`SELECT   group_name, start, end_, activity_name
		 FROM     character_slots
		 WHERE    group_name = $1
		 ORDER BY start`, characterRow.GroupName,
	); err != nil {
		return nil, err
	}
	slots, err := unmarshallCharacterSlotsFromRows(characterSlotsRows)
	if err != nil {
		return nil, err
	}

	return sm.UnmarshallCharacterFromDB(
		characterRow.GroupName,
		characterRow.Username,
		skills,
		timeLocalOrNil(characterRow.StartedAt),
		slots,
	)
}

func (r *pgCharactersRepository) update(
	ctx context.Context,
	ex sqlx.ExtContext,
	character *sm.Character,
) error {
	var err error
	if err = r.requireExecResult(ex.ExecContext(ctx,
		`UPDATE characters 
		 SET    started_at = $2
		 WHERE  group_name = $1`, character.GroupName, character.StartedAt,
	)); err != nil {
		return err
	}

	if err = r.requireExecResult(ex.ExecContext(ctx,
		`DELETE FROM character_skills WHERE group_name = $1`, character.GroupName,
	)); err != nil {
		return err
	}

	if err = r.requireExecResult(sqlx.NamedExecContext(ctx, ex,
		`INSERT INTO
			character_skills (group_name, skill_type, points)
		 VALUES (:group_name, :skill_type, :points)`,
		marshallCharacterSkillsToRows(character.GroupName, character.Skills),
	)); err != nil {
		return err
	}

	if err = r.requireExecResult(ex.ExecContext(ctx,
		`DELETE FROM character_slots WHERE group_name = $1`, character.GroupName,
	)); err != nil {
		return err
	}

	if err = r.requireExecResult(sqlx.NamedExecContext(ctx, ex,
		`INSERT INTO
			character_slots (group_name, start, end_, activity_name) 
		 VALUES (:group_name, :start, :end_, :activity_name)`,
		marshallCharacterSlotsToRows(character.GroupName, character.Slots),
	)); err != nil {
		return err
	}

	return nil
}

func (r *pgCharactersRepository) requireExecResult(res sql.Result, err error) error {
	if err != nil {
		return err
	}

	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if aff == 0 {
		return sql.ErrNoRows
	}

	return nil
}

type characterRow struct {
	GroupName string     `db:"group_name"`
	Username  string     `db:"username"`
	StartedAt *time.Time `db:"started_at"`
}

func marshallCharacterToRow(c *sm.Character) characterRow {
	return characterRow{
		GroupName: c.GroupName,
		Username:  c.Username,
		StartedAt: timeUTCOrNil(c.StartedAt),
	}
}

type characterSkillRow struct {
	GroupName string `db:"group_name"`
	SkillType string `db:"skill_type"`
	Points    int    `db:"points"`
}

func marshallCharacterSkillToRow(groupName string, skillType sm.SkillType, points int) characterSkillRow {
	return characterSkillRow{
		GroupName: groupName,
		SkillType: skillType.String(),
		Points:    points,
	}
}

func marshallCharacterSkillsToRows(groupName string, skills map[sm.SkillType]int) []characterSkillRow {
	res := make([]characterSkillRow, 0, len(skills))
	for k, v := range skills {
		res = append(res, marshallCharacterSkillToRow(groupName, k, v))
	}
	return res
}

func unmarshallCharacterSkillsFromRows(ss []characterSkillRow) (map[sm.SkillType]int, error) {
	res := make(map[sm.SkillType]int, len(ss))
	for _, v := range ss {
		st, err := sm.NewSkillTypeFromString(v.SkillType)
		if err != nil {
			return nil, err
		}
		res[st] = v.Points
	}
	return res, nil
}

type characterSlotRow struct {
	GroupName    string    `db:"group_name"`
	Start        time.Time `db:"start"`
	End          time.Time `db:"end_"`
	ActivityName *string   `db:"activity_name"`
}

func marshallCharacterSlotToRow(groupName string, s *sm.Slot) characterSlotRow {
	return characterSlotRow{
		GroupName:    groupName,
		Start:        s.Start.UTC(),
		End:          s.End.UTC(),
		ActivityName: s.Whom,
	}
}

func marshallCharacterSlotsToRows(groupName string, ss []*sm.Slot) []characterSlotRow {
	res := make([]characterSlotRow, len(ss))
	for i, s := range ss {
		res[i] = marshallCharacterSlotToRow(groupName, s)
	}
	return res
}

func unmarshallCharacterSlotFromRow(a characterSlotRow) (*sm.Slot, error) {
	return sm.UnmarshallSlotFromDB(a.Start.Local(), a.End.Local(), a.ActivityName)
}

func unmarshallCharacterSlotsFromRows(cs []characterSlotRow) ([]*sm.Slot, error) {
	res := make([]*sm.Slot, len(cs))
	for i, s := range cs {
		slot, err := unmarshallCharacterSlotFromRow(s)
		if err != nil {
			return nil, err
		}
		res[i] = slot
	}
	return res, nil
}

func timeUTCOrNil(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	v := t.UTC()
	return &v
}

func timeLocalOrNil(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	v := t.Local()
	return &v
}
