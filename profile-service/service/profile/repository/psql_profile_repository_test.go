package repository

import (
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofrs/uuid"
	"github.com/jariwat/p_project/profile-service/models"
	_profile "github.com/jariwat/p_project/profile-service/service/profile"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ptrUUID() *uuid.UUID {
	id, _ := uuid.NewV4()
	return &id
}

func TestFetchProfiles(t *testing.T) {
	// Create sqlmock database connection and gorm DB instance
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewPsqlProfileRepository(gormDB)

	// Prepare paginator and params
	paginator := &models.Paginator{
		Page:    1,
		PerPage: 2,
	}
	searchTerm := "Phanes"
	params := _profile.GetProfilesParams{
		SearchWord: &searchTerm,
	}

	limit := paginator.PerPage
	var profileID1 = ptrUUID()
	var profileID2 = ptrUUID()

	// Mock count query
	likeQuery := "%" + strings.ToLower(strings.ReplaceAll(searchTerm, " ", "")) + "%"
	countQuery := `SELECT count(*) FROM "profile" WHERE LOWER(REPLACE(CONCAT_WS('', first_name, middle_name, last_name), ' ', '')) LIKE $1`
	mock.ExpectQuery(regexp.QuoteMeta(countQuery)).
		WithArgs(likeQuery).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

	selectQuery := `SELECT * FROM "profile" WHERE LOWER(REPLACE(CONCAT_WS('', first_name, middle_name, last_name), ' ', '')) LIKE $1 LIMIT $2`
	mock.ExpectQuery(regexp.QuoteMeta(selectQuery)).
		WithArgs(likeQuery, limit).
		WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "middle_name", "last_name"}).
			AddRow(profileID1, "SeiA", "F", "Phanes").
			AddRow(profileID2, "AliZe", "", "Phanes"))

	preloadQuery := `SELECT * FROM "skill" WHERE "skill"."profile_id" IN ($1,$2)`
	mock.ExpectQuery(regexp.QuoteMeta(preloadQuery)).
		WithArgs(profileID1, profileID2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "profile_id", "skill", "detail"}).
			AddRow(ptrUUID(), profileID1, "Go", "Advanced").
			AddRow(ptrUUID(), profileID2, "Python", "Intermediate"))

	profiles, err := repo.FetchProfiles(params, paginator)
	assert.NoError(t, err)
	assert.Len(t, profiles, 2)

	assert.Equal(t, 3, paginator.TotalRows)

	// Verify expectations
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestFetchProfileById(t *testing.T) {
	// Setup mock DB
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Prepare repo
	repo := NewPsqlProfileRepository(gormDB)

	// Generate test UUID
	profileID := ptrUUID()

	// Mock profile query
	profileQuery := `SELECT * FROM "profile" WHERE id = $1 ORDER BY "profile"."id" LIMIT $2`
	mock.ExpectQuery(regexp.QuoteMeta(profileQuery)).
		WithArgs(profileID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "first_name"}).
			AddRow(profileID, "SeiA"))

	// Mock preload Skills query
	skillsQuery := `SELECT * FROM "skill" WHERE "skill"."profile_id" = $1`
	mock.ExpectQuery(regexp.QuoteMeta(skillsQuery)).
		WithArgs(profileID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "profile_id", "name"}).
			AddRow(ptrUUID(), profileID, "Go").
			AddRow(ptrUUID(), profileID, "Python"))

	// Call function
	result, err := repo.FetchProfileById(profileID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, profileID, result.ID)
	assert.Equal(t, "SeiA", result.FirstName)

	// Expectation check
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateProfile(t *testing.T) {
	// Setup mock DB
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Create the repository
	repo := NewPsqlProfileRepository(gormDB)

	// Prepare test profile
	profileID := ptrUUID()
	middleName := "F"
	profile := &models.Profile{
		ID:         profileID,
		FirstName:  "AliZe",
		MiddleName: &middleName,
		LastName:   "Phanes",
		Gender:     "FEMALE",
		Class:      "Queen",
		Skills: []*models.Skill{
			{
				ID:        ptrUUID(),
				ProfileID: profileID,
				Skill:     "Go",
				Detail:    "Advanced",
			},
		},
	}
	profile.SetCreatedAt()
	profile.SetUpdatedAt()
	profile.Skills[0].SetCreatedAt()
	profile.Skills[0].SetUpdatedAt()

	// Expect transaction BEGIN
	mock.ExpectBegin()

	// Expect INSERT INTO "profile"
	mock.ExpectExec(`INSERT INTO "profile"`).
		WithArgs(profile.ID, profile.FirstName, profile.MiddleName, profile.LastName, profile.Gender, profile.Class, profile.CreatedAt, profile.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))
	
	// Expect INSERT INTO "skill" for each skill
	for _, skill := range profile.Skills {
		mock.ExpectExec(`INSERT INTO "skill"`).
			WithArgs(skill.ID, skill.ProfileID, skill.Skill, skill.Detail, skill.CreatedAt, skill.UpdatedAt).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}

	// Expect COMMIT
	mock.ExpectCommit()

	// Call function
	err = repo.CreateProfile(profile)
	assert.NoError(t, err)

	// Validate all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateProfile(t *testing.T) {
	// Setup mock DB
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewPsqlProfileRepository(gormDB)

	// Prepare test data
	profileID := ptrUUID()
	middleName := "F"

	profile := &models.Profile{
		ID:         profileID,
		FirstName:  "SeiA",
		MiddleName: &middleName,
		LastName:   "Phanes",
		Gender:     "MALE",
		Class:      "A",
		Skills: []*models.Skill{
			{
				ID:        ptrUUID(),
				ProfileID: profileID,
				Skill:     "Go",
				Detail:    "Advanced",
			},
		},
	}

	profile.SetUpdatedAt()
	profile.Skills[0].SetCreatedAt()
	profile.Skills[0].SetUpdatedAt()


	// Begin transaction
	mock.ExpectBegin()

	// Expect update query with map of columns
	updateQuery := `UPDATE "profile" SET "class"=$1,"first_name"=$2,"gender"=$3,"last_name"=$4,"middle_name"=$5,"updated_at"=$6 WHERE id = $7`
	mock.ExpectExec(regexp.QuoteMeta(updateQuery)).
		WithArgs(
			profile.Class,
			profile.FirstName,
			profile.Gender,
			profile.LastName,
			profile.MiddleName,
			profile.UpdatedAt,
			profileID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Expect delete skills by profile_id
	deleteQuery := `DELETE FROM "skill" WHERE profile_id = $1`
	mock.ExpectExec(regexp.QuoteMeta(deleteQuery)).
		WithArgs(profileID).
		WillReturnResult(sqlmock.NewResult(1, 2)) // assuming 2 rows deleted

	// Expect insert for each skill
	for _, skill := range profile.Skills {
		insertSkillQuery := `INSERT INTO "skill" ("id","profile_id","skill","detail","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6)`
		mock.ExpectExec(regexp.QuoteMeta(insertSkillQuery)).
			WithArgs(sqlmock.AnyArg(), skill.ProfileID, skill.Skill, skill.Detail, skill.CreatedAt, skill.UpdatedAt).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}

	// Commit transaction
	mock.ExpectCommit()

	// Call UpdateProfile
	err = repo.UpdateProfile(profile)
	assert.NoError(t, err)

	// Verify all expectations met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteProfile(t *testing.T) {
	// Setup sqlmock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Wrap sqlmock in gorm DB
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Prepare repo
	repo := NewPsqlProfileRepository(gormDB)

	// Create test UUID
	profileID := ptrUUID()

	// Begin transaction
	mock.ExpectBegin()

	// Expect DELETE FROM "profile" WHERE "id" = $1
	deleteQuery := `DELETE FROM "profile" WHERE "profile"."id" = $1`
	mock.ExpectExec(regexp.QuoteMeta(deleteQuery)).
		WithArgs(profileID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Commit transaction
	mock.ExpectCommit()

	// Call DeleteProfile
	err = repo.DeleteProfile(profileID)
	assert.NoError(t, err)

	// Check all expectations met
	assert.NoError(t, mock.ExpectationsWereMet())
}