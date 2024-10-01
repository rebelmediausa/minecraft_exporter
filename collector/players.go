// Copyright 2010 Rebel Media
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build !noplayer
// +build !noplayer

package collector

import (
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rebelcore/minecraft_exporter/collector/utils"
)

type playerCollector struct {
	playersOnline *prometheus.Desc
	logger        *slog.Logger
}

func init() {
	registerCollector("players", defaultEnabled, NewPlayerCollector)
}

func NewPlayerCollector(logger *slog.Logger) (Collector, error) {
	const subsystem = "players"
	playersOnline := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "online"),
		"Minecraft players online.",
		[]string{"username", "dimension", "x", "y", "z", "experience"}, nil,
	)
	return &playerCollector{
		playersOnline: playersOnline,
		logger:        logger,
	}, nil
}

func getPlayerPosition(username string) []string {
	playerPOSFilter := regexp.MustCompile(`has the following entity data: (.*)`)
	getPlayerPOS := fmt.Sprintf("data get entity @p[name=%s] Pos", username)
	rawData := utils.GetRCON(getPlayerPOS)
	rawPlayerPOS := playerPOSFilter.FindStringSubmatch(rawData)[1]
	PlayerPOS := strings.Split(strings.Trim(strings.Trim(strings.ReplaceAll(rawPlayerPOS, " ", ""), "["), "]"), ",")

	return PlayerPOS
}

func getPlayerDimension(username string) string {
	playerDimensionFilter := regexp.MustCompile(`has the following entity data: (.*)`)
	getPlayerDimension := fmt.Sprintf("data get entity @p[name=%s] Dimension", username)
	rawData := utils.GetRCON(getPlayerDimension)
	rawPlayerDimension := playerDimensionFilter.FindStringSubmatch(rawData)[1]

	playerDimension := strings.Split(strings.Trim(strings.Trim(rawPlayerDimension, "\""), "\""), ":")

	return playerDimension[1]
}

func getPlayerXP(username string) string {
	playerXPFilter := regexp.MustCompile(`has the following entity data: (.*)`)
	getPlayerXP := fmt.Sprintf("data get entity @p[name=%s] XpLevel", username)
	rawData := utils.GetRCON(getPlayerXP)
	playerXP := playerXPFilter.FindStringSubmatch(rawData)[1]

	return playerXP
}

func (c *playerCollector) Update(ch chan<- prometheus.Metric) error {
	playerFilter := regexp.MustCompile(`players online: (.*)`)
	rawData := utils.GetRCON("list")

	if len(playerFilter.FindStringSubmatch(rawData)[1]) != 0 {
		rawPlayers := playerFilter.FindStringSubmatch(rawData)[1]
		Players := strings.Split(strings.Join(strings.Fields(rawPlayers), ""), ",")

		playerDimension := ""
		playerX := ""
		playerY := ""
		playerZ := ""
		playerXP := ""

		for _, player := range Players {
			c.logger.Debug("Minecraft user active", "Value", player)

			PlayerPOS := getPlayerPosition(player)
			playerX = strings.ReplaceAll(PlayerPOS[0], "d", "")
			playerY = strings.ReplaceAll(PlayerPOS[1], "d", "")
			playerZ = strings.ReplaceAll(PlayerPOS[2], "d", "")

			playerDimension = getPlayerDimension(player)

			playerXP = getPlayerXP(player)

			ch <- prometheus.MustNewConstMetric(c.playersOnline,
				prometheus.GaugeValue,
				1,
				player,
				playerDimension,
				playerX,
				playerY,
				playerZ,
				playerXP,
			)
		}
	}

	return nil
}
