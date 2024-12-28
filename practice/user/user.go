package user

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/akmalfairuz/df-practice/practice/lang"
	"github.com/akmalfairuz/df-practice/practice/model"
	"github.com/df-mc/atomic"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/scoreboard"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"golang.org/x/text/language"
	"strings"
	"sync"
	"time"
)

type User struct {
	id int

	p    *world.EntityHandle
	s    *session.Session
	conn session.Conn

	xuid string
	name string

	rankName      string
	rankExpiredAt int64

	lang language.Tag

	closed atomic.Bool

	currentFFAArena atomic.Value[any]
	currentGame     atomic.Value[any]

	lastWhisperFromXUID atomic.Value[string]

	clicksMu sync.Mutex
	clicks   []time.Time

	comboCounter             int
	comboCounterMu           sync.Mutex
	lastComboCounterModified time.Time

	lastHitReachDistance         float64
	lastHitReachDistanceModified time.Time
	lastHitReachDistanceMu       sync.Mutex
}

func New(p *player.Player) *User {

	s := playerSession(p)
	conn := sessionConn(s)

	locale := lang.ToLangTag(p.Locale().String())

	u := &User{
		p:      p.H(),
		xuid:   p.XUID(),
		name:   p.Name(),
		s:      s,
		conn:   conn,
		lang:   locale,
		clicks: make([]time.Time, 0),
	}

	u.lastWhisperFromXUID.Store("")
	return u
}

func (u *User) Lang() language.Tag {
	return u.lang
}

func (u *User) Load() error {
	userData, err := u.loadUserData()
	if err != nil {
		return fmt.Errorf("failed to load user data: %w", err)
	}
	u.id = userData.ID
	u.rankName = userData.RankName
	u.rankExpiredAt = userData.RankExpireAt

	_ = u.SynchronizeLastSeen()
	return nil
}

func (u *User) loadUserData() (model.User, error) {
	userData, err := userRepository.FindByXUID(u.xuid)
	if errors.Is(err, sql.ErrNoRows) {
		if _, err := userRepository.Create(model.CreateUser{
			XUID:        u.xuid,
			DisplayName: u.name,
		}); err != nil {
			return model.User{}, fmt.Errorf("failed to create user: %w", err)
		}

		userData, err = userRepository.FindByXUID(u.xuid)
	}

	if err != nil {
		return model.User{}, fmt.Errorf("failed to find user: %w", err)
	}

	if userData.DisplayName != u.name {
		if err := userRepository.SetDisplayName(userData.ID, u.name); err != nil {
			return model.User{}, fmt.Errorf("detected player changed username, but failed to set display name: %w", err)
		}
	}

	return userData, nil
}

func (u *User) SynchronizeLastSeen() error {
	return userRepository.SynchronizeLastSeen(u.id)
}

func (u *User) Disconnect(message string) {
	u.s.Disconnect(message)
}

func (u *User) Closed() bool {
	return u.closed.Load()
}

func (u *User) Close() error {
	if u.closed.CAS(false, true) {
		RemoveByXUID(u.xuid)
		_ = u.conn.Close()
	}
	return nil
}

func (u *User) WritePacket(pk packet.Packet) error {
	return u.conn.WritePacket(pk)
}

func (u *User) Messaget(translationName string, args ...any) {
	u.s.SendMessage(lang.Translatef(u.lang, translationName, args...))
}

func (u *User) SendScoreboard(lines []string) {
	lines2 := make([]string, 0, len(lines))
	lines2 = append(lines2, "")
	for _, line := range lines {
		lines2 = append(lines2, line)
	}
	lines2 = append(lines2, "")
	lines2 = append(lines2, u.Translatef("scoreboard.footer"))

	u.SendScoreboardRaw(u.Translatef("scoreboard.title"), lines2)
}

func (u *User) SendScoreboardRaw(title string, lines []string) {
	sb := scoreboard.New(title)
	blankLinesCount := 0
	for _, line := range lines {
		if line == "" {
			_, _ = sb.WriteString(" " + strings.Repeat("§r", blankLinesCount+1))
			blankLinesCount++
			continue
		}
		if line == "<empty>" {
			continue
		}
		_, _ = sb.WriteString(line)
	}
	u.s.RemoveScoreboard()
	u.s.SendScoreboard(sb)
}

func (u *User) Session() *session.Session {
	return u.s
}

func (u *User) Translatef(translationName string, args ...any) string {
	return lang.Translatef(u.lang, translationName, args...)
}

func (u *User) XUID() string {
	return u.xuid
}

func (u *User) Name() string {
	return u.name
}

func (u *User) CurrentFFAArena() any {
	return u.currentFFAArena.Load()
}

func (u *User) SetCurrentFFAArena(a any) {
	u.currentFFAArena.Store(a)
}

func (u *User) EntityHandle() *world.EntityHandle {
	return u.p
}

func (u *User) Conn() session.Conn {
	return u.conn
}

func (u *User) EntityRuntimeID() uint64 {
	return 1
}

func (u *User) OnReceiveWhisper(sender *User, message string) {
	u.lastWhisperFromXUID.Store(sender.XUID())

	u.Messaget("whisper.received", sender.Name(), message)
}

func (u *User) OnSendWhisper(target *User, message string) {
	u.Messaget("whisper.sent", target.Name(), message)
}

func (u *User) ReplyWhisperToXUID() string {
	return u.lastWhisperFromXUID.Load()
}

func (u *User) CurrentGame() any {
	return u.currentGame.Load()
}

func (u *User) SetCurrentGame(g any) {
	u.currentGame.Store(g)
}

func (u *User) Ping() int {
	return int(u.s.Latency().Milliseconds())
}

func (u *User) RankName() string {
	return u.rankName
}

func (u *User) Player(tx *world.Tx) (*player.Player, bool) {
	ent, ok := u.EntityHandle().Entity(tx)
	if !ok {
		return nil, false
	}
	return ent.(*player.Player), true
}

// RemoveOldClicks removes all clicks that are older than 1 second.
func (u *User) RemoveOldClicks() {
	u.clicksMu.Lock()
	defer u.clicksMu.Unlock()
	newClicks := make([]time.Time, 0)

	for _, click := range u.clicks {
		if time.Since(click) < time.Second {
			newClicks = append(newClicks, click)
		}
	}

	u.clicks = newClicks
}

func (u *User) HandleClientClick() {
	u.clicksMu.Lock()
	u.clicks = append(u.clicks, time.Now())
	u.clicksMu.Unlock()
}

func (u *User) CPS() int {
	clicks := 0
	now := time.Now()
	u.clicksMu.Lock()
	for i := 0; i < len(u.clicks); i++ {
		if now.Sub(u.clicks[i]) <= time.Second {
			clicks++
		}
	}
	u.clicksMu.Unlock()
	return clicks
}

func (u *User) ComboCounter() int {
	return u.comboCounter
}

func (u *User) AddComboCounter() {
	u.comboCounterMu.Lock()
	u.comboCounter++
	u.lastComboCounterModified = time.Now()
	u.comboCounterMu.Unlock()
}

func (u *User) ResetComboCounter() {
	u.comboCounterMu.Lock()
	u.comboCounter = 0
	u.lastComboCounterModified = time.Now()
	u.comboCounterMu.Unlock()
}

func (u *User) LastComboCounterModified() time.Time {
	return u.lastComboCounterModified
}

func (u *User) LastHitReachDistance() float64 {
	u.lastHitReachDistanceMu.Lock()
	defer u.lastHitReachDistanceMu.Unlock()
	if time.Since(u.lastHitReachDistanceModified) > time.Second {
		return 0
	}
	return u.lastHitReachDistance
}

func (u *User) SetLastHitReachDistance(d float64) {
	u.lastHitReachDistanceMu.Lock()
	u.lastHitReachDistance = d
	u.lastHitReachDistanceModified = time.Now()
	u.lastHitReachDistanceMu.Unlock()
}

func (u *User) LastHitReachDistanceModified() time.Time {
	return u.lastHitReachDistanceModified
}

func (u *User) SendPVPInfoTip() {
	cps := u.CPS()
	combo := u.ComboCounter()
	reach := u.LastHitReachDistance()
	if reach > 3.0 {
		reach = 3.0
	}

	if cps == 0 && combo == 0 && reach <= 0.1 {
		return
	}
	u.Session().SendTip(u.Translatef("pvp.info.tip", cps, combo, reach))
}
