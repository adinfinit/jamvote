package event

import (
	"context"
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
		team.EventID = eventid
		team.ID = TeamID(keys[i].IntID())
	}
	return teams, err
}

func allBallots(ctx context.Context, eventkey *datastore.Key) ([]*Ballot, error) {
	var ballots []*Ballot
	q := datastore.NewQuery("Team").Ancestor(eventkey)
	keys, err := q.GetAll(ctx, &ballots)
	for i, ballot := range ballots {
		ballot.ID = keys[i]
	}
	return ballots, err
}

func createTeamInfos(teams []*TeamInfo, ballots []*Ballot) []*TeamInfo {
	infomap := map[TeamID]*TeamInfo{}
	for _, team := range teams {
		info := &TeamInfo{}
		info.Team = team
		infomap[team.ID] = info
	}

	for _, ballot := range ballots {
		info := infomap[ballot.Team]
		info.Ballots = append(info.Ballots, ballot)
		if ballot.Completed {
			info.Complete++
		}
		info.Pending++
	}

	infos := make([]*TeamInfo, 0, len(infomap))
	for _, info := range infomap {
		infos = append(infos, info)
	}
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

func (repo *Datastore) GetBallots(userid user.UserID, eventid EventID) (complete, incomplete []*BallotInfo, err error) {
	err := datastore.RunInTransaction(repo.Context, func(ctx netcontext.Context) error {
		eventkey := datastore.NewKey(ctx, "Event", string(event.ID), 0, nil)

		ballots, err := allBallots(ctx, eventkey)
		if err != nil {
			return nil, err
		}

		teams, err := allTeams(ctx, eventkey)
		if err != nil {
			return nil, err
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

		infos := createTeamInfos(teams, ballots)
		sort.Slice(infos, func(i, k int) bool {
			if infos[i].Pending == infos[k].Pending {
				return infos[i].Complete < infos[k].Complete
			}
			return infos[i].Pending < infos[k].Pending
		})

		var needIncomplete int
		if len(complete) >= 3 {
			needIncomplete = 1
		} else {
			needIncomplete = 3 - len(complete)
		}

		newBallots := []*Ballot{}
		for _, info := range infos {
			if len(incomplete) >= needIncomplete {
				break
			}
			if info.HasReviewer(userid) {
				continue
			}

			// TODO add empty ballot
		}

		// put empty ballots to datastore

		return err
	}, nil)
	return complete, incomplete, datastoreError(err)
}

func (repo *Datastore) SubmitBallot(eventid EventID, ballot *Ballot) error {
	return nil
}

func (repo *Datastore) UserBallots(userid user.UserID, eventid EventID) ([]*Ballot, error) {
	return nil, nil
}

func (repo *Datastore) UserBallot(userid user.UserID, eventid EventID, teamid TeamID) (*Ballot, error) {
	return nil, nil
}

func (repo *Datastore) AllBallots(eventid EventID) ([]*Ballot, error) {
	return nil, nil
}
