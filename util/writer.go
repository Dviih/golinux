/*
 *     Execute binaries on bare Linux.
 *     Copyright (C) 2025  Dviih
 *
 *     This program is free software: you can redistribute it and/or modify
 *     it under the terms of the GNU Affero General Public License as published
 *     by the Free Software Foundation, either version 3 of the License, or
 *     (at your option) any later version.
 *
 *     This program is distributed in the hope that it will be useful,
 *     but WITHOUT ANY WARRANTY; without even the implied warranty of
 *     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *     GNU Affero General Public License for more details.
 *
 *     You should have received a copy of the GNU Affero General Public License
 *     along with this program.  If not, see <https://www.gnu.org/licenses/>.
 *
 */

package util

import "errors"

type Writer struct {
	data []byte
}

func (writer *Writer) Write(data []byte) (int, error) {
	writer.data = append(writer.data, data...)
	return len(data), nil
}

func (writer *Writer) Len() int {
	return len(writer.data)
}

func (writer *Writer) Data() []byte {
	data := make([]byte, len(writer.data))

	copy(data, writer.data)
	return data
}

func (writer *Writer) Error(err error) error {
	if writer.data == nil {
		return err
	}

	return errors.Join(err, errors.New(string(writer.data)))
}
