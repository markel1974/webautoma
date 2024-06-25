/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package invaders

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strconv"
	"strings"
)

const (
	highScoreFilename  = "hs"
	highScoreSeparator = ":"
	maxHighScores      = 5
)

type HighScore struct {
	score int
	name  string
}

type ByScore []*HighScore

func (a ByScore) Len() int           { return len(a) }
func (a ByScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByScore) Less(i, j int) bool { return a[i].score < a[j].score }

func (g *Invaders) loadHighScores() {
	data, err := ioutil.ReadFile(highScoreFilename)
	if err != nil {
		log.Fatalln(err)
	}
	lines := strings.Split(string(data), "\n")
	if len(lines) == 0 {
		return
	}
	for _, l := range lines {
		parts := strings.Split(l, highScoreSeparator)
		if len(parts) != 2 {
			log.Println("high scores file has been corrupted - please correct/delete it")
			continue
		} else if len(parts[0]) < 3 || len(parts[0]) > 10 {
			log.Println("high score file has been corrupted (name too long/short) - please correct/delete it")
			continue
		}
		if i, err := strconv.Atoi(parts[1]); err == nil {
			if i < 0 {
				log.Println("negative high score found - data corrupted - please correct/delete the hs file")
				continue
			}
			g.highscores = append(g.highscores, &HighScore{i, parts[0]})
		} else {
			log.Fatalln(err)
		}
	}

	sort.Sort(sort.Reverse(ByScore(g.highscores)))
	if len(g.highscores) > 5 {
		g.highscores = g.highscores[:5]
	}
}

func (g *Invaders) checkHighScores() {
	if len(g.highscores) < maxHighScores || g.player.score > g.highscores[len(g.highscores)-1].score {
		name := "TODO NAME!!!"
		g.highscores = append(g.highscores, &HighScore{g.player.score, name})
		sort.Sort(sort.Reverse(ByScore(g.highscores)))
		if len(g.highscores) > maxHighScores {
			g.highscores = append([]*HighScore(nil), g.highscores[:maxHighScores]...)
		}

		data := ""
		for i, score := range g.highscores {
			data += fmt.Sprintf("%s%s%d", score.name, highScoreSeparator, score.score)
			if i != len(g.highscores)-1 {
				data += "\n"
			}
		}
		_ = ioutil.WriteFile(highScoreFilename, []byte(data), 0666)
	}
}

// loadHighScoresData
func _(data []byte) error {
	highScores := make([]*HighScore, 0)
	lines := strings.Split(string(data), "\n")
	if len(lines) == 0 {
		return fmt.Errorf("no highscore data in file")
	}
	for _, l := range lines {
		parts := strings.Split(l, highScoreSeparator)
		if len(parts) != 2 {
			return fmt.Errorf("corrupted highscore data")
		}
		if i, err := strconv.Atoi(parts[1]); err == nil {
			if i < 0 {
				return fmt.Errorf("negative highscore - data corrupted")
			}
			highScores = append(highScores, &HighScore{i, parts[0]})
		} else {
			return err
		}
	}

	sort.Sort(sort.Reverse(ByScore(highScores)))
	if len(highScores) > 5 {
		highScores = highScores[:5]
	}
	return nil
}
