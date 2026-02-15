package datastoredb

import (
	"context"
	"fmt"
	"sort"

	"cloud.google.com/go/datastore"

	"github.com/adinfinit/jamvote/event"
	"github.com/adinfinit/jamvote/user"
)

// Events contains methods for managing events.
type Events struct {
	Context context.Context
	Client  *datastore.Client
}

// eventsError converts a datastore error to domain error.
func eventsError(err error) error {
	if err == datastore.ErrNoSuchEntity {
		return event.ErrNotExists
	}
	return err
}

// newEventKey returns event key associated with eventid.
func newEventKey(eventid event.EventID) *datastore.Key {
	return datastore.NameKey("Event", string(eventid), nil)
}

// newTeamKey returns event key associated with event and teamid.
func newTeamKey(eventkey *datastore.Key, teamid event.TeamID) *datastore.Key {
	return datastore.IDKey("Team", int64(teamid), eventkey)
}

// newBallotKey returns event key associated with event, voter and team.
func newBallotKey(eventkey *datastore.Key, voter user.UserID, votingFor event.TeamID) *datastore.Key {
	id := fmt.Sprintf("%v-%v", voter, votingFor)
	return datastore.NameKey("Ballot", id, eventkey)
}

// List returns all events.
func (repo *Events) List() ([]*event.Event, error) {
	var events []*event.Event

	q := datastore.NewQuery("Event")
	keys, err := repo.Client.GetAll(repo.Context, q, &events)
	for i, ev := range events {
		ev.ID = event.EventID(keys[i].Name)
	}

	return events, err
}

// Create creates a new event.
func (repo *Events) Create(ev *event.Event) error {
	_, err := repo.Client.RunInTransaction(repo.Context, func(tx *datastore.Transaction) error {
		eventkey := newEventKey(ev.ID)

		existing := &event.Event{}
		err := tx.Get(eventkey, existing)
		if err != datastore.ErrNoSuchEntity {
			if err == nil {
				return event.ErrExists
			}
			return err
		}

		_, err = tx.Put(eventkey, ev)
		return err
	})

	return eventsError(err)
}

// ByID retrieves an event by ID.
func (repo *Events) ByID(eventid event.EventID) (*event.Event, error) {
	ev := &event.Event{}
	if appCache.Get("Event_"+eventid.String(), ev) {
		return ev, nil
	}

	ev = &event.Event{}
	ev.ID = eventid
	eventkey := newEventKey(eventid)
	err := repo.Client.Get(repo.Context, eventkey, ev)
	if err == nil {
		appCache.Set("Event_"+eventid.String(), ev)
	}

	return ev, eventsError(err)
}

// Update updates an existing event.
func (repo *Events) Update(ev *event.Event) error {
	eventkey := newEventKey(ev.ID)
	_, err := repo.Client.Put(repo.Context, eventkey, ev)
	if err == nil {
		appCache.Set("Event_"+ev.ID.String(), ev)
	}
	return eventsError(err)
}

// CreateTeam creates a new team.
func (repo *Events) CreateTeam(eventid event.EventID, team *event.Team) (event.TeamID, error) {
	eventkey := newEventKey(eventid)
	incompletekey := datastore.IncompleteKey("Team", eventkey)
	teamkey, err := repo.Client.Put(repo.Context, incompletekey, team)
	team.ID = event.TeamID(teamkey.ID)
	return team.ID, eventsError(err)
}

// UpdateTeam updates a team.
func (repo *Events) UpdateTeam(eventid event.EventID, team *event.Team) error {
	eventkey := newEventKey(eventid)
	teamkey := datastore.IDKey("Team", int64(team.ID), eventkey)
	_, err := repo.Client.Put(repo.Context, teamkey, team)
	return eventsError(err)
}

// TeamByID retrieves a team by ID.
func (repo *Events) TeamByID(eventid event.EventID, teamid event.TeamID) (*event.Team, error) {
	team := &event.Team{}
	eventkey := newEventKey(eventid)
	teamkey := datastore.IDKey("Team", int64(teamid), eventkey)
	err := repo.Client.Get(repo.Context, teamkey, team)
	team.EventID = eventid
	team.ID = teamid
	return team, eventsError(err)
}

// DeleteTeam deletes a team.
func (repo *Events) DeleteTeam(eventid event.EventID, teamid event.TeamID) error {
	eventkey := newEventKey(eventid)
	teamkey := datastore.IDKey("Team", int64(teamid), eventkey)

	err := repo.Client.Delete(repo.Context, teamkey)
	return eventsError(err)
}

// TeamsByUser returns all teams associated with the user.
func (repo *Events) TeamsByUser(userid user.UserID) ([]*event.EventTeam, error) {
	var allTeamsList []*event.Team
	q := datastore.NewQuery("Team")
	keys, err := repo.Client.GetAll(repo.Context, q, &allTeamsList)
	if err != nil {
		return nil, err
	}

	events, err := repo.List()
	if err != nil {
		return nil, err
	}

	teams := []*event.EventTeam{}
	for i, team := range allTeamsList {
		if !team.HasMemberID(userid) {
			continue
		}

		team.EventID = event.EventID(keys[i].Parent.Name)
		team.ID = event.TeamID(keys[i].ID)

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
// If tx is non-nil, the query runs inside that transaction.
func (repo *Events) allTeams(eventkey *datastore.Key, tx *datastore.Transaction) ([]*event.Team, error) {
	var teams []*event.Team
	q := datastore.NewQuery("Team").Ancestor(eventkey)
	if tx != nil {
		q = q.Transaction(tx)
	}
	keys, err := repo.Client.GetAll(repo.Context, q, &teams)
	for i, team := range teams {
		team.EventID = event.EventID(eventkey.Name)
		team.ID = event.TeamID(keys[i].ID)
	}
	return teams, err
}

// someTeams retrieves all teams specified in the teamids.
func (repo *Events) someTeams(eventkey *datastore.Key, teamids []event.TeamID) ([]*event.Team, error) {
	var keys []*datastore.Key
	var teams []*event.Team
	for _, teamid := range teamids {
		keys = append(keys, newTeamKey(eventkey, teamid))
		teams = append(teams, &event.Team{})
	}

	err := repo.Client.GetMulti(repo.Context, keys, teams)
	for i, team := range teams {
		team.EventID = event.EventID(eventkey.Name)
		team.ID = event.TeamID(keys[i].ID)
	}

	return teams, err
}

// allBallots returns all event ballots.
// If tx is non-nil, the query runs inside that transaction.
func (repo *Events) allBallots(eventkey *datastore.Key, tx *datastore.Transaction) ([]*event.Ballot, error) {
	var ballots []*event.Ballot
	q := datastore.NewQuery("Ballot").Ancestor(eventkey)
	if tx != nil {
		q = q.Transaction(tx)
	}
	keys, err := repo.Client.GetAll(repo.Context, q, &ballots)
	for i, ballot := range ballots {
		ballot.ID = keys[i]
	}
	return ballots, err
}

// teamBallots returns all event ballots in a team.
func (repo *Events) teamBallots(eventkey *datastore.Key, teamid event.TeamID) ([]*event.Ballot, error) {
	var ballots []*event.Ballot
	q := datastore.NewQuery("Ballot").Ancestor(eventkey).FilterField("Team", "=", teamid)
	keys, err := repo.Client.GetAll(repo.Context, q, &ballots)
	for i, ballot := range ballots {
		ballot.ID = keys[i]
	}
	return ballots, err
}

// userBallots returns all event ballots of an user.
func (repo *Events) userBallots(eventkey *datastore.Key, userid user.UserID) ([]*event.Ballot, error) {
	var ballots []*event.Ballot
	q := datastore.NewQuery("Ballot").Ancestor(eventkey).FilterField("Voter", "=", userid)
	keys, err := repo.Client.GetAll(repo.Context, q, &ballots)
	for i, ballot := range ballots {
		ballot.ID = keys[i]
	}
	return ballots, err
}

// userBallot returns a specific event ballot.
func (repo *Events) userBallot(eventkey *datastore.Key, userid user.UserID, teamid event.TeamID) (*event.Ballot, error) {
	ballot := &event.Ballot{}
	ballot.ID = newBallotKey(eventkey, userid, teamid)
	err := repo.Client.Get(repo.Context, ballot.ID, ballot)
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
	eventkey := newEventKey(eventid)
	teams, err := repo.allTeams(eventkey, nil)
	return teams, eventsError(err)
}

// CreateIncompleteBallots creates new incomplete ballots for a user.
func (repo *Events) CreateIncompleteBallots(eventid event.EventID, userid user.UserID) (complete, incomplete []*event.BallotInfo, err error) {
	const FirstBatchCount = 3

	//TODO: extract transaction from here
	_, txErr := repo.Client.RunInTransaction(repo.Context, func(tx *datastore.Transaction) error {
		eventkey := newEventKey(eventid)

		ballots, err := repo.allBallots(eventkey, tx)
		if err != nil {
			return err
		}

		teams, err := repo.allTeams(eventkey, tx)
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
				ID:        newBallotKey(eventkey, userid, teamresult.Team.ID),
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
			_, err = tx.PutMulti(createKeys, createBallots)
		}

		return err
	})

	return complete, incomplete, eventsError(txErr)
}

// SubmitBallot submits a ballot.
func (repo *Events) SubmitBallot(eventid event.EventID, ballot *event.Ballot) error {
	eventkey := newEventKey(eventid)
	ballot.ID = newBallotKey(eventkey, ballot.Voter, ballot.Team)
	_, err := repo.Client.Put(repo.Context, ballot.ID, ballot)
	return eventsError(err)
}

// UserBallot retrieves a user ballot.
func (repo *Events) UserBallot(eventid event.EventID, userid user.UserID, teamid event.TeamID) (*event.Ballot, error) {
	eventkey := newEventKey(eventid)
	ballot, err := repo.userBallot(eventkey, userid, teamid)
	return ballot, eventsError(err)
}

// UserBallots retrieves all user ballots.
func (repo *Events) UserBallots(eventid event.EventID, userid user.UserID) ([]*event.BallotInfo, error) {
	eventkey := newEventKey(eventid)
	ballots, err := repo.userBallots(eventkey, userid)
	if err != nil {
		return nil, eventsError(err)
	}

	teamids := []event.TeamID{}
	for _, ballot := range ballots {
		teamids = append(teamids, ballot.Team)
	}

	teams, err := repo.someTeams(eventkey, teamids)
	if err != nil {
		return nil, eventsError(err)
	}

	return createBallotInfos(teams, ballots), nil
}

// Results retrieves results for an event.
func (repo *Events) Results(eventid event.EventID) ([]*event.TeamResult, error) {
	eventkey := newEventKey(eventid)

	ballots, err := repo.allBallots(eventkey, nil)
	if err != nil {
		return nil, eventsError(err)
	}

	teams, err := repo.allTeams(eventkey, nil)
	if err != nil {
		return nil, eventsError(err)
	}

	results := createTeamResults(teams, ballots)

	currentEvent, err := repo.ByID(eventid)

	if err != nil {
		return nil, eventsError(err)
	}

	for _, result := range results {
		result.Average, result.JammerAverage, result.JudgeAverage = event.AverageScores(result.Ballots, currentEvent)
	}

	return results, nil
}

// Ballots retrieves all event ballots.
func (repo *Events) Ballots(eventid event.EventID) ([]*event.Ballot, error) {
	eventkey := newEventKey(eventid)
	ballots, err := repo.allBallots(eventkey, nil)
	return ballots, eventsError(err)
}

// TeamBallots retrieves all event ballots for a team.
func (repo *Events) TeamBallots(eventid event.EventID, teamid event.TeamID) ([]*event.Ballot, error) {
	eventkey := newEventKey(eventid)
	ballots, err := repo.teamBallots(eventkey, teamid)
	return ballots, eventsError(err)
}
