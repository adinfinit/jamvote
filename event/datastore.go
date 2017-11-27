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

func (repo *Datastore) ByID(id EventID) (*Event, error) { return nil, nil }
func (repo *Datastore) Update(event *Event) error       { return nil }

func (repo *Datastore) CreateTeam(id EventID, team *Team) (TeamID, error) { return 0, nil }
func (repo *Datastore) TeamByID(id EventID, teamid TeamID) (*Team, error) { return nil, nil }
func (repo *Datastore) Teams(id EventID) ([]*Team, error)                 { return nil, nil }

func (repo *Datastore) SubmitBallot(id EventID, ballot *Ballot) error   { return nil }
func (repo *Datastore) Ballots(id EventID) ([]*Ballot, error)           { return nil, nil }
func (repo *Datastore) LeastBallots(id EventID, n int) ([]*Team, error) { return nil, nil }
