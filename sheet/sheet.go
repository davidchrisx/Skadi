package sheet

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/davidchrisx/skadi"

	log "github.com/sirupsen/logrus"
	"gopkg.in/Iwark/spreadsheet.v2"
	"github.com/pkg/errors"
)

func Run(m *skadi.Match) error {
	service, err := spreadsheet.NewService()
	if err != nil {
		return errors.Wrap(err, "failed to get service")
	}
	ss, err := service.FetchSpreadsheet("1jSFPTXN2Eam75vaMTgBRw404Ca1LQcBbWPkL3sw-ULM")
	if err != nil {
		return errors.Wrap(err, "failed to fetch spreadsheet")
	}
	info, err := ss.SheetByTitle("Info")
	if err != nil {
		return errors.Wrap(err, "failed to fetch info sheet")
	}
	mSheet, err := ss.SheetByTitle("Matches")
	if err != nil {
		return errors.Wrap(err, "failed to fetch matches sheet")
	}
	r, err := strconv.Atoi(info.Columns[0][1].Value)
	if err != nil {
		return errors.Wrap(err, "failed to convert info columns")
	}
	err = FillMatch(mSheet, r+1, m)
	if err != nil {
		return errors.Wrap(err, "failed to fill match sheet")
	}
	log.Infoln("Fill Match OK")

	pSheet, err := ss.SheetByTitle("MatchPlayerData")
	if err != nil {
		return errors.Wrap(err, "failed to fetch match player sheet")
	}
	r, err = strconv.Atoi(info.Columns[1][1].Value)
	if err != nil {
		return err
	}
	err = FillPlayer(pSheet, r+1, m)
	if err != nil {
		return errors.Wrap(err, "failed to fill match player sheet")
	}
	log.Infoln("Fill Player OK")

	pfSheet, err := ss.SheetByTitle("FantasyAllMatch")
	if err != nil {
		return errors.Wrap(err, "failed to fetch fantasy sheet")
	}
	r, err = strconv.Atoi(info.Columns[2][1].Value)
	if err != nil {
		return err
	}
	err = FillFantasy(pfSheet, r+1, m)
	return err
}

func FillMatch(s *spreadsheet.Sheet, r int, m *skadi.Match) error {
	s.Update(r, 0, strconv.FormatUint(m.MatchID, 10))
	s.Update(r, 1, m.Winner)
	s.Update(r, 2, m.Duration)
	s.Update(r, 3, m.Teams[0].TeamName)
	s.Update(r, 4, m.Teams[1].TeamName)
	s.Update(r, 5, arrayToString(m.Teams[0].Pick, ", "))
	s.Update(r, 6, arrayToString(m.Teams[0].Ban, ", "))
	s.Update(r, 7, arrayToString(m.Teams[1].Pick, ", "))
	s.Update(r, 8, arrayToString(m.Teams[1].Ban, ", "))
	return s.Synchronize()
}

func FillFantasy(s *spreadsheet.Sheet, r int, m *skadi.Match) error {
	for i := 0; i < 10; i++ {
		T := m.Teams[i/5]
		P := T.Players[i%5]
		s.Update(r+i, 0, fmt.Sprintf("%d", P.AccountID))
		s.Update(r+i, 1, T.TeamName)
		s.Update(r+i, 2, fmt.Sprintf("%d", m.MatchID))
		s.Update(r+i, 3, fmt.Sprintf("%d", P.Kill))
		s.Update(r+i, 4, fmt.Sprintf("%d", P.Death))
		s.Update(r+i, 5, fmt.Sprintf("%d", P.LastHit+P.Deny))
		s.Update(r+i, 6, fmt.Sprintf("%d", P.GPM))
		s.Update(r+i, 7, fmt.Sprintf("%d", P.TowerKill))
		s.Update(r+i, 8, fmt.Sprintf("%d", P.RoshanKill))
		s.Update(r+i, 9, fmt.Sprintf("%.2f", P.TeamFight))
		s.Update(r+i, 10, fmt.Sprintf("%d", P.ObsPlaced))
		s.Update(r+i, 11, fmt.Sprintf("%d", P.CampStacked))
		s.Update(r+i, 12, fmt.Sprintf("%d", P.RunePickup))
		s.Update(r+i, 13, fmt.Sprintf("%d", P.FirstBlood))
		s.Update(r+i, 14, fmt.Sprintf("%.2f", P.StunTime))
		F := float32(3) + float32(P.Kill)*(0.3) - float32(P.Death)*(0.3) + float32(P.LastHit+P.Deny)*(0.003) + float32(P.GPM)*(0.002) + float32(P.TowerKill) + float32(P.RoshanKill) + float32(P.TeamFight)*(3) + float32(P.ObsPlaced)*(0.5) + float32(P.CampStacked)*(0.5) + float32(P.RunePickup)*(0.25) + float32(P.FirstBlood)*(4) + float32(P.StunTime)*(0.05)
		s.Update(r+i, 15, fmt.Sprintf("%.1f", F))
	}
	return s.Synchronize()
}

func FillPlayer(s *spreadsheet.Sheet, r int, m *skadi.Match) error {
	for i := 0; i < 10; i++ {
		T := m.Teams[i/5]
		P := T.Players[i%5]
		s.Update(r+i, 0, fmt.Sprintf("%d", m.MatchID))
		s.Update(r+i, 1, T.TeamName)
		s.Update(r+i, 2, fmt.Sprintf("%d", P.AccountID))
		s.Update(r+i, 3, fmt.Sprintf("%d", P.HeroID))
		s.Update(r+i, 4, uArrayToString(P.Items, ", "))
		s.Update(r+i, 5, fmt.Sprintf("%d", P.Level))
		s.Update(r+i, 6, fmt.Sprintf("%d", P.Kill))
		s.Update(r+i, 7, fmt.Sprintf("%d", P.Death))
		s.Update(r+i, 8, fmt.Sprintf("%d", P.Assist))
		s.Update(r+i, 9, fmt.Sprintf("%d", P.LastHit))
		s.Update(r+i, 10, fmt.Sprintf("%d", P.Deny))
		s.Update(r+i, 11, fmt.Sprintf("%d", P.HeroDamage))
		s.Update(r+i, 12, fmt.Sprintf("%d", P.TowerDamage))
		s.Update(r+i, 13, fmt.Sprintf("%d", P.DamageTaken))
		s.Update(r+i, 14, fmt.Sprintf("%d", P.HeroHealing))
		s.Update(r+i, 15, fmt.Sprintf("%d", P.GPM))
		s.Update(r+i, 16, fmt.Sprintf("%d", P.XPM))
		s.Update(r+i, 17, fmt.Sprintf("%.2f", P.TeamFight))
		s.Update(r+i, 18, fmt.Sprintf("%d", P.FirstBlood))
		s.Update(r+i, 19, fmt.Sprintf("%d", P.RunePickup))
		s.Update(r+i, 20, fmt.Sprintf("%d", P.ObsPlaced))
		s.Update(r+i, 21, fmt.Sprintf("%d", P.SenPlaced))
		s.Update(r+i, 22, fmt.Sprintf("%d", P.CampStacked))
		s.Update(r+i, 23, fmt.Sprintf("%d", P.TowerKill))
		s.Update(r+i, 24, fmt.Sprintf("%d", P.RoshanKill))
		s.Update(r+i, 25, fmt.Sprintf("%.2f", P.StunTime))
	}
	return s.Synchronize()
}

func uArrayToString(A []uint32, delim string) string {

	var buffer bytes.Buffer
	for i := 0; i < len(A); i++ {
		buffer.WriteString(strconv.FormatUint(uint64(A[i]), 10))
		if i != len(A)-1 {
			buffer.WriteString(delim)
		}
	}

	return buffer.String()
}

func arrayToString(A []int32, delim string) string {

	var buffer bytes.Buffer
	for i := 0; i < len(A); i++ {
		buffer.WriteString(strconv.FormatInt(int64(A[i]), 10))
		if i != len(A)-1 {
			buffer.WriteString(delim)
		}
	}

	return buffer.String()
}
