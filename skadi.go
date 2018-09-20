package skadi

type Match struct {
	MatchID 	uint64
	Winner 		string
	Duration	string
	Teams		[]Team
}

type Team struct {
	TeamName	string
	Side		string
	Pick		[]int32
	Ban			[]int32
	Players		[]Player
}

type Player struct {
	IngameName	string
	SteamID		uint32
	AccountID	uint32
	HeroID		int32
	Items		[]uint32
	Level		int32
	Kill		int32
	Death		int32
	Assist		int32
	Deny		int32
	LastHit		int32
	HeroHealing	uint32
	HeroDamage	uint32
	DamageTaken	uint32
	TowerDamage	uint32
	TeamFight	float32
	FirstBlood	int32
	RunePickup	int32
	ObsPlaced	int32
	SenPlaced	int32
	Gold		int32
	Exp			int32
	StunTime	float32
	CampStacked	int32
	RoshanKill	int32
	TowerKill	int32
	KDA			float32
	GPM			int
	XPM			int
}