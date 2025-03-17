package repository

import (
        "testing"
        "time"
        "github.com/DATA-DOG/go-sqlmock"
)

func TestGetNoteByUserId(t *testing.T) {
        db, mock, err := sqlmock.New()
        if err != nil {
                t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
        }
        defer db.Close()

        repo := &NoteRepository{DB: db}
        userId := 123

        rows := sqlmock.NewRows([]string{"id", "title", "created_at", "updated_at"}).
                AddRow(1, "Test Note 1", time.Now(), time.Now()).
                AddRow(2, "Test Note 2", time.Now(), time.Now())

        mock.ExpectQuery("SELECT id, title, created_at, updated_at FROM note n WHERE n.active = true AND n.user_id = \\$1").
                WithArgs(userId).
                WillReturnRows(rows)

        notes, err := repo.GetNoteByUserId(userId)
        if err != nil {
                t.Errorf("error was not expected while getting notes: %s", err)
        }

        if notes == nil {
                t.Errorf("expected notes, but got nil")
				return
        }
		
		t.Logf("Notes: %+v", notes)

        if len(notes.Notes) != 2 {
                t.Errorf("expected 2 notes, but got %d", len(notes.Notes))
        }

        if err := mock.ExpectationsWereMet(); err != nil {
                t.Errorf("there were unfulfilled expectations: %s", err)
        }
}