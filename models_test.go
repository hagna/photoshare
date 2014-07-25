package photoshare

import (
	"database/sql"
	"testing"
)

func TestGetIfNotNone(t *testing.T) {

	config, _ := newAppConfig()
	tdb := makeTestDB(config)
	defer tdb.clean()

	userDS := newUserDataStore(tdb.dbMap)
	photoDS := newPhotoDataStore(tdb.dbMap)

	user := &user{Name: "tester", Email: "tester@gmail.com", Password: "test"}

	if err := userDS.insert(user); err != nil {
		t.Error(err)
		return
	}
	photo := &photo{Title: "test", OwnerID: user.ID, Filename: "test.jpg"}
	if err := photoDS.insert(photo); err != nil {
		t.Error(err)
		return
	}

	photo, err := photoDS.get(photo.ID)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetIfNone(t *testing.T) {

	config, _ := newAppConfig()
	tdb := makeTestDB(config)
	defer tdb.clean()

	_, err := newPhotoDataStore(tdb.dbMap).get(1)
	if err != sql.ErrNoRows {
		t.Error(err)
		return
	}

}

func TestSearchPhotos(t *testing.T) {
	config, _ := newAppConfig()
	tdb := makeTestDB(config)
	defer tdb.clean()

	photoDS := newPhotoDataStore(tdb.dbMap)
	userDS := newUserDataStore(tdb.dbMap)

	user := &user{Name: "tester", Email: "tester@gmail.com", Password: "test"}
	if err := userDS.insert(user); err != nil {
		t.Error(err)
		return
	}
	photo := &photo{Title: "test", OwnerID: user.ID, Filename: "test.jpg"}
	if err := photoDS.insert(photo); err != nil {
		t.Error(err)
		return
	}
	result, err := photoDS.search(newPage(1), "test")
	if err != nil {
		t.Error(err)
		return
	}

	if len(result.Items) != 1 {
		t.Error("There should be 1 photo")
	}
}
func TestAllPhotos(t *testing.T) {
	config, _ := newAppConfig()
	tdb := makeTestDB(config)
	defer tdb.clean()

	photoDS := newPhotoDataStore(tdb.dbMap)
	userDS := newUserDataStore(tdb.dbMap)

	user := &user{Name: "tester", Email: "tester@gmail.com", Password: "test"}
	if err := userDS.insert(user); err != nil {
		t.Error(err)
		return
	}
	photo := &photo{Title: "test", OwnerID: user.ID, Filename: "test.jpg"}
	if err := photoDS.insert(photo); err != nil {
		t.Error(err)
		return
	}
	result, err := photoDS.all(newPage(1), "")
	if err != nil {
		t.Error(err)
		return
	}

	if len(result.Items) != 1 {
		t.Error("There should be 1 photo")
	}
}

func TestCanEdit(t *testing.T) {
	user := &user{ID: 1}
	photo := &photo{ID: 1, OwnerID: 1}

	if photo.canEdit(user) {
		t.Error("Non-authenticated should not be able to edit")
	}

	user.IsAuthenticated = true

	if !photo.canEdit(user) {
		t.Error("User should be able to edit")
	}

	photo.OwnerID = 2

	if photo.canEdit(user) {
		t.Error("User should not be able to edit")
	}

	user.IsAdmin = true
	if !photo.canEdit(user) {
		t.Error("Admin should be able to edit")
	}
}

func TestHasVoted(t *testing.T) {

	u := &user{}
	if u.hasVoted(1) {
		t.Error("The user has not voted yet")
	}

	u.registerVote(1)
	if !u.hasVoted(1) {
		t.Error("The user should have voted")
	}
}