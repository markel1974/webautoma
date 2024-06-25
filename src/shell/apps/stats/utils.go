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

package stats

import (
	"io/ioutil"
	"strconv"
	"strings"
)

func bToMb(b uint64) float64 {
	return float64(b) / 1024 / 1024
}

func getCPUSample() (uint64, uint64) {
	var idle uint64 = 0  //(rand.Intn(max - min) + min)
	var total uint64 = 0 //(rand.Intn(max - min) + min)

	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return idle, total
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err == nil {
					total += val // tally up all the numbers to get total ticks
					if i == 4 {  // idle is the 5th field in the cpu line
						idle = val
					}
				}
			}
			return idle, total
		}
	}
	return total, total
}
