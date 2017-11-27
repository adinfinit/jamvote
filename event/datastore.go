package event

import (
	"context"

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

// TODO
func (repo *Datastore) Update(event *Event) error { return nil }

func (repo *Datastore) CreateTeam(eventid EventID, team *Team) (TeamID, error) {
	eventkey := datastore.NewKey(repo.Context, "Event", string(eventid), 0, nil)
	incompletekey := datastore.NewIncompleteKey(repo.Context, "Team", eventkey)
	teamkey, err := datastore.Put(repo.Context, incompletekey, team)
	team.ID = TeamID(teamkey.IntID())
	return team.ID, datastoreError(err)

}

// TODO
func (repo *Datastore) UpdateTeam(team *Team) error { return nil }

func (repo *Datastore) TeamByID(eventid EventID, teamid TeamID) (*Team, error) {
	team := &Team{}
	eventkey := datastore.NewKey(repo.Context, "Event", string(eventid), 0, nil)
	teamkey := datastore.NewKey(repo.Context, "Team", "", int64(teamid), eventkey)
	err := datastore.Get(repo.Context, teamkey, team)
	team.EventID = eventid
	team.ID = teamid
	return team, datastoreError(err)
}

func (repo *Datastore) Teams(eventid EventID) ([]*Team, error) {
	var teams []*Team
	eventkey := datastore.NewKey(repo.Context, "Event", string(eventid), 0, nil)
	q := datastore.NewQuery("Team").Ancestor(eventkey)
	keys, err := q.GetAll(repo.Context, &teams)
	for i, team := range teams {
		team.EventID = eventid
		team.ID = TeamID(keys[i].IntID())
	}
	return teams, datastoreError(err)
}

// TODO
func (repo *Datastore) SubmitBallot(eventid EventID, ballot *Ballot) error { return nil }

// TODO
func (repo *Datastore) Ballots(eventid EventID) ([]*Ballot, error) { return nil, nil }

// TODO
func (repo *Datastore) LeastBallots(eventid EventID, n int) ([]*Team, error) { return nil, nil }
