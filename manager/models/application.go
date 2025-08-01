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

package models

import "gorm.io/gorm"

type Application struct {
	BaseModel
	Name     string  `gorm:"column:name;type:varchar(256);index:uk_application_name,unique;not null;comment:name" json:"name"`
	URL      string  `gorm:"column:url;not null;comment:url" json:"url"`
	BIO      string  `gorm:"column:bio;type:varchar(1024);comment:biography" json:"bio"`
	Priority JSONMap `gorm:"column:priority;not null;comment:download priority" json:"priority"`
	UserID   uint    `gorm:"comment:user id" json:"user_id"`
	User     User    `gorm:"-" json:"user"`
}

func (a *Application) AfterFind(tx *gorm.DB) (err error) {
	if a.UserID == 0 || a.UserID == a.User.ID {
		return nil
	}

	var user User
	if err := tx.First(&user, a.UserID).Error; err != nil {
		return err
	}

	a.User = user
	return nil
}
