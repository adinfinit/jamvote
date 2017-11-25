package event

import "context"

type Datastore struct {
	Context context.Context
}

func (repo *Datastore) List() ([]*Event, error)         { return nil, nil }
func (repo *Datastore) Create(event *Event) error       { return nil }
func (repo *Datastore) ByID(id EventID) (*Event, error) { return nil, nil }
func (repo *Datastore) Update(event *Event) error       { return nil }

func (repo *Datastore) CreateTeam(id EventID, team *Team) (TeamID, error) { return 0, nil }
func (repo *Datastore) TeamByID(id EventID, teamid TeamID) (*Team, error) { return nil, nil }
func (repo *Datastore) Teams(id EventID) ([]*Team, error)                 { return nil, nil }

func (repo *Datastore) SubmitBallot(id EventID, ballot *Ballot) error   { return nil }
func (repo *Datastore) Ballots(id EventID) ([]*Ballot, error)           { return nil, nil }
func (repo *Datastore) LeastBallots(id EventID, n int) ([]*Team, error) { return nil, nil }
