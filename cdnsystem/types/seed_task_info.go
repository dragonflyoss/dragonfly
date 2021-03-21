/*
 *     Copyright 2020 The Dragonfly Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package types

//type UrlMeta struct {
//	// md5 of file content downloaded from url
//	Md5 string `json:"md5,omitempty"`
//	// downloading range of resource file
//	Range string `json:"range,omitempty"`
//}

type SeedTask struct {
	TaskId           string            `json:"taskId,omitempty"`
	Url              string            `json:"url,omitempty"`
	TaskUrl          string            `json:"taskUrl,omitempty"`
	SourceFileLength int64             `json:"sourceFileLength,omitempty"`
	CdnFileLength    int64             `json:"cdnFileLength,omitempty"`
	PieceSize        int32             `json:"pieceSize,omitempty"`
	Header           map[string]string `json:"header,omitempty"`
	CdnStatus        string            `json:"cdnStatus,omitempty"`
	PieceTotal       int32             `json:"pieceTotal,omitempty"`
	RequestMd5       string            `json:"requestMd5,omitempty"`
	SourceRealMd5    string            `json:"sourceRealMd5,omitempty"`
	PieceMd5Sign     string            `json:"pieceMd5Sign,omitempty"`
}

const (

	// TaskInfoCdnStatusWAITING captures enum value "WAITING"
	TaskInfoCdnStatusWAITING string = "WAITING"

	// TaskInfoCdnStatusRUNNING captures enum value "RUNNING"
	TaskInfoCdnStatusRUNNING string = "RUNNING"

	// TaskInfoCdnStatusFAILED captures enum value "FAILED"
	TaskInfoCdnStatusFAILED string = "FAILED"

	// TaskInfoCdnStatusSuccess captures enum value "SUCCESS"
	TaskInfoCdnStatusSuccess string = "SUCCESS"

	// TaskInfoCdnStatusSourceERROR captures enum value "SOURCE_ERROR"
	TaskInfoCdnStatusSourceERROR string = "SOURCE_ERROR"
)
