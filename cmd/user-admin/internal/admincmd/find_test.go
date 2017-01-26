// Copyright 2016 Canonical Ltd.

package admincmd_test

import (
	"encoding/json"
	"time"

	"github.com/juju/idmclient/params"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
)

type findSuite struct {
	commandSuite
}

var _ = gc.Suite(&findSuite{})

func (s *findSuite) TestFindEmail(c *gc.C) {
	runf := s.RunServer(c, &handler{
		queryUsers: func(req *params.QueryUsersRequest) ([]string, error) {
			if req.Email == "bob@example.com" {
				return []string{"bob"}, nil
			}
			return []string{}, nil
		},
	})
	stdout := CheckSuccess(c, runf, "find", "-a", "admin.agent", "-e", "bob@example.com")
	c.Assert(stdout, gc.Equals, "bob\n")
}

func (s *findSuite) TestFindEmailNotFound(c *gc.C) {
	runf := s.RunServer(c, &handler{
		queryUsers: func(req *params.QueryUsersRequest) ([]string, error) {
			if req.Email == "alice@example.com" {
				return []string{"alice"}, nil
			}
			return []string{}, nil
		},
	})
	stdout := CheckSuccess(c, runf, "find", "-a", "admin.agent", "-e", "bob@example.com")
	c.Assert(stdout, gc.Equals, "")
}

func (s *findSuite) TestFindNoParameters(c *gc.C) {

	runf := s.RunServer(c, &handler{
		queryUsers: func(req *params.QueryUsersRequest) ([]string, error) {
			return []string{"alice", "bob", "charlie"}, nil
		},
	})
	stdout := CheckSuccess(c, runf, "find", "-a", "admin.agent", "--format", "json")
	var usernames []string
	err := json.Unmarshal([]byte(stdout), &usernames)
	c.Assert(err, gc.Equals, nil)
	c.Assert(usernames, jc.DeepEquals, []string{"alice", "bob", "charlie"})
}

func (s *findSuite) TestFindLastLoginTime(c *gc.C) {
	var gotTime time.Time
	runf := s.RunServer(c, &handler{
		queryUsers: func(req *params.QueryUsersRequest) ([]string, error) {
			if err := gotTime.UnmarshalText([]byte(req.LastLoginSince)); err != nil {
				return nil, err
			}
			return []string{"alice", "bob", "charlie"}, nil
		},
	})
	stdout := CheckSuccess(c, runf, "find", "-a", "admin.agent", "--format", "json", "--last-login", "30")
	var usernames []string
	err := json.Unmarshal([]byte(stdout), &usernames)
	c.Assert(err, gc.Equals, nil)
	c.Assert(usernames, jc.DeepEquals, []string{"alice", "bob", "charlie"})
	t := time.Now().AddDate(0, 0, -30)
	c.Assert(t.Sub(gotTime), jc.LessThan, time.Second)
}

func (s *findSuite) TestFindLastDischargeTime(c *gc.C) {
	var gotTime time.Time
	runf := s.RunServer(c, &handler{
		queryUsers: func(req *params.QueryUsersRequest) ([]string, error) {
			if err := gotTime.UnmarshalText([]byte(req.LastDischargeSince)); err != nil {
				return nil, err
			}
			return []string{"alice", "bob", "charlie"}, nil
		},
	})
	stdout := CheckSuccess(c, runf, "find", "-a", "admin.agent", "--format", "json", "--last-discharge", "20")
	var usernames []string
	err := json.Unmarshal([]byte(stdout), &usernames)
	c.Assert(err, gc.Equals, nil)
	c.Assert(usernames, jc.DeepEquals, []string{"alice", "bob", "charlie"})
	t := time.Now().AddDate(0, 0, -20)
	c.Assert(t.Sub(gotTime), jc.LessThan, time.Second)
}