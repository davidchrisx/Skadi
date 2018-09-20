package proccess

import (
	"log"
	"os"
	"strings"
	"fmt"
	"time"
	"strconv"

	"github.com/davidchrisx/skadi"

	"github.com/dotabuff/manta"
	"github.com/dotabuff/manta/dota"
)

func Run(file *os.File) (*skadi.Match, error) {
	Players := [10]skadi.Player{}
	Teams := [2]skadi.Team{}
	Match := &skadi.Match{}

	durSec := 0
	draftHero := [10]int32{}
	bannedHero := [12]int32{}
	postGame := false
	var gameStartTime time.Duration
	var gameEndTime time.Duration

	p, err := manta.NewStreamParser(file)
	if err != nil {
		return Match, err
	}

	p.OnEntity(func(pe *manta.Entity, pet manta.EntityOp) error {
		if postGame {
			return nil
		}
		if pe.GetClassName() == "CDOTAGamerulesProxy" {
			if v, ok := pe.GetFloat32("m_pGameRules.m_flGameStartTime"); ok {
				gameStartTime = time.Duration(v) * time.Second
			}
			if v, ok := pe.GetFloat32("m_pGameRules.m_flGameEndTime"); ok {
				gameEndTime = time.Duration(v) * time.Second
			}
			draftStage, _ := pe.GetInt32("m_pGameRules.m_nGameState")
			if draftStage == 2 {
				for i := 0; i < 10; i++ {
					draftHero[i], _ = pe.GetInt32(fmt.Sprintf("m_pGameRules.m_SelectedHeroes.%04d", i))
				}
				for i := 0; i < 12; i++ {
					bannedHero[i], _ = pe.GetInt32(fmt.Sprintf("m_pGameRules.m_BannedHeroes.%04d", i))
				}
			}
		}
		if pe.GetClassName() == "CDOTA_PlayerResource" {
			for i := 0; i < 10; i++ {
				heroID, _ := pe.GetInt32("m_vecPlayerTeamData.000" + strconv.Itoa(i) + ".m_nSelectedHeroID")
				if heroID > 0 {
					Players[i].HeroID = heroID
					Players[i].Level, _ = pe.GetInt32(fmt.Sprintf("m_vecPlayerTeamData.%04d.m_iLevel", i))
					Players[i].Assist, _ = pe.GetInt32(fmt.Sprintf("m_vecPlayerTeamData.%04d.m_iAssists", i))
					Players[i].Death, _ = pe.GetInt32(fmt.Sprintf("m_vecPlayerTeamData.%04d.m_iDeaths", i))
					Players[i].Kill, _ = pe.GetInt32(fmt.Sprintf("m_vecPlayerTeamData.%04d.m_iKills", i))
					Players[i].FirstBlood, _ = pe.GetInt32(fmt.Sprintf("m_vecPlayerTeamData.%04d.m_iFirstBloodClaimed", i))
					Players[i].TeamFight, _ = pe.GetFloat32(fmt.Sprintf("m_vecPlayerTeamData.%04d.m_flTeamFightParticipation", i))
				}
			}
		}
		if pe.GetClassName() == "CDOTA_DataRadiant" {
			for i := 0; i < 5; i++ {
				Players[i].Deny, _ = pe.GetInt32(fmt.Sprintf("m_vecDataTeam.%04d.m_iDenyCount", i))
				Players[i].LastHit, _ = pe.GetInt32(fmt.Sprintf("m_vecDataTeam.%04d.m_iLastHitCount", i))
				Players[i].ObsPlaced, _ = pe.GetInt32(fmt.Sprintf("m_vecDataTeam.%04d.m_iObserverWardsPlaced", i))
				Players[i].SenPlaced, _ = pe.GetInt32(fmt.Sprintf("m_vecDataTeam.%04d.m_iSentryWardsPlaced", i))
				Players[i].RunePickup, _ = pe.GetInt32(fmt.Sprintf("m_vecDataTeam.%04d.m_iRunePickups", i))
				Players[i].Gold, _ = pe.GetInt32(fmt.Sprintf("m_vecDataTeam.%04d.m_iTotalEarnedGold", i))
				Players[i].Exp, _ = pe.GetInt32(fmt.Sprintf("m_vecDataTeam.%04d.m_iTotalEarnedXP", i))
				Players[i].StunTime, _ = pe.GetFloat32(fmt.Sprintf("m_vecDataTeam.%04d.m_fStuns", i))
				Players[i].CampStacked, _ = pe.GetInt32(fmt.Sprintf("m_vecDataTeam.%04d.m_iCampsStacked", i))
				Players[i].TowerKill, _ = pe.GetInt32(fmt.Sprintf("m_vecDataTeam.%04d.m_iTowerKills", i))
				Players[i].RoshanKill, _ = pe.GetInt32(fmt.Sprintf("m_vecDataTeam.%04d.m_iRoshanKills", i))
			}
		}
		if pe.GetClassName() == "CDOTA_DataDire" {
			for j := 0; j < 5; j++ {
				i := j + 5
				Players[i].Deny, _ = pe.GetInt32(fmt.Sprintf("m_vecDataTeam.%04d.m_iDenyCount", j))
				Players[i].LastHit, _ = pe.GetInt32(fmt.Sprintf("m_vecDataTeam.%04d.m_iLastHitCount", j))
				Players[i].ObsPlaced, _ = pe.GetInt32(fmt.Sprintf("m_vecDataTeam.%04d.m_iObserverWardsPlaced", j))
				Players[i].SenPlaced, _ = pe.GetInt32(fmt.Sprintf("m_vecDataTeam.%04d.m_iSentryWardsPlaced", j))
				Players[i].RunePickup, _ = pe.GetInt32(fmt.Sprintf("m_vecDataTeam.%04d.m_iRunePickups", j))
				Players[i].Gold, _ = pe.GetInt32(fmt.Sprintf("m_vecDataTeam.%04d.m_iTotalEarnedGold", j))
				Players[i].Exp, _ = pe.GetInt32(fmt.Sprintf("m_vecDataTeam.%04d.m_iTotalEarnedXP", j))
				Players[i].StunTime, _ = pe.GetFloat32(fmt.Sprintf("m_vecDataTeam.%04d.m_fStuns", j))
				Players[i].CampStacked, _ = pe.GetInt32(fmt.Sprintf("m_vecDataTeam.%04d.m_iCampsStacked", j))
				Players[i].TowerKill, _ = pe.GetInt32(fmt.Sprintf("m_vecDataTeam.%04d.m_iTowerKills", j))
				Players[i].RoshanKill, _ = pe.GetInt32(fmt.Sprintf("m_vecDataTeam.%04d.m_iRoshanKills", j))
			}
		}
		return nil
	})

	p.Callbacks.OnCDemoFileInfo(func(m *dota.CDemoFileInfo) error{
		log.Println("Executing CDemoFileInfo callback")
		Info := m.GetGameInfo().GetDota()
		Teams[0].TeamName = Info.GetRadiantTeamTag()
		Teams[1].TeamName = Info.GetDireTeamTag()
		Match.MatchID = Info.GetMatchId()
		Match.Winner = getWinner(int(Info.GetGameWinner()))
		return nil
	})

	p.Callbacks.OnCDOTAMatchMetadataFile(func(m *dota.CDOTAMatchMetadataFile) error {
		log.Println("Executing OnCDOTAMatchMetadataFile")
		meta := m.GetMetadata().GetTeams()
		k := 0
		for i := 0; i < len(meta); i++ {
			players := meta[i].GetPlayers()
			for j := 0; j < len(players); j++ {
				player := players[j]
				Players[k].AccountID = player.GetAccountId()
				Snapshots := player.GetInventorySnapshot()
				LastSnap := Snapshots[len(Snapshots) - 1]
				Players[k].Items = LastSnap.GetItemId()
				k++
			}
		}
		return nil
	})
	combatLogs := []*dota.CMsgDOTACombatLogEntry{}
	p.Callbacks.OnCMsgDOTACombatLogEntry(func(m *dota.CMsgDOTACombatLogEntry) error {
		combatLogs = append(combatLogs, m)
		return nil
	})

	err = p.Start()

	heroIndex := []int{}
	for i := 0; i < 30 ; i ++ {
	s ,_ := p.LookupStringByIndex("CombatLogNames", int32(i))
	if strings.Contains(s ,"npc_dota_hero"){
			heroIndex = append(heroIndex, i)		
		}
	}

	for _, m := range combatLogs {
		source := m.GetDamageSourceName()
		target := m.GetTargetSourceName()
		value := m.GetValue()
		for i, hi := range heroIndex {
			if hi == int(source) {
				source = uint32(i)
				break
			}
		}
		for i, hi := range heroIndex {
			if hi == int(target) {
				target = uint32(i)
				break
			}
		}

		if m.GetType().String() == "DOTA_COMBATLOG_GAME_STATE" && value == 6 {
			postGame = true
			continue
		}

		if m.GetType().String() == "DOTA_COMBATLOG_DAMAGE" {
			if (m.GetIsAttackerHero() || m.GetIsAttackerIllusion()) && m.GetIsTargetHero() && !m.GetIsTargetIllusion() {
				Players[source].HeroDamage += value
				Players[target].DamageTaken += value
				continue
			}
			if (m.GetIsAttackerHero() || m.GetIsAttackerIllusion()) && m.GetIsTargetBuilding() {
				Players[source].TowerDamage += value
				continue
			}
		}

		if m.GetType().String() == "DOTA_COMBATLOG_HEAL" {
			if m.GetIsAttackerHero() && m.GetIsTargetHero() && !m.GetIsTargetIllusion() && !m.GetTargetIsSelf() {
				Players[source].HeroHealing += value
			}
		
		}
	}

	durSec = int(gameEndTime.Seconds() - gameStartTime.Seconds())
	durInt := time.Duration(durSec) * time.Second
	Match.Duration = durInt.String()
	for i := range Players {
		Players[i].GPM = countGPM(int(Players[i].Gold), durSec)
		Players[i].XPM = countXPM(int(Players[i].Exp), durSec)
		Players[i].KDA = countKDA(Players[i].Kill, Players[i].Death, Players[i].Assist)
		Teams[i/5].Players = append(Teams[i/5].Players, Players[i])
	}
	for i := range Teams {
		Match.Teams = append(Match.Teams, Teams[i])
	}
	for i := range draftHero {
		Match.Teams[i/5].Pick = append(Match.Teams[i/5].Pick, draftHero[i])
	}
	for i := range bannedHero {
		Match.Teams[i/6].Ban = append(Match.Teams[i/6].Ban, bannedHero[i])

	}
	return Match, err
}

func getWinner(t int) string {
	if t == 3 {
		return "Dire"
	} 
	return "Radiant"
}

func countKDA(k, d, a int32) float32 {
	if d == 0 {
		return float32(k + a)
	}
	return (float32(k) + float32(a)) / float32(d)
}

func countXPM(xp int, sec int) int {
	if sec == 0 {
		return 0
	}
	return xp * 60 / sec
}

func countGPM(gold int, sec int) int {
	if sec == 0 {
		return 0
	}
	return gold * 60 / sec
}
