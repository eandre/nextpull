package boss

import "github.com/eandre/lunar-shim/hbd"

type EncounterID int

type Raid struct {
	Name       string
	InstanceID hbd.InstanceID
	Bosses     []*Boss
}

type Boss struct {
	EncounterID    EncounterID
	Name           string
	X, Y           hbd.WorldCoord
	Raid           *Raid
	CustomDuration float32
}

func (b *Boss) FightDuration() float32 {
	if b.CustomDuration != 0 {
		return b.CustomDuration
	}
	return DefaultDuration
}

const DefaultDuration = 5 * 60

var Raids = []*Raid{
	&Raid{
		Name:       "Hellfire Citadel",
		InstanceID: 1448,
		Bosses: []*Boss{
			&Boss{
				Name:        "Hellfire Assault",
				EncounterID: 1778,
				X:           -767.1,
				Y:           3972.5,
			},
			&Boss{
				Name:           "Shadow-Lord Iskar",
				EncounterID:    1788,
				X:              2497.2,
				Y:              4040,
				CustomDuration: 4 * 60,
			},
		},
	},
}

var UnknownRaid = &Raid{
	Name:       "Unknown Raid",
	InstanceID: 0,
	Bosses: []*Boss{
		&Boss{
			Name:        "Unknown",
			EncounterID: 0,
			X:           0,
			Y:           0,
		},
	},
}

var UnknownBoss = UnknownRaid.Bosses[0]

func Find(eid EncounterID) *Boss {
	for _, raid := range Raids {
		for _, boss := range raid.Bosses {
			if boss.EncounterID == eid {
				return boss
			}
		}
	}
	return UnknownBoss
}

func FindClosest() *Boss {
	x, y, inst := hbd.PlayerWorldPosition()
	for _, raid := range Raids {
		if raid.InstanceID == inst {
			var closestDist float32 = -1
			var closestBoss *Boss
			for _, boss := range raid.Bosses {
				dist := hbd.WorldDistance(inst, x, y, boss.X, boss.Y)
				if closestBoss == nil || dist < closestDist {
					closestBoss = boss
					closestDist = dist
				}
			}
			return closestBoss
		}
	}
	return nil
}

func init() {
	// Populate back-references
	for _, raid := range Raids {
		for _, boss := range raid.Bosses {
			boss.Raid = raid
		}
	}
	UnknownBoss.Raid = UnknownRaid
}
