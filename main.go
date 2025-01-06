package main

import (
	"fmt"
	"github.com/akmalfairuz/df-practice/internal/meta"
	"github.com/akmalfairuz/df-practice/practice"
	"github.com/akmalfairuz/df-practice/practice/command"
	"github.com/akmalfairuz/df-practice/practice/config"
	"github.com/akmalfairuz/legacy-version/legacyver"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft"
	"log/slog"
)

func main() {
	log := slog.Default()

	userConfig := server.DefaultConfig()
	userConfig.Server.Name = config.Get().Server.Name
	userConfig.Server.DisableJoinQuitMessages = true

	userConfig.Players.SaveData = false
	userConfig.Players.MaxCount = config.Get().Server.MaxPlayers
	userConfig.Network.Address = config.Get().Server.ListenAddress
	userConfig.Server.AuthEnabled = config.Get().Server.AuthEnabled

	serverConfig, err := userConfig.Config(log)
	if err != nil {
		panic(err)
	}

	serverConfig.Allower = practice.NewAllower()
	serverConfig.Listeners = []func(conf server.Config) (server.Listener, error){
		func(conf server.Config) (server.Listener, error) {
			cfg := minecraft.ListenConfig{
				ErrorLog:               log,
				MaximumPlayers:         conf.MaxPlayers,
				StatusProvider:         statusProvider{name: conf.Name},
				AuthenticationDisabled: conf.AuthDisabled,
				ResourcePacks:          conf.Resources,
				Biomes:                 biomes(),
				TexturePacksRequired:   conf.ResourcesRequired,
				AcceptedProtocols: []minecraft.Protocol{
					legacyver.New748(),
					legacyver.New729(),
					legacyver.New712(),
					legacyver.New686(),
					legacyver.New685(),
					legacyver.New671(),
				},
			}
			l, err := cfg.Listen("raknet", userConfig.Network.Address)
			if err != nil {
				return nil, fmt.Errorf("create minecraft listener: %w", err)
			}
			conf.Log.Info("Server running on %v", "addr", l.Addr())
			return listener{l}, nil
		},
	}
	srv := serverConfig.New()

	meta.Set("server", srv)

	cmd.Register(cmd.New("ban", "Ban a player from the server.", nil, command.Ban{}))
	cmd.Register(cmd.New("pardon", "Pardon a banned player.", []string{"unban"}, command.Pardon{}))
	cmd.Register(cmd.New("whisper", "Send a private message to a player.", []string{"w", "msg"}, command.Whisper{}))
	cmd.Register(cmd.New("reply", "Reply to the last private message received.", []string{"r"}, command.Reply{}))
	cmd.Register(cmd.New("lobby", "Teleport to the lobby.", []string{"hub"}, command.Lobby{}))
	cmd.Register(cmd.New("gamemode", "Change the game mode of a player.", []string{"gm"}, command.GameMode{}))
	cmd.Register(cmd.New("teleport", "Teleport to a player or location.", []string{"tp"}, command.TeleportToTarget{}, command.TeleportToPos{}, command.TeleportTargetToTarget{}, command.TeleportTargetToPos{}))
	cmd.Register(cmd.New("duel", "Send a duel request to a player.", nil, command.Duel{}))

	pr := practice.New(log, srv)
	pr.Run()
}

// listener is a Listener implementation that wraps around a minecraft.Listener so that it can be listened on by
// Server.
type listener struct {
	*minecraft.Listener
}

// Accept blocks until the next connection is established and returns it. An error is returned if the Listener was
// closed using Close.
func (l listener) Accept() (session.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	return conn.(session.Conn), err
}

// Disconnect disconnects a connection from the Listener with a reason.
func (l listener) Disconnect(conn session.Conn, reason string) error {
	return l.Listener.Disconnect(conn.(*minecraft.Conn), reason)
}

// statusProvider handles the way the server shows up in the server list. The
// online players and maximum players are not changeable from outside the
// server, but the server name may be changed at any time.
type statusProvider struct {
	name string
}

// ServerStatus returns the player count, max players and the server's name as
// a minecraft.ServerStatus.
func (s statusProvider) ServerStatus(playerCount, maxPlayers int) minecraft.ServerStatus {
	return minecraft.ServerStatus{
		ServerName:  s.name,
		PlayerCount: playerCount,
		MaxPlayers:  maxPlayers,
	}
}

// ashyBiome represents a biome that has any form of ash.
type ashyBiome interface {
	// Ash returns the ash and white ash of the biome.
	Ash() (ash float64, whiteAsh float64)
}

// sporingBiome represents a biome that has blue or red spores.
type sporingBiome interface {
	// Spores returns the blue and red spores of the biome.
	Spores() (blueSpores float64, redSpores float64)
}

// biomes builds a mapping of all biome definitions of the server, ready to be set in the biomes field of the server
// listener.
func biomes() map[string]any {
	definitions := make(map[string]any)
	for _, b := range world.Biomes() {
		definition := map[string]any{
			"name_hash":   b.String(), // This isn't actually a hash despite what the field name may suggest.
			"temperature": float32(b.Temperature()),
			"downfall":    float32(b.Rainfall()),
			"rain":        b.Rainfall() > 0,
		}
		if a, ok := b.(ashyBiome); ok {
			ash, whiteAsh := a.Ash()
			definition["ash"], definition["white_ash"] = float32(ash), float32(whiteAsh)
		}
		if s, ok := b.(sporingBiome); ok {
			blueSpores, redSpores := s.Spores()
			definition["blue_spores"], definition["red_spores"] = float32(blueSpores), float32(redSpores)
		}
		definitions[b.String()] = definition
	}
	return definitions
}
