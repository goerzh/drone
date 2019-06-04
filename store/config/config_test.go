package config

import (
	"context"
	"database/sql"
	"github.com/drone/drone/core"
	"github.com/drone/drone/store/repos"
	"github.com/drone/drone/store/shared/db/dbtest"
	"testing"
)

var noContext = context.TODO()

func TestConfig(t *testing.T) {
	conn, err := dbtest.Connect()
	if err != nil {
		t.Error(err)
		return
	}

	defer func() {
		dbtest.Reset(conn)
		dbtest.Disconnect(conn)
	}()

	// seeds the database with a dummy repository.
	repo := &core.Repository{UID: "1", Slug: "octocat/hello-world"}
	repos := repos.New(conn)
	if err := repos.Create(noContext, repo); err != nil {
		t.Error(err)
	}

	store := New(conn).(*configStore)
	t.Run("Create", testConfigCreate(store, repos, repo))
}

func testConfigCreate(store *configStore, repos core.RepositoryStore, repo *core.Repository) func(t *testing.T) {
	return func(t *testing.T) {
		item := &core.Config{
			RepoID: repo.ID,
			After:  "some commit HASH",
			Data:   "correct-horse-battery-staple",
			Kind:   "demo",
		}
		err := store.Create(noContext, item)
		if err != nil {
			t.Error(err)
		}
		if item.ID == 0 {
			t.Errorf("Want config ID assigned, got %d", item.ID)
		}

		t.Run("Find", testConfigFind(store, item))
		t.Run("FindName", testConfigFindAfter(store, repo))
		t.Run("List", testConfigList(store, repo))
		t.Run("Update", testConfigUpdate(store, repo))
		t.Run("Delete", testConfigDelete(store, repo))
		t.Run("FindOrExist", testConfigFindAfterOrExist(store, repo))
		t.Run("UpdateOrCreate", testConfigUpdateOrCreate(store, repos, repo))
		t.Run("Fkey", testConfigForeignKey(store, repos, repo))
	}
}

func testConfigFind(store *configStore, config *core.Config) func(t *testing.T) {
	return func(t *testing.T) {
		item, err := store.Find(noContext, config.ID)
		if err != nil {
			t.Error(err)
		} else {
			t.Run("Fields", testConfig(item))
		}
	}
}

func testConfigFindAfter(store *configStore, repo *core.Repository) func(t *testing.T) {
	return func(t *testing.T) {
		item, err := store.FindAfter(noContext, repo.ID, "some commit HASH")
		if err != nil {
			t.Error(err)
		} else {
			t.Run("Fields", testConfig(item))
		}
	}
}

func testConfigList(store *configStore, repo *core.Repository) func(t *testing.T) {
	return func(t *testing.T) {
		list, err := store.List(noContext, repo.ID)
		if err != nil {
			t.Error(err)
			return
		}
		if got, want := len(list), 1; got != want {
			t.Errorf("Want count %d, got %d", want, got)
		} else {
			t.Run("Fields", testConfig(list[0]))
		}
	}
}

func testConfigUpdate(store *configStore, repo *core.Repository) func(t *testing.T) {
	return func(t *testing.T) {
		before, err := store.FindAfter(noContext, repo.ID, "some commit HASH")
		if err != nil {
			t.Error(err)
			return
		}
		err = store.Update(noContext, before)
		if err != nil {
			t.Error(err)
			return
		}
		after, err := store.Find(noContext, before.ID)
		if err != nil {
			t.Error(err)
			return
		}
		if after == nil {
			t.Fail()
		}
	}
}

func testConfigDelete(store *configStore, repo *core.Repository) func(t *testing.T) {
	return func(t *testing.T) {
		config, err := store.FindAfter(noContext, repo.ID, "some commit HASH")
		if err != nil {
			t.Error(err)
			return
		}
		err = store.Delete(noContext, config)
		if err != nil {
			t.Error(err)
			return
		}
		_, err = store.Find(noContext, config.ID)
		if got, want := err, sql.ErrNoRows; got != want {
			t.Errorf("Want sql.ErrNoRows, got %v", got)
			return
		}
	}
}

func testConfigFindAfterOrExist(store *configStore, repo *core.Repository) func(t *testing.T) {
	return func(t *testing.T) {
		config, err := store.FindAfterOrExist(noContext, repo.ID, "some commit HASH")
		if err != nil {
			t.Error(err)
			return
		}

		if config != nil {
			t.Errorf("Want %v, got %v", nil, config)
		}
	}
}

func testConfigForeignKey(store *configStore, repos core.RepositoryStore, repo *core.Repository) func(t *testing.T) {
	return func(t *testing.T) {
		item := &core.Config{
			RepoID: repo.ID,
			After:  "some commit HASH",
			Data:   "correct-horse-battery-staple",
			Kind:   "demo",
		}
		store.Create(noContext, item)
		before, _ := store.List(noContext, repo.ID)
		if len(before) == 0 {
			t.Errorf("Want non-empty config list")
			return
		}

		err := repos.Delete(noContext, repo)
		if err != nil {
			t.Error(err)
			return
		}
		after, _ := store.List(noContext, repo.ID)
		if len(after) != 0 {
			t.Errorf("Want empty config list")
		}
	}
}

func testConfigUpdateOrCreate(store *configStore, repos core.RepositoryStore, repo *core.Repository) func(t *testing.T) {
	return func(t *testing.T) {
		item := &core.Config{
			RepoID: repo.ID,
			After:  "some commit HASH",
			Data:   "correct-horse-battery-staple",
			Kind:   "demo",
		}

		err := store.UpdateOrCreate(noContext, item)
		if err != nil {
			t.Error(err)
			return
		}

		if item.ID == 0 {
			t.Errorf("Want config ID assigned, got %d", item.ID)
		}

		err = store.Update(noContext, item)
		if err != nil {
			t.Error(err)
			return
		}
		after, err := store.Find(noContext, item.ID)
		if err != nil {
			t.Error(err)
			return
		}
		if after == nil {
			t.Fail()
		}
	}
}

func testConfig(item *core.Config) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := item.After, "some commit HASH"; got != want {
			t.Errorf("Want config after %q, got %q", want, got)
		}
		if got, want := item.Data, "correct-horse-battery-staple"; got != want {
			t.Errorf("Want config data %q, got %q", want, got)
		}
		if got, want := item.Kind, "demo"; got != want {
			t.Errorf("Want config kind %q, got %q", want, got)
		}
	}
}
