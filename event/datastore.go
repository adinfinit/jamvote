package event

import (
	"context"
	"fmt"
	"sort"

	"github.com/adinfinit/jamvote/user"
	netcontext "golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type Datastore struct {
	Context context.Context
}

func datastoreError(err error) error {
	if err == datastore.ErrNoSuchEntity {
		return ErrNotExists
	}
	return err
}

func newEventKey(ctx context.Context, eventid EventID) *datastore.Key {
	return datastore.NewKey(ctx, "Event", string(eventid), 0, nil)
}
func newTeamKey(ctx context.Context, eventkey *datastore.Key, teamid TeamID) *datastore.Key {
	return datastore.NewKey(ctx, "Team", "", int64(teamid), eventkey)
}
func newBallotKey(ctx context.Context, eventkey *datastore.Key, voter user.UserID, votingFor TeamID) *datastore.Key {
	id := fmt.Sprintf("%v-%v", voter, votingFor)
	return datastore.NewKey(ctx, "Ballot", id, 0, eventkey)
}

func (repo *Datastore) List() ([]*Event, error) {
	var events []*Event

	q := datastore.NewQuery("Event")
	keys, err := q.GetAll(repo.Context, &events)
	for i, event := range events {
		event.ID = EventID(keys[i].StringID())
	}

	return events, err
}

func (repo *Datastore) Create(event *Event) error {
	err := datastore.RunInTransaction(repo.Context, func(ctx netcontext.Context) error {
		eventkey := datastore.NewKey(ctx, "Event", string(event.ID), 0, nil)

		existing := &Event{}
		err := datastore.Get(ctx, eventkey, existing)
		if err != datastore.ErrNoSuchEntity {
			if err == nil {
				return ErrExists
			}
			return err
		}

		_, err = datastore.Put(ctx, eventkey, event)
		return err
	}, nil)

	return datastoreError(err)
}

func (repo *Datastore) ByID(eventid EventID) (*Event, error) {
	event := &Event{}
	event.ID = eventid
	eventkey := datastore.NewKey(repo.Context, "Event", string(eventid), 0, nil)
	err := datastore.Get(repo.Context, eventkey, event)
	return event, datastoreError(err)
}

func (repo *Datastore) Update(event *Event) error {
	eventkey := datastore.NewKey(repo.Context, "Event", string(event.ID), 0, nil)
	_, err := datastore.Put(repo.Context, eventkey, event)
	return datastoreError(err)
}

func (repo *Datastore) CreateTeam(eventid EventID, team *Team) (TeamID, error) {
	eventkey := datastore.NewKey(repo.Context, "Event", string(eventid), 0, nil)
	incompletekey := datastore.NewIncompleteKey(repo.Context, "Team", eventkey)
	teamkey, err := datastore.Put(repo.Context, incompletekey, team)
	team.ID = TeamID(teamkey.IntID())
	return team.ID, datastoreError(err)
}

func (repo *Datastore) UpdateTeam(eventid EventID, team *Team) error {
	eventkey := datastore.NewKey(repo.Context, "Event", string(eventid), 0, nil)
	teamkey := datastore.NewKey(repo.Context, "Team", "", int64(team.ID), eventkey)
	_, err := datastore.Put(repo.Context, teamkey, team)
	return datastoreError(err)
}

func (repo *Datastore) TeamByID(eventid EventID, teamid TeamID) (*Team, error) {
	team := &Team{}
	eventkey := datastore.NewKey(repo.Context, "Event", string(eventid), 0, nil)
	teamkey := datastore.NewKey(repo.Context, "Team", "", int64(teamid), eventkey)
	err := datastore.Get(repo.Context, teamkey, team)
	team.EventID = eventid
	team.ID = teamid
	return team, datastoreError(err)
}

func allTeams(ctx context.Context, eventkey *datastore.Key) ([]*Team, error) {
	var teams []*Team
	q := datastore.NewQuery("Team").Ancestor(eventkey)
	keys, err := q.GetAll(ctx, &teams)
	for i, team := range teams {
		team.EventID = EventID(eventkey.StringID())
		team.ID = TeamID(keys[i].IntID())
	}
	return teams, err
}

func someTeams(ctx context.Context, eventkey *datastore.Key, teamids []TeamID) ([]*Team, error) {
	var keys []*datastore.Key
	var teams []*Team
	for _, teamid := range teamids {
		keys = append(keys, datastore.NewKey(ctx, "Team", "", int64(teamid), eventkey))
		teams = append(teams, &Team{})
	}

	err := datastore.GetMulti(ctx, keys, teams)
	for i, team := range teams {
		team.EventID = EventID(eventkey.StringID())
		team.ID = TeamID(keys[i].IntID())
	}

	return teams, err
}

func allBallots(ctx context.Context, eventkey *datastore.Key) ([]*Ballot, error) {
	var ballots []*Ballot
	q := datastore.NewQuery("Ballot").Ancestor(eventkey)
	keys, err := q.GetAll(ctx, &ballots)
	for i, ballot := range ballots {
		ballot.ID = keys[i]
	}
	return ballots, err
}

func userBallots(ctx context.Context, eventkey *datastore.Key, userid user.UserID) ([]*Ballot, error) {
	var ballots []*Ballot
	q := datastore.NewQuery("Ballot").Ancestor(eventkey).Filter("Voter =", userid)
	keys, err := q.GetAll(ctx, &ballots)
	for i, ballot := range ballots {
		ballot.ID = keys[i]
	}
	return ballots, err
}

func userBallot(ctx context.Context, eventkey *datastore.Key, userid user.UserID, teamid TeamID) (*Ballot, error) {
	ballot := &Ballot{}
	ballot.ID = newBallotKey(ctx, eventkey, userid, teamid)
	err := datastore.Get(ctx, ballot.ID, ballot)
	return ballot, err
}

func createTeamResults(teams []*Team, ballots []*Ballot) []*TeamResult {
	cross := map[TeamID]*TeamResult{}
	for _, team := range teams {
		res := &TeamResult{}
		res.Team = team
		cross[team.ID] = res
	}

	for _, ballot := range ballots {
		res := cross[ballot.Team]
		res.Ballots = append(res.Ballots, ballot)
		if ballot.Completed {
			res.Complete++
		}
		res.Pending++
	}

	results := make([]*TeamResult, 0, len(cross))
	for _, res := range cross {
		results = append(results, res)
	}
	return results
}

func createBallotInfos(teams []*Team, ballots []*Ballot) []*BallotInfo {
	infos := make([]*BallotInfo, 0, len(ballots))
	for _, ballot := range ballots {
		infos = append(infos, &BallotInfo{
			Ballot: ballot,
			Team:   findTeam(teams, ballot.Team),
		})
	}

	sort.Slice(infos, func(i, k int) bool {
		return infos[i].Ballot.Index < infos[k].Ballot.Index
	})

	return infos
}

func findTeam(teams []*Team, id TeamID) *Team {
	for _, team := range teams {
		if team.ID == id {
			return team
		}
	}
	return nil
}

func (repo *Datastore) Teams(eventid EventID) ([]*Team, error) {
	eventkey := datastore.NewKey(repo.Context, "Event", string(eventid), 0, nil)
	teams, err := allTeams(repo.Context, eventkey)
	return teams, datastoreError(err)
}

func (repo *Datastore) CreateIncompleteBallots(eventid EventID, userid user.UserID) (complete, incomplete []*BallotInfo, err error) {
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
				info := &BallotInfo{
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

		// user has not complete all previous
		if len(incomplete) > 0 {
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
		if len(complete) >= 3 {
			needIncomplete = 1
		} else {
			needIncomplete = 3 - len(complete)
		}

		createKeys := []*datastore.Key{}
		createBallots := []*Ballot{}
		for _, teamresult := range teamresults {
			if len(incomplete) >= needIncomplete {
				break
			}
			if teamresult.HasReviewer(userid) {
				continue
			}
			if !teamresult.HasSubmitted() {
				continue
			}

			ballot := &Ballot{
				ID:        newBallotKey(ctx, eventkey, userid, teamresult.Team.ID),
				Voter:     userid,
				Team:      teamresult.Team.ID,
				Index:     int64(len(complete) + len(incomplete)),
				Completed: false,
				Aspects:   DefaultAspects,
			}

			ballotinfo := &BallotInfo{
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

	return complete, incomplete, datastoreError(err)
}

func (repo *Datastore) SubmitBallot(eventid EventID, ballot *Ballot) error {
	eventkey := datastore.NewKey(repo.Context, "Event", string(eventid), 0, nil)
	ballot.ID = newBallotKey(repo.Context, eventkey, ballot.Voter, ballot.Team)
	_, err := datastore.Put(repo.Context, ballot.ID, ballot)
	return datastoreError(err)
}

func (repo *Datastore) UserBallot(eventid EventID, userid user.UserID, teamid TeamID) (*Ballot, error) {
	eventkey := datastore.NewKey(repo.Context, "Event", string(eventid), 0, nil)
	ballot, err := userBallot(repo.Context, eventkey, userid, teamid)
	return ballot, datastoreError(err)
}

func (repo *Datastore) UserBallots(eventid EventID, userid user.UserID) ([]*BallotInfo, error) {
	eventkey := datastore.NewKey(repo.Context, "Event", string(eventid), 0, nil)
	ballots, err := userBallots(repo.Context, eventkey, userid)
	if err != nil {
		return nil, err
	}

	teamids := []TeamID{}
	for _, ballot := range ballots {
		teamids = append(teamids, ballot.Team)
	}

	teams, err := someTeams(repo.Context, eventkey, teamids)
	if err != nil {
		return nil, err
	}

	return createBallotInfos(teams, ballots), nil
}

func (repo *Datastore) Results(eventid EventID) ([]*TeamResult, error) {
	eventkey := datastore.NewKey(repo.Context, "Event", string(eventid), 0, nil)

	ballots, err := allBallots(repo.Context, eventkey)
	if err != nil {
		return nil, err
	}

	teams, err := allTeams(repo.Context, eventkey)
	if err != nil {
		return nil, err
	}

	return createTeamResults(teams, ballots), nil
}
