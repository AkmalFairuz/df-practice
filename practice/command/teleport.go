package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// TeleportToPos is a command that teleports the sender to a specific position.
type TeleportToPos struct {
	onlyAdmin

	Pos mgl64.Vec3 `name:"pos"`
}

// TeleportToTarget is a command that teleports the sender to a specific target.
type TeleportToTarget struct {
	onlyAdmin

	Target onlineTarget `name:"target"`
}

// TeleportTargetToTarget is a command that teleports a target to another target.
type TeleportTargetToTarget struct {
	onlyAdmin

	From onlineTarget `name:"from"`
	To   onlineTarget `name:"to"`
}

// TeleportTargetToPos is a command that teleports a target to a specific position.
type TeleportTargetToPos struct {
	onlyAdmin

	Target onlineTarget `name:"target"`
	Pos    mgl64.Vec3   `name:"pos"`
}

func (t TeleportToPos) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	p := s.(*player.Player)
	p.Teleport(t.Pos)
	messaget(s, "command.teleport.success", p.Name(), formatPos(t.Pos))
}

func (t TeleportToTarget) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	fromWorld := tx.World()
	targetWorld := t.Target.User().World()

	if fromWorld == targetWorld {
		ent, ok := t.Target.User().EntityHandle().Entity(tx)
		if ok {
			p := ent.(*player.Player)
			p.Teleport(s.(*player.Player).Position())
			messaget(s, "command.teleport.success", p.Name(), t.Target.User().Name())
			return
		}
		messaget(s, "command.an.error.occurred")
		return
	}

	tx.RemoveEntity(s.(world.Entity))

	targetWorld.Exec(func(tx *world.Tx) {
		tx.AddEntity(s.(world.Entity).H())

		ent, ok := t.Target.User().EntityHandle().Entity(tx)
		if ok {
			p := ent.(*player.Player)
			p.Teleport(s.(*player.Player).Position())
			messaget(s, "command.teleport.success", p.Name(), t.Target.User().Name())
			return
		}

		messaget(s, "command.an.error.occurred")
	})
}

func (t TeleportTargetToTarget) Run(s cmd.Source, o *cmd.Output, _ *world.Tx) {
	if t.To.User().World() == t.From.User().World() {
		t.To.User().World().Exec(func(tx *world.Tx) {
			fromEnt, ok := t.From.User().EntityHandle().Entity(tx)
			if !ok {
				return
			}

			toEnt, ok := t.To.User().EntityHandle().Entity(tx)
			if !ok {
				return
			}

			toEnt.(*player.Player).Teleport(fromEnt.(*player.Player).Position())
			messaget(s, "command.teleport.success", t.To.User().Name(), t.From.User().Name())
		})
		return
	}

	t.From.User().World().Exec(func(tx *world.Tx) {
		fromEnt, ok := t.From.User().EntityHandle().Entity(tx)
		if !ok {
			return
		}

		tx.RemoveEntity(fromEnt)
		t.To.User().World().Exec(func(tx *world.Tx) {
			tx.AddEntity(fromEnt.H())

			toEnt, ok := t.To.User().EntityHandle().Entity(tx)
			if !ok {
				return
			}

			toEnt.(*player.Player).Teleport(fromEnt.(*player.Player).Position())
			messaget(s, "command.teleport.success", t.To.User().Name(), t.From.User().Name())
		})
	})
}

func (t TeleportTargetToPos) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	t.Target.ExecutePlayer(func(p *player.Player, ok bool) {
		if !ok {
			return
		}

		p.Teleport(t.Pos)

		messaget(s, "command.teleport.success", p.Name(), formatPos(t.Pos))
	})
}
