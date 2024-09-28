package adapters

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/zhikh23/pgutils"

	"github.com/zhikh23/sm-instruction/internal/domain/sm"
)

type pgActivitiesRepository struct {
	db *sqlx.DB
}

func NewPGActivitiesRepository() (sm.ActivitiesRepository, func() error) {
	uri := os.Getenv("DATABASE_URL")
	if uri == "" {
		panic("DATABASE_URL environment variable not set")
	}
	db := sqlx.MustConnect("postgres", uri)

	return &pgActivitiesRepository{db: db}, db.Close
}

func (r *pgActivitiesRepository) Save(
	ctx context.Context,
	activity *sm.Activity,
) error {
	return pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		var err error
		if err = r.requireExecResult(tx.NamedExecContext(ctx,
			`INSERT INTO
					activities (name, description, location, skills, max_points)
			 VALUES (:name, :description, :location, :skills, :max_points)`,
			marshallActivityToRow(activity),
		)); pgutils.IsUniqueViolationError(err) {
			return sm.ErrActivityAlreadyExists
		} else if err != nil {
			return err
		}

		if len(activity.Admins) > 0 {
			if err = r.requireExecResult(tx.NamedExecContext(ctx,
				`INSERT INTO
				admins (activity_name, username)
			 VALUES (:activity_name, :username)`,
				marshallAdminsToRows(activity.Name, activity.Admins),
			)); err != nil {
				return err
			}
		}

		if len(activity.Slots()) > 0 {
			if err = r.requireExecResult(tx.NamedExecContext(ctx,
				`INSERT INTO
					activity_slots (activity_name, start, end_, group_name) 
			 	 VALUES (:activity_name, :start, :end_, :group_name)`,
				marshallActivitySlotsToRows(activity.Name, activity.Slots()),
			)); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *pgActivitiesRepository) Activity(
	ctx context.Context,
	activityName string,
) (*sm.Activity, error) {
	var res *sm.Activity
	var err error
	if err = pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		res, err = r.activity(ctx, tx, activityName)
		return err
	}); errors.Is(err, sql.ErrNoRows) {
		return nil, sm.ErrActivityNotFound
	} else if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *pgActivitiesRepository) ActivityByAdmin(
	ctx context.Context,
	adminUsername string,
) (*sm.Activity, error) {
	var res *sm.Activity
	var err error
	if err = pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		res, err = r.activityByAdmin(ctx, tx, adminUsername)
		return err
	}); errors.Is(err, sql.ErrNoRows) {
		return nil, sm.ErrActivityNotFound
	} else if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *pgActivitiesRepository) AvailableActivities(
	ctx context.Context,
) ([]*sm.Activity, error) {
	var res []*sm.Activity
	var err error
	if err = pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		res, err = r.availableActivities(ctx, tx)
		return err
	}); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *pgActivitiesRepository) UpdateSlots(
	ctx context.Context,
	activityUUID string,
	updateFn func(innerCtx context.Context, activity *sm.Activity) error,
) error {
	return pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		activity, err := r.activity(ctx, tx, activityUUID)
		if errors.Is(err, sql.ErrNoRows) {
			return sm.ErrActivityNotFound
		} else if err != nil {
			return err
		}

		err = updateFn(ctx, activity)
		if err != nil {
			return err
		}

		return r.updateSlots(ctx, tx, activity)
	})
}

func (r *pgActivitiesRepository) activity(
	ctx context.Context,
	qx sqlx.QueryerContext,
	activityName string,
) (*sm.Activity, error) {
	var err error

	var activityRow activityRow
	if err = sqlx.GetContext(ctx, qx, &activityRow,
		`SELECT name, description, location, skills, max_points 
		 FROM   activities
		 WHERE  name = $1`, activityName,
	); err != nil {
		return nil, err
	}

	var adminsRows []adminRow
	if err = sqlx.SelectContext(ctx, qx, &adminsRows,
		`SELECT activity_name, username 
		 FROM admins 
		 WHERE activity_name = $1`, activityRow.Name,
	); err != nil {
		return nil, err
	}
	admins, err := unmarshallAdminsFromRows(adminsRows)
	if err != nil {
		return nil, err
	}

	var slotsRows []activitySlotRow
	if err = sqlx.SelectContext(ctx, qx, &slotsRows,
		`SELECT activity_name, start, end_, group_name
		 FROM activity_slots
		 WHERE activity_name = $1
		 ORDER BY start`, activityRow.Name,
	); err != nil {
		return nil, err
	}
	slots, err := unmarshallActivitySlotsFromRows(slotsRows)
	if err != nil {
		return nil, err
	}

	return sm.UnmarshallActivityFromDB(
		activityRow.Name,
		activityRow.Description,
		activityRow.Location,
		admins,
		activityRow.Skills,
		activityRow.MaxPoints,
		slots,
	)
}

func (r *pgActivitiesRepository) activityByAdmin(
	ctx context.Context,
	qx sqlx.QueryerContext,
	adminUsername string,
) (*sm.Activity, error) {
	var err error

	var activityRow activityRow
	if err = sqlx.GetContext(ctx, qx, &activityRow,
		`SELECT name, description, location, skills, max_points 
		 FROM   activities AS activity
				LEFT JOIN admins AS admin 
					   ON admin.activity_name = activity.name
		 WHERE admin.username = $1`, adminUsername,
	); err != nil {
		return nil, err
	}

	var adminsRows []adminRow
	if err = sqlx.SelectContext(ctx, qx, &adminsRows,
		`SELECT activity_name, username 
		 FROM admins 
		 WHERE activity_name = $1`, activityRow.Name,
	); err != nil {
		return nil, err
	}
	admins, err := unmarshallAdminsFromRows(adminsRows)
	if err != nil {
		return nil, err
	}

	var slotsRows []activitySlotRow
	if err = sqlx.SelectContext(ctx, qx, &slotsRows,
		`SELECT activity_name, start, end_, group_name
		 FROM activity_slots
		 WHERE activity_name = $1`, activityRow.Name,
	); err != nil {
		return nil, err
	}
	slots, err := unmarshallActivitySlotsFromRows(slotsRows)
	if err != nil {
		return nil, err
	}

	return sm.UnmarshallActivityFromDB(
		activityRow.Name,
		activityRow.Description,
		activityRow.Location,
		admins,
		activityRow.Skills,
		activityRow.MaxPoints,
		slots,
	)
}

func (r *pgActivitiesRepository) availableActivities(
	ctx context.Context,
	qx sqlx.QueryerContext,
) ([]*sm.Activity, error) {
	var err error

	var activityRows []activityRow
	if err = sqlx.SelectContext(ctx, qx, &activityRows,
		`SELECT name, description, location, skills, max_points 
		 FROM   activities AS activity
		 WHERE  location IS NOT NULL`,
	); err != nil {
		return nil, err
	}

	activities := make([]*sm.Activity, 0, len(activityRows))
	for _, activityRow := range activityRows {
		var slotsRows []activitySlotRow
		if err = sqlx.SelectContext(ctx, qx, &slotsRows,
			`SELECT activity_name, start, end_, group_name
			 FROM activity_slots
			 WHERE activity_name = $1
			 ORDER BY start`, activityRow.Name,
		); err != nil {
			return nil, err
		}
		if len(slotsRows) == 0 {
			continue
		}

		slots, err := unmarshallActivitySlotsFromRows(slotsRows)
		if err != nil {
			return nil, err
		}

		var adminsRows []adminRow
		if err = sqlx.SelectContext(ctx, qx, &adminsRows,
			`SELECT activity_name, username 
			 FROM admins 
			 WHERE activity_name = $1`, activityRow.Name,
		); err != nil {
			return nil, err
		}
		admins, err := unmarshallAdminsFromRows(adminsRows)
		if err != nil {
			return nil, err
		}

		activity, err := sm.UnmarshallActivityFromDB(
			activityRow.Name,
			activityRow.Description,
			activityRow.Location,
			admins,
			activityRow.Skills,
			activityRow.MaxPoints,
			slots,
		)
		if err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}
	return activities, nil
}

func (r *pgActivitiesRepository) updateSlots(
	ctx context.Context,
	ex sqlx.ExecerContext,
	activity *sm.Activity,
) error {
	var err error
	for _, slot := range activity.Slots() {
		if err = r.requireExecResult(ex.ExecContext(ctx,
			`UPDATE activity_slots SET group_name = $3 WHERE activity_name = $1 AND start = $2`,
			activity.Name, slot.Start.UTC(), slot.Whom,
		)); err != nil {
			return err
		}
	}
	return nil
}

func (r *pgActivitiesRepository) requireExecResult(res sql.Result, err error) error {
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

type activityRow struct {
	Name        string         `db:"name"`
	Description *string        `db:"description"`
	Location    *string        `db:"location"`
	Skills      pq.StringArray `db:"skills"`
	MaxPoints   int            `db:"max_points"`
}

func marshallActivityToRow(a *sm.Activity) activityRow {
	pqSkills := make(pq.StringArray, len(a.Skills))
	for i, s := range a.Skills {
		pqSkills[i] = s.String()
	}
	return activityRow{
		Name:        a.Name,
		Description: a.Description,
		Location:    a.Location,
		Skills:      pqSkills,
		MaxPoints:   a.MaxPoints,
	}
}

type adminRow struct {
	ActivityName string `db:"activity_name"`
	Username     string `db:"username"`
}

func marshallAdminToRow(activityName string, a sm.User) adminRow {
	return adminRow{
		ActivityName: activityName,
		Username:     a.Username,
	}
}

func unmarshallAdminFromRow(a adminRow) (sm.User, error) {
	return sm.UnmarshallUserFromDB(a.Username, sm.Administrator.String())
}

func unmarshallAdminsFromRows(as []adminRow) ([]sm.User, error) {
	res := make([]sm.User, len(as))
	for i, a := range as {
		user, err := unmarshallAdminFromRow(a)
		if err != nil {
			return nil, err
		}
		res[i] = user
	}
	return res, nil
}

func marshallAdminsToRows(activityName string, as []sm.User) []adminRow {
	res := make([]adminRow, len(as))
	for i, a := range as {
		res[i] = marshallAdminToRow(activityName, a)
	}
	return res
}

type activitySlotRow struct {
	ActivityName string    `db:"activity_name"`
	Start        time.Time `db:"start"`
	End          time.Time `db:"end_"`
	GroupName    *string   `db:"group_name"`
}

func marshallActivitySlotToRow(activityName string, s *sm.Slot) activitySlotRow {
	return activitySlotRow{
		ActivityName: activityName,
		Start:        s.Start.UTC(),
		End:          s.End.UTC(),
		GroupName:    s.Whom,
	}
}

func marshallActivitySlotsToRows(activityName string, ss []*sm.Slot) []activitySlotRow {
	res := make([]activitySlotRow, len(ss))
	for i, s := range ss {
		res[i] = marshallActivitySlotToRow(activityName, s)
	}
	return res
}

func unmarshallActivitySlotFromRow(a activitySlotRow) (*sm.Slot, error) {
	return sm.UnmarshallSlotFromDB(a.Start.Local(), a.End.Local(), a.GroupName)
}

func unmarshallActivitySlotsFromRows(as []activitySlotRow) ([]*sm.Slot, error) {
	res := make([]*sm.Slot, len(as))
	for i, a := range as {
		slot, err := unmarshallActivitySlotFromRow(a)
		if err != nil {
			return nil, err
		}
		res[i] = slot
	}
	return res, nil
}
