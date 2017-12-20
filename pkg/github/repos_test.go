package github

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/sourcegraph/go-github/github"
	"sourcegraph.com/sourcegraph/sourcegraph/pkg/actor"
	"sourcegraph.com/sourcegraph/sourcegraph/pkg/api/legacyerr"
	"sourcegraph.com/sourcegraph/sourcegraph/pkg/rcache"
)

type RoundTripperFunc func(*http.Request) (*http.Response, error)

func (f RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

const testGHCachePrefix = "__test__gh_pub"

func resetCache(t *testing.T) {
	rcache.SetupForTest(testGHCachePrefix)
	reposGithubPublicCache = rcache.NewWithTTL(testGHCachePrefix, 1000)
}

// TestRepos_Get_nocache tests the behavior of Repos.Get.
func TestRepos_Get(t *testing.T) {
	resetCache(t)

	var calledGet bool
	MockRoundTripper = RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		calledGet = true
		body, err := json.Marshal(&github.Repository{
			ID:       github.Int(123),
			Name:     github.String("repo"),
			FullName: github.String("owner/repo"),
			Owner:    &github.User{ID: github.Int(1)},
			CloneURL: github.String("https://github.com/owner/repo.git"),
		})
		if err != nil {
			t.Fatal(err)
		}
		return &http.Response{
			Request:    req,
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(body)),
		}, nil
	})

	repo, err := GetRepo(context.Background(), "github.com/owner/repo")
	if err != nil {
		t.Fatal(err)
	}
	if repo == nil {
		t.Error("repo == nil")
	}
	if !calledGet {
		t.Error("!calledGet, expected to miss cache")
	}

	// Test that repo is cached (and therefore NOT fetched) from client on second request.
	calledGet = false
	repo, err = GetRepo(context.Background(), "github.com/owner/repo")
	if err != nil {
		t.Fatal(err)
	}
	if repo == nil {
		t.Error("repo == nil")
	}
	if calledGet {
		t.Error("calledGet, expected to hit cache")
	}
}

// TestRepos_Get_nonexistent tests the behavior of Repos.Get when called
// on a repo that does not exist.
func TestRepos_Get_nonexistent(t *testing.T) {
	resetCache(t)

	MockRoundTripper = RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			Request:    req,
			StatusCode: http.StatusNotFound,
			Body:       ioutil.NopCloser(bytes.NewReader(nil)),
		}, nil
	})

	nonexistentRepo := "github.com/owner/repo"
	repo, err := GetRepo(context.Background(), nonexistentRepo)
	if legacyerr.ErrCode(err) != legacyerr.NotFound {
		t.Fatal(err)
	}
	if repo != nil {
		t.Error("repo != nil")
	}
}

// TestRepos_Get_publicnotfound tests we correctly cache 404 responses. If we
// are not authed as a user, all private repos 404 which counts towards our
// rate limit. This test will ensure unauthed clients cache 404, but authed
// users do not use the cache
func TestRepos_Get_publicnotfound(t *testing.T) {
	resetCache(t)

	calledGetMissing := false
	mockGetMissing := RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		calledGetMissing = true
		return &http.Response{
			Request:    req,
			StatusCode: http.StatusNotFound,
			Body:       ioutil.NopCloser(bytes.NewReader(nil)),
		}, nil
	})

	privateRepo := "github.com/owner/repo"

	// An unauthed user won't be able to see the repo
	MockRoundTripper = mockGetMissing
	ctx := actor.WithActor(context.Background(), &actor.Actor{})
	if _, err := GetRepo(ctx, privateRepo); legacyerr.ErrCode(err) != legacyerr.NotFound {
		t.Fatal(err)
	}
	if !calledGetMissing {
		t.Fatal("did not call repos.Get when it should not be cached")
	}

	// If we repeat the call, we should hit the cache
	calledGetMissing = false
	if _, err := GetRepo(ctx, privateRepo); legacyerr.ErrCode(err) != legacyerr.NotFound {
		t.Fatal(err)
	}
	if calledGetMissing {
		t.Fatal("should have hit cache")
	}

	// Ensure the repo is still missing for unauthed users
	calledGetMissing = false
	MockRoundTripper = mockGetMissing
	ctx = actor.WithActor(context.Background(), &actor.Actor{})
	if _, err := GetRepo(ctx, privateRepo); legacyerr.ErrCode(err) != legacyerr.NotFound {
		t.Fatal(err)
	}
	if calledGetMissing {
		t.Fatal("should have hit cache")
	}
}

// TestRepos_Get_authednocache tests authed users do add public repos to the cache
func TestRepos_Get_authednocache(t *testing.T) {
	resetCache(t)

	calledGet := false
	MockRoundTripper = RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		calledGet = true
		body, err := json.Marshal(&github.Repository{
			ID:       github.Int(123),
			Name:     github.String("repo"),
			FullName: github.String("owner/repo"),
			Owner:    &github.User{ID: github.Int(1)},
			CloneURL: github.String("https://github.com/owner/repo.git"),
			Private:  github.Bool(false),
		})
		if err != nil {
			t.Fatal(err)
		}
		return &http.Response{
			Request:    req,
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(body)),
		}, nil
	})

	repo := "github.com/owner/repo"

	authedGet := func() bool {
		calledGet = false
		ctx := actor.WithActor(context.Background(), &actor.Actor{UID: "1", Login: "test"})
		_, err := GetRepo(ctx, repo)
		if err != nil {
			t.Fatal(err)
		}
		return calledGet
	}
	unauthedGet := func() bool {
		calledGet = false
		ctx := actor.WithActor(context.Background(), &actor.Actor{})
		_, err := GetRepo(ctx, repo)
		if err != nil {
			t.Fatal(err)
		}
		return calledGet
	}

	// An authed user will populate the empty cache
	if !authedGet() {
		t.Fatal("authed should populate empty cache")
	}

	// An unauthed user should now get from cache
	if unauthedGet() {
		t.Fatal("unauthed should get from cache")
	}

	// The authed user should also be able to get public repo from the cache
	if authedGet() {
		t.Fatal("authed should not get from cache")
	}
}
