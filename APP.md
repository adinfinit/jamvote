# JamVote - Application Overview

JamVote is a web application for organizing and voting on game jam entries. It manages the full lifecycle of a jam: event creation, team registration, game submission, voting, and result reveal.

## User Roles

- **Visitor** -- unauthenticated user, can view public event listings and results.
- **Jammer** -- authenticated user approved for an event. Can create teams, submit games, and vote.
- **Judge** -- authenticated user marked as judge for an event. Votes are weighted separately.
- **Admin** -- can create/edit events, manage users, approve jammers/judges, reveal results.

## Authentication

Login is via Google OAuth2. In development mode a simple username form is available instead. After first login a user record is created and the user is prompted to fill in their name.

## Pages

### Public / Site

| Route | Page | Description |
|---|---|---|
| `/` | Event List | Home page. Lists active events and finished events grouped by year. |
| `/static/*` | Static Assets | CSS (`main.css`), JS (`jam.js`), favicon. |
| `/about` | About | Markdown page describing JamVote. |
| `/about/jamming` | Jamming Guide | How to register, submit, and vote. |
| `/about/scoring` | Scoring Guide | Explanation of the five scoring aspects and statistics. |
| `/about/organizing` | Organizing Guide | Best practices for running a jam event. |

### Authentication

| Route | Page | Description |
|---|---|---|
| `/user/login` | Login | Shows login provider links. In dev mode shows a username form. |
| `/auth/login` | OAuth Start | Redirects to Google OAuth2 consent screen. |
| `/auth/callback` | OAuth Callback | Exchanges auth code for token, stores credentials in session. |
| `/auth/development-login` | Dev Login | POST-only. Creates a dev session with the given username. |
| `/auth/logout` | Logout | Clears session, redirects to home. |
| `/user/logged-in` | Post-Login | Redirects new users to profile edit, returning users to home. |
| `/user/logout` | Logout | Triggers auth logout. |

### User / Profile

| Route | Page | Description |
|---|---|---|
| `/user` | My Profile | Redirects logged-in user to their edit page. |
| `/users` | User List | Table of all registered users (requires login). |
| `/user/{userid}` | Profile | Displays name, social links, and teams the user has been on. |
| `/user/{userid}/edit` | Edit Profile | Form: name, email, Facebook, GitHub. Admins can toggle admin flag. |

### Event Management

| Route | Page | Description |
|---|---|---|
| `/event/create` | Create Event | Admin form: name, slug, theme, judge percentage, start/end dates, info text. |
| `/event/{eventid}` | Event Dashboard | Main event page. Shows theme, countdown timers (start, end, voting open/close), team submission warnings, and a "Vote" button when voting is active. |
| `/event/{eventid}/edit` | Edit Event | Admin form: theme, judge%, lifecycle flags (registration, voting, closed, revealed), dates, info, organizers. |
| `/event/{eventid}/jammers` | Manage Jammers | Admin page. Checkbox table to approve/remove jammers and toggle judge status. |
| `/event/{eventid}/linking` | User Linking | Admin page. Resolves mismatches between entered member names and registered users; shows unregistered and unapproved members. |
| `/event/{eventid}/linking-approve-all` | Auto-Approve | Admin action. Automatically links and approves all matching users. |

### Teams

| Route | Page | Description |
|---|---|---|
| `/event/{eventid}/teams` | Team List | Table of all teams with members and game names. |
| `/event/{eventid}/team/create` | Create Team | Form: team name (with randomizer), game name (with randomizer), members (up to 9). |
| `/event/{eventid}/team/{teamid}` | Team Detail | Shows game info, member list, download/jam/facebook links, voting results with violin plots, and voter comments. |
| `/event/{eventid}/team/{teamid}/edit` | Edit Team | Form: team name, game details (name, info, links), member list, non-competing flag. Admin can delete team. |
| `/event/{eventid}/team/{teamid}/delete` | Delete Team | Admin-only. Removes team during registration phase. |

### Voting

| Route | Page | Description |
|---|---|---|
| `/event/{eventid}/voting` | Voting Queue | Shows the user's pending games to vote on and a table of completed votes with scores. |
| `/event/{eventid}/fill-queue` | Fill Queue | Assigns next games to vote on (3 initially, then 1 at a time). Redirects to voting page. |
| `/event/{eventid}/vote/{teamid}` | Vote | Five slider inputs (Theme, Enjoyment, Aesthetics, Innovation, Bonus) each with an optional comment, plus a submit button. Votes can be edited until voting closes. |
| `/event/{eventid}/progress` | Voting Progress | Bar charts showing aggregate and per-team voting completion. Admin gets auto-refresh. |

### Results

| Route | Page | Description |
|---|---|---|
| `/event/{eventid}/reveal` | Reveal | Admin triggers reveal. Shows top 5 results in a dramatic countdown display. |
| `/event/{eventid}/results` | Results | Three ranked tables: overall, jammer average, judge average. Non-competing entries shown separately. |
| `/event/{eventid}/ballots.csv` | Export Ballots | Admin-only CSV download of all ballots. |

## Data Models

### User

| Field | Type | Description |
|---|---|---|
| ID | int64 | Unique identifier. |
| Name | string | Display name. |
| Email | string | Email address. |
| Admin | bool | Site-wide admin flag. |
| Facebook | string | Facebook profile URL. |
| Github | string | GitHub profile URL. |
| NewUser | bool | True until profile is first saved. |

### Event

| Field | Type | Description |
|---|---|---|
| ID | string | URL slug (lowercase alphanumeric + hyphens). |
| Name | string | Display name. |
| Theme | string | Jam theme. |
| Info | string | Description / rules (plain text, rendered as paragraphs). |
| Created | time | Creation timestamp. |
| StartTime | time | Jam start. |
| EndTime | time | Jam end. |
| VotingOpens | time | When voting opens. |
| VotingCloses | time | When voting closes. |
| JudgePercentage | float64 | Weight of judge scores (0-100). |
| Registration | bool | Team registration open. |
| Voting | bool | Voting open. |
| Closed | bool | Voting closed. |
| Revealed | bool | Results publicly visible. |
| Organizers | []UserID | Users who organize this event. |
| Jammers | []UserID | Approved jammers. |
| Judges | []UserID | Designated judges. |

### Team

| Field | Type | Description |
|---|---|---|
| EventID | string | Parent event. |
| ID | int64 | Unique team identifier. |
| Name | string | Team name. |
| Members | []Member | List of members (each has UserID and Name). |
| Game.Name | string | Submitted game name. |
| Game.Info | string | Game description. |
| Game.Noncompeting | bool | Excluded from rankings. |
| Game.Link.Jam | string | Jam page URL. |
| Game.Link.Download | string | Download URL. |
| Game.Link.Facebook | string | Facebook post URL. |

### Ballot

| Field | Type | Description |
|---|---|---|
| Voter | UserID | Who cast this vote. |
| Team | TeamID | Which team is being voted on. |
| Index | int64 | Queue order. |
| Completed | bool | Has been submitted. |
| Theme | Aspect | Theme score (1-5) + comment. |
| Enjoyment | Aspect | Enjoyment score (1-5) + comment. |
| Aesthetics | Aspect | Aesthetics score (1-5) + comment. |
| Innovation | Aspect | Innovation score (1-5) + comment. |
| Bonus | Aspect | Bonus score (0-2.5) + comment. |
| Overall | Aspect | Calculated: (Theme+Enjoyment+Aesthetics+Innovation+Bonus) / 4.5, clamped to 1-5. |

## Scoring System

Each voter rates five aspects. Overall is a weighted average: `(Theme + Enjoyment + Aesthetics + Innovation + Bonus) / 4.5`, clamped to 1-5.

When an event has a judge percentage `p` (0-100), final scores blend jammer and judge averages: `final = jammer_avg * (1 - p/100) + judge_avg * (p/100)`.

Self-votes (voting on your own team) are excluded from results.

## Event Lifecycle

1. **Setup** -- Admin creates event, sets theme, dates, judge percentage.
2. **Registration** -- `Registration=true`. Users sign up, create teams, submit games.
3. **Linking** -- Admin resolves user-team mismatches, approves jammers/judges.
4. **Voting** -- `Voting=true`. Jammers and judges vote on games through a queue system (3 initial assignments, then 1 at a time).
5. **Closed** -- `Closed=true`. No more votes accepted. Admin reviews progress.
6. **Reveal** -- `Revealed=true`. Top 5 shown in dramatic reveal page, then full results become public.

## Template Structure

All pages extend a shared layout (`site-common.html`) that provides:

- Navigation bar with links to home, about, login/logout
- Event sub-navigation (when viewing an event): dashboard, teams, voting, progress, results, admin links
- Flash messages (errors and notifications)
- Admin bar with event lifecycle toggles
- Footer

Templates use Go's `html/template` with custom functions for date formatting, markdown rendering, math operations, and violin plot SVG generation.

## Tech Stack

- **Backend:** Go, `net/http` with gorilla/mux router, gorilla/sessions for cookie sessions
- **Database:** Google Cloud Datastore
- **Auth:** Google OAuth2 (with dev mode fallback)
- **Frontend:** Server-rendered HTML templates, vanilla CSS, vanilla JS
- **Hosting:** Google Cloud Run
- **Secrets:** Google Secret Manager (with env var fallback)
