// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package migrations

import (
	"fmt"
	"os"
	"path"

	"code.gitea.io/gitea/modules/setting"

	"xorm.io/builder"
	"xorm.io/xorm"
)

func removeAttachmentMissedRepo(x *xorm.Engine) error {
	type Attachment struct {
		UUID string `xorm:"uuid"`
	}
	var start int
	attachments := make([]*Attachment, 0, 50)
	for {
		err := x.Select("uuid").Where(builder.NotIn("release_id", builder.Select("id").From("`release`"))).
			And("release_id > 0").
			OrderBy("id").Limit(50, start).Find(&attachments)
		if err != nil {
			return err
		}

		for i := 0; i < len(attachments); i++ {
			uuid := attachments[i].UUID
			if err = os.RemoveAll(path.Join(setting.AttachmentPath, uuid[0:1], uuid[1:2], uuid)); err != nil {
				fmt.Printf("Error: %v", err)
			}
		}

		if len(attachments) < 50 {
			break
		}
		start += 50
		attachments = attachments[:0]
	}

	_, err := x.Exec("DELETE FROM attachment WHERE release_id > 0 AND release_id NOT IN (SELECT id FROM `release`)")
	return err
}
