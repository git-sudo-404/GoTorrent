/*
 * MIT License
 *
 * Copyright (c) 2026 git-sudo-404
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * Of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * Copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * Copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING, BUT NOT LIMITED TO, THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES, OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package torrent

import "fmt"

func checkFieldPresenceAndType[T any](dict map[string]any, field string) (T, error) {
	var zero T
	value, ok := dict[field]
	if !ok {
		return zero, fmt.Errorf("Meta-Info missing needed field : %s", field)
	}
	if v, ok := value.(T); !ok {
		return zero, fmt.Errorf("Undesired field value type  , Field : %s ", field)
	} else {
		return v, nil
	}
}

func checkOptionalFieldPresenceAndType[T any](dict map[string]any, field string) (T, bool, error) {
	var zero T
	value, ok := dict[field]
	if !ok {
		return zero, false, nil
	}
	if v, ok := value.(T); !ok {
		return zero, false, fmt.Errorf("Undesired field value type , Field : %s ", field)
	} else {
		return v, true, nil
	}
}
