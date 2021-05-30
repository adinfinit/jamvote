package datastoredb

import (
	"context"
	"fmt"
	"sort"

	netcontext "golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"

	"github.com/adinfinit/jamvote/event"
	"github.com/adinfinit/jamvote/user"
)

// Events contains methods for managing events.
type Events struct {
	Context context.Context
}

// eventsError converts a datastore error to domain error.
func eventsError(err error) error {
	if err == datastore.ErrNoSuchEntity {
		return event.ErrNotExists
	}
	return err
}

// newEventKey returns event key associated with eventid.
func newEventKey(ctx context.Context, eventid event.EventID) *datastore.Key {
	return datastore.NewKey(ctx, "Event", string(eventid), 0, nil)
}

// newTeamKey returns event key associated with event and teamid.
func newTeamKey(ctx context.Context, eventkey *datastore.Key, teamid event.TeamID) *datastore.Key {
	return datastore.NewKey(ctx, "Team", "", int64(teamid), eventkey)
}

// newBallotKey returns event key associated with event, voter and team.
func newBallotKey(ctx context.Context, eventkey *datastore.Key, voter user.UserID, votingFor event.TeamID) *datastore.Key {
	id := fmt.Sprintf("%v-%v", voter, votingFor)
	return datastore.NewKey(ctx, "Ballot", id, 0, eventkey)
}

// List returns all events.
func (repo *Events) List() ([]*event.Event, error) {
	var events []*event.Event

	q := datastore.NewQuery("Event")
	keys, err := q.GetAll(repo.Context, &events)
	for i, ev := range events {
		ev.ID = event.EventID(keys[i].StringID())
	}

	return events, err
}

// Create creates a new event.
func (repo *Events) Create(ev *event.Event) error {
	err := datastore.RunInTransaction(repo.Context, func(ctx netcontext.Context) error {
		eventkey := datastore.NewKey(ctx, "Event", string(ev.ID), 0, nil)

		existing := &event.Event{}
		err := datastore.Get(ctx, eventkey, existing)
		if err != datastore.ErrNoSuchEntity {
			if err == nil {
				return event.ErrExists
			}
			return err
		}

		_, err = datastore.Put(ctx, eventkey, ev)
		return err
	}, nil)

	return eventsError(err)
}

// ByID retrieves an event by ID.
func (repo *Events) ByID(eventid event.EventID) (*event.Event, error) {
	ev := &event.Event{}
	if _, err := memcache.Gob.Get(repo.Context, "Event_"+eventid.String(), ev); err == nil {
		return ev, nil
	}

	ev = &event.Event{}
	ev.ID = eventid
	eventkey := datastore.NewKey(repo.Context, "Event", string(eventid), 0, nil)
	err := datastore.Get(repo.Context, eventkey, ev)

	memcache.Gob.Add(repo.Context, &memcache.Item{
		Key:    "Event_" + eventid.String(),
		Object: ev,
	})

	return ev, eventsError(err)
}

// Update updates an existing event.
func (repo *Events) Update(ev *event.Event) error {
	eventkey := datastore.NewKey(repo.Context, "Event", string(ev.ID), 0, nil)
	_, err := datastore.Put(repo.Context, eventkey, ev)
	if err == nil {
		memcache.Gob.Set(repo.Context, &memcache.Item{
			Key:    "Event_" + ev.ID.String(),
			Object: ev,
		})
	}
	return eventsError(err)
}

// CreateTeam creates a new team.
func (repo *Events) CreateTeam(eventid event.EventID, team *event.Team) (event.TeamID, error) {
	eventkey := datastore.NewKey(repo.Context, "Event", string(eventid), 0, nil)
	incompletekey := datastore.NewIncompleteKey(repo.Context, "Team", eventkey)
	teamkey, err := datastore.Put(repo.Context, incompletekey, team)
	team.ID = event.TeamID(teamkey.IntID())
	return team.ID, eventsError(err)
}

// UpdateTeam updates a team.
func (repo *Events) UpdateTeam(eventid event.EventID, team *event.Team) error {
	eventkey := datastore.NewKey(repo.Context, "Event", string(eventid), 0, nil)
	teamkey := datastore.NewKey(repo.Context, "Team", "", int64(team.ID), eventkey)
	_, err := datastore.Put(repo.Context, teamkey, team)
	return eventsError(err)
}

// TeamByID retrieves a team by ID.
func (repo *Events) TeamByID(eventid event.EventID, teamid event.TeamID) (*event.Team, error) {
	team := &event.Team{}
	eventkey := datastore.NewKey(repo.Context, "Event", string(eventid), 0, nil)
	teamkey := datastore.NewKey(repo.Context, "Team", "", int64(teamid), eventkey)
	err := datastore.Get(repo.Context, teamkey, team)
	team.EventID = eventid
	team.ID = teamid
	return team, eventsError(err)
}

// DeleteTeam deletes a team.
func (repo *Events) DeleteTeam(eventid event.EventID, teamid event.TeamID) error {
	eventkey := datastore.NewKey(repo.Context, "Event", string(eventid), 0, nil)
	teamkey := datastore.NewKey(repo.Context, "Team", "", int64(teamid), eventkey)

	err := datastore.Delete(repo.Context, teamkey)
	return eventsError(err)
}

// TeamsByUser returns all teams associated with the user.
func (repo *Events) TeamsByUser(userid user.UserID) ([]*event.EventTeam, error) {
	var allTeams []*event.Team
	q := datastore.NewQuery("Team")
	keys, err := q.GetAll(repo.Context, &allTeams)
	if err != nil {
		return nil, err
	}

	events, err := repo.List()
	if err != nil {
		return nil, err
	}

	teams := []*event.EventTeam{}
	for i, team := range allTeams {
		if !team.HasMemberID(userid) {
			continue
		}

		team.EventID = event.EventID(keys[i].Parent().StringID())
		team.ID = event.TeamID(keys[i].IntID())

		var ev *event.Event
		for _, e := range events {
			if e.ID == team.EventID {
				ev = e
				break
			}
		}

		teams = append(teams, &event.EventTeam{
			Event: *ev,
			Team:  *team,
		})
	}

	return teams, err
}

// allTeams retrieves all teams in an event.
func allTeams(ctx context.Context, eventkey *datastore.Key) ([]*event.Team, error) {
	var teams []*event.Team
	q := datastore.NewQuery("Team").Ancestor(eventkey)
	keys, err := q.GetAll(ctx, &teams)
	for i, team := range teams {
		team.EventID = event.EventID(eventkey.StringID())
		team.ID = event.TeamID(keys[i].IntID())
	}
	return teams, err
}

// someTeams retrieves all teams specified in the teamids.
func someTeams(ctx context.Context, eventkey *datastore.Key, teamids []event.TeamID) ([]*event.Team, error) {
	var keys []*datastore.Key
	var teams []*event.Team
	for _, teamid := range teamids {
		keys = append(keys, datastore.NewKey(ctx, "Team", "", int64(teamid), eventkey))
		teams = append(teams, &event.Team{})
	}

	err := datastore.GetMulti(ctx, keys, teams)
	for i, team := range teams {
		team.EventID = event.EventID(eventkey.StringID())
		team.ID = event.TeamID(keys[i].IntID())
	}

	return teams, err
}

// allBallots returns all event ballots.
func allBallots(ctx context.Context, eventkey *datastore.Key) ([]*event.Ballot, error) {
	var ballots []*event.Ballot
	q := datastore.NewQuery("Ballot").Ancestor(eventkey)
	keys, err := q.GetAll(ctx, &ballots)
	for i, ballot := range ballots {
		ballot.ID = keys[i]
	}
	return ballots, err
}

// teamBallots returns all event ballots in a team.
func teamBallots(ctx context.Context, eventkey *datastore.Key, teamid event.TeamID) ([]*event.Ballot, error) {
	var ballots []*event.Ballot
	q := datastore.NewQuery("Ballot").Ancestor(eventkey).Filter("Team =", teamid)
	keys, err := q.GetAll(ctx, &ballots)
	for i, ballot := range ballots {
		ballot.ID = keys[i]
	}
	return ballots, err
}

// userBallots returns all event ballots of an user.
func userBallots(ctx context.Context, eventkey *datastore.Key, userid user.UserID) ([]*event.Ballot, error) {
	var ballots []*event.Ballot
	q := datastore.NewQuery("Ballot").Ancestor(eventkey).Filter("Voter =", userid)
	keys, err := q.GetAll(ctx, &ballots)
	for i, ballot := range ballots {
		ballot.ID = keys[i]
	}
	return ballots, err
}

// userBallot returns a specific event ballot.
func userBallot(ctx context.Context, eventkey *datastore.Key, userid user.UserID, teamid event.TeamID) (*event.Ballot, error) {
	ballot := &event.Ballot{}
	ballot.ID = newBallotKey(ctx, eventkey, userid, teamid)
	err := datastore.Get(ctx, ballot.ID, ballot)
	return ballot, err
}

// createTeamResults summarizes teams and ballots into TeamResult.
func createTeamResults(teams []*event.Team, ballots []*event.Ballot) []*event.TeamResult {
	cross := map[event.TeamID]*event.TeamResult{}
	for _, team := range teams {
		res := &event.TeamResult{}
		res.Team = team
		cross[team.ID] = res
	}

	for _, ballot := range ballots {
		res := cross[ballot.Team]
		if res.Team.HasMemberID(ballot.Voter) {
			res.MemberBallots = append(res.MemberBallots, ballot)
			continue
		}

		res.Ballots = append(res.Ballots, ballot)
		if ballot.Completed {
			res.Complete++
		}
		res.Pending++
	}

	results := make([]*event.TeamResult, 0, len(cross))
	for _, res := range cross {
		results = append(results, res)
	}
	return results
}

// createBallotInfos associates ballots with teams.
func createBallotInfos(teams []*event.Team, ballots []*event.Ballot) []*event.BallotInfo {
	infos := make([]*event.BallotInfo, 0, len(ballots))
	for _, ballot := range ballots {
		infos = append(infos, &event.BallotInfo{
			Ballot: ballot,
			Team:   findTeam(teams, ballot.Team),
		})
	}

	sort.Slice(infos, func(i, k int) bool {
		return infos[i].Ballot.Index < infos[k].Ballot.Index
	})

	return infos
}

// findTeam finds team with the specified id from the slice.
func findTeam(teams []*event.Team, id event.TeamID) *event.Team {
	for _, team := range teams {
		if team.ID == id {
			return team
		}
	}
	return nil
}

// Teams returns all teams in an event.
func (repo *Events) Teams(eventid event.EventID) ([]*event.Team, error) {
	eventkey := datastore.NewKey(repo.Context, "Event", string(eventid), 0, nil)
	teams, err := allTeams(repo.Context, eventkey)
	return teams, eventsError(err)
}

// CreateIncompleteBallots creates new incomplete ballots for a user.
func (repo *Events) CreateIncompleteBallots(eventid event.EventID, userid user.UserID) (complete, incomplete []*event.BallotInfo, err error) {
	const FirstBatchCount = 3

	//TODO: extract tranaction from here
	err = datastore.RunInTransaction(repo.Context, func(ctx netcontext.Context) error {
		eventkey := datastore.NewKey(ctx, "Event", string(eventid), 0, nil)

		ballots, err := allBallots(ctx, eventkey)
		if err != nil {
			return err
		}

		teams, err := allTeams(ctx, eventkey)
		if err != nil {
			return err
		}

		for _, ballot := range ballots {
			if ballot.Voter == userid {
				info := &event.BallotInfo{
					Team:   findTeam(teams, ballot.Team),
					Ballot: ballot,
				}

				if ballot.Completed {
					complete = append(complete, info)
				} else {
					incomplete = append(incomplete, info)
				}
			}
		}

		// user has not completed first batch?
		hasFullBatch := len(complete)+len(incomplete) >= FirstBatchCount
		if hasFullBatch && len(incomplete) > 0 {
			return nil
		}

		teamresults := createTeamResults(teams, ballots)
		sort.Slice(teamresults, func(i, k int) bool {
			if teamresults[i].Pending == teamresults[k].Pending {
				return teamresults[i].Complete < teamresults[k].Complete
			}
			return teamresults[i].Pending < teamresults[k].Pending
		})

		// TODO: don't hardcode and move this logic to service level
		var needIncomplete int
		if len(complete) >= FirstBatchCount {
			needIncomplete = 1
		} else {
			needIncomplete = FirstBatchCount
		}

		createKeys := []*datastore.Key{}
		createBallots := []*event.Ballot{}
		for _, teamresult := range teamresults {
			if len(incomplete) >= needIncomplete {
				break
			}
			if teamresult.HasReviewer(userid) {
				continue
			}
			if teamresult.HasMemberID(userid) {
				continue
			}
			if !teamresult.HasSubmitted() {
				continue
			}

			ballot := &event.Ballot{
				ID:        newBallotKey(ctx, eventkey, userid, teamresult.Team.ID),
				Voter:     userid,
				Team:      teamresult.Team.ID,
				Index:     int64(len(complete) + len(incomplete)),
				Completed: false,
				Aspects:   event.DefaultAspects,
			}

			ballotinfo := &event.BallotInfo{
				Team:   teamresult.Team,
				Ballot: ballot,
			}

			createKeys = append(createKeys, ballot.ID)
			createBallots = append(createBallots, ballot)
			incomplete = append(incomplete, ballotinfo)
		}

		if len(createKeys) > 0 {
			_, err = datastore.PutMulti(ctx, createKeys, createBallots)
		}

		return err
	}, nil)

	return complete, incomplete, eventsError(err)
}

// SubmitBallot submits a ballot.
func (repo *Events) SubmitBallot(eventid event.EventID, ballot *event.Ballot) error {
	eventkey := datastore.NewKey(repo.Context, "Event", string(eventid), 0, nil)
	ballot.ID = newBallotKey(repo.Context, eventkey, ballot.Voter, ballot.Team)
	_, err := datastore.Put(repo.Context, ballot.ID, ballot)
	return eventsError(err)
}

// UserBallot retrieves a user ballot.
func (repo *Events) UserBallot(eventid event.EventID, userid user.UserID, teamid event.TeamID) (*event.Ballot, error) {
	eventkey := datastore.NewKey(repo.Context, "Event", string(eventid), 0, nil)
	ballot, err := userBallot(repo.Context, eventkey, userid, teamid)
	return ballot, eventsError(err)
}

// UserBallots retrieves all user ballots.
func (repo *Events) UserBallots(eventid event.EventID, userid user.UserID) ([]*event.BallotInfo, error) {
	eventkey := datastore.NewKey(repo.Context, "Event", string(eventid), 0, nil)
	ballots, err := userBallots(repo.Context, eventkey, userid)
	if err != nil {
		return nil, eventsError(err)
	}

	teamids := []event.TeamID{}
	for _, ballot := range ballots {
		teamids = append(teamids, ballot.Team)
	}

	teams, err := someTeams(repo.Context, eventkey, teamids)
	if err != nil {
		return nil, eventsError(err)
	}

	return createBallotInfos(teams, ballots), nil
}

// Results retrieves results for an event.
func (repo *Events) Results(eventid event.EventID) ([]*event.TeamResult, error) {
	eventkey := datastore.NewKey(repo.Context, "Event", string(eventid), 0, nil)

	ballots, err := allBallots(repo.Context, eventkey)
	if err != nil {
		return nil, eventsError(err)
	}

	teams, err := allTeams(repo.Context, eventkey)
	if err != nil {
		return nil, eventsError(err)
	}

	results := createTeamResults(teams, ballots)
	for _, result := range results {
		result.Average = event.AverageScores(result.Ballots)
	}

	return results, nil
}

// Ballots retrieves all event ballots.
func (repo *Events) Ballots(eventid event.EventID) ([]*event.Ballot, error) {
	eventkey := datastore.NewKey(repo.Context, "Event", string(eventid), 0, nil)
	ballots, err := allBallots(repo.Context, eventkey)
	return ballots, eventsError(err)
}

// TeamBallots retrieves all event ballots for a team.
func (repo *Events) TeamBallots(eventid event.EventID, teamid event.TeamID) ([]*event.Ballot, error) {
	eventkey := datastore.NewKey(repo.Context, "Event", string(eventid), 0, nil)
	ballots, err := teamBallots(repo.Context, eventkey, teamid)
	return ballots, eventsError(err)
}
