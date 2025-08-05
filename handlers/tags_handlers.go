package handlers

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/aukawut/BackendCompoundHt/model"
	"github.com/gofiber/fiber/v2"
	"github.com/skip2/go-qrcode"
)

func GetTagsByDate(c *fiber.Ctx) error {
	start := c.Query("start")
	end := c.Query("end")
	active := c.Query("active")

	db, err := OpenConnectDatabase()
	if err != nil {
		return c.JSON(fiber.Map{"err": true, "msg": err.Error()})
	}

	tags := []model.TagsList{}

	rows, errQuery := db.Query(fmt.Sprintf(`SELECT [Id]
      ,[OLD_TAG]
      ,[TAG_NO]
      ,[PART_NO]
      ,[QTY]
      ,[BATCH]
      ,[BASKET]
      ,[PO_NO]
	  ,[LOT_NO]
      ,[LOCATION]
      ,[ACTIVE]
      ,[CREATED_AT]
      ,[CREATED_BY]
      ,[UPDATED_AT]
      ,[UPDATED_BY]
  FROM [DB_COMPOUND].[dbo].[TBL_TAGS] WHERE CONVERT(DATE,CREATED_AT) BETWEEN '%s' AND '%s' AND [ACTIVE] = @active`, start, end), sql.Named("active", active))
	if errQuery != nil {
		return c.JSON(fiber.Map{"err": true, "msg": errQuery.Error()})
	}

	for rows.Next() {
		tag := model.TagsList{}

		err := rows.Scan(
			&tag.Id,
			&tag.OLD_TAG,
			&tag.TAG_NO,
			&tag.PART_NO,
			&tag.QTY,
			&tag.BATCH,
			&tag.BASKET,
			&tag.PO_NO,
			&tag.LOT_NO,
			&tag.LOCATION,
			&tag.ACTIVE,
			&tag.CREATED_AT,
			&tag.CREATED_BY,
			&tag.UPDATED_AT,
			&tag.UPDATED_BY,
		)

		if err != nil {
			return c.JSON(fiber.Map{"err": true, "msg": err.Error()})
		} else {
			tags = append(tags, tag)
		}
	}
	defer rows.Close()

	if len(tags) > 0 {
		return c.JSON(fiber.Map{"err": false, "msg": "", "status": "Ok", "results": tags})
	} else {
		return c.JSON(fiber.Map{"err": true, "msg": "Not Found"})
	}
}

func CheckDuplicatedTag(c *fiber.Ctx) error {

	// Open database connnection
	db, err := OpenConnectDatabase()

	count := 0

	tag := c.Query("tagNo")

	if err != nil {
		return c.JSON(fiber.Map{"err": true, "msg": err.Error()})
	}

	errRow := db.QueryRow(`SELECT COUNT(*) as [TAG_AMOUNT] FROM [dbo].[TBL_TAGS] WHERE [OLD_TAG] = @tag AND [ACTIVE] = 'Y'`, sql.Named("tag", tag)).Scan(&count)

	if count > 0 {
		return c.JSON(fiber.Map{"err": true, "msg": "Tag is duplicated!"})
	}
	defer db.Close()

	return c.JSON(fiber.Map{"err": false, "msg": "Tag is ok", "status": "Ok", "errRow": errRow})
}

func GetRunningPo() (string, error) {
	// Open Database connection
	db, err := OpenConnectDatabase()
	rows := 0

	po := ""
	prefix := `PO99/`

	if err != nil {
		return "", err
	}
	row := db.QueryRow(`SELECT COUNT(*) as [AMOUNT] FROM [dbo].[TBL_TAGS]`).Scan(&rows)
	if rows == 0 {
		fmt.Println(row)
		return prefix + "0000001", nil
	}
	rowPo := db.QueryRow(`SELECT CONCAT('PO99/',FORMAT(ISNULL(MAX(CAST(RIGHT([PO_NO],7) as INT)) + 1,1),'0000000')) as [PO] FROM [dbo].[TBL_TAGS]`).Scan(&po)
	fmt.Println(rowPo)
	return po, nil
}

func SaveTags(c *fiber.Ctx) error {

	var req model.SaveTagsBody
	errBody := c.BodyParser(&req)

	if errBody != nil {
		return c.JSON(fiber.Map{"err": true, "msg": errBody.Error()})
	}

	// Get PO Number
	poNo, _ := GetRunningPo()

	if poNo == "" {
		return c.JSON(fiber.Map{"err": true, "msg": "Generate PO Number is fail!"})
	}

	// Open Database connection
	db, err := OpenConnectDatabase()

	if err != nil {
		return c.JSON(fiber.Map{"err": true, "msg": err.Error()})
	}

	if len(req.Tags) > 0 {

		for _, tag := range req.Tags {

			newTag := fmt.Sprintf(`%s|%s`, poNo, tag.TagNo)

			_, errInsert := db.Exec(`INSERT INTO [dbo].[TBL_TAGS] ([LOT_NO],[OLD_TAG],[TAG_NO],[PART_NO],[QTY],[BATCH],[BASKET],[PO_NO],[ACTIVE],[CREATED_AT],[CREATED_BY],[LOCATION]) 
			VALUES (@lotNo,@oldTag,@newTagNo,@partNo,@qty,@batch,@basket,@po,'Y',GETDATE(),@actionBy,@location)`,
				sql.Named("lotNo", tag.Lot),
				sql.Named("oldTag", tag.TagNo),
				sql.Named("newTagNo", newTag),
				sql.Named("partNo", tag.PartNo),
				sql.Named("qty", tag.Qty),
				sql.Named("batch", tag.Batch),
				sql.Named("basket", tag.Basket),
				sql.Named("po", poNo),
				sql.Named("actionBy", req.ActionBy),
				sql.Named("location", req.Location),
			)

			if errInsert != nil {
				return c.JSON(fiber.Map{"err": true, "msg": errInsert.Error()})
			}

		}

	} else {
		return c.JSON(fiber.Map{"err": true, "msg": "Tags is required!"})
	}

	defer db.Close()

	return c.JSON(fiber.Map{"err": false, "msg": "Tag saved!", "status": "Ok"})
}

func GenerateQRCode(c *fiber.Ctx) error {
	tag := c.Query("tagNo")
	size := c.Query("size")

	sizeNum, _ := strconv.Atoi(size)

	if tag == "" {
		return c.JSON(fiber.Map{"err": true, "msg": "Tag No is required!"})
	}

	qr, err := qrcode.Encode(tag, qrcode.Medium, sizeNum)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate QR code",
		})
	}
	c.Set("Content-Type", "image/png")

	return c.Send(qr)

}

func GetCompoundTagDetail(c *fiber.Ctx) error {
	tagNo := c.Query("tagNo")

	db, err := OpenConnectDatabase()
	if err != nil {
		return c.JSON(fiber.Map{"err": true, "msg": err.Error()})
	}
	stmt := `SELECT [Id]
      ,[OLD_TAG]
      ,[TAG_NO]
      ,[PART_NO]
      ,[QTY]
      ,[BATCH]
      ,[BASKET]
      ,[PO_NO]
	  ,[LOT_NO]
      ,[LOCATION]
      ,[ACTIVE]
      ,[CREATED_AT]
      ,[CREATED_BY]
      ,[UPDATED_AT]
      ,[UPDATED_BY]
  FROM [DB_COMPOUND].[dbo].[TBL_TAGS]`

	tags := []model.TagsList{}

	if strings.Contains(tagNo, "|") {
		stmt += ` WHERE [TAG_NO] = @tagNo`
	} else {
		stmt += ` WHERE [OLD_TAG] = @tagNo`
	}

	rows, errQuery := db.Query(stmt, sql.Named("tagNo", tagNo))
	if errQuery != nil {
		return c.JSON(fiber.Map{"err": true, "msg": errQuery.Error()})
	}

	for rows.Next() {
		tag := model.TagsList{}

		err := rows.Scan(
			&tag.Id,
			&tag.OLD_TAG,
			&tag.TAG_NO,
			&tag.PART_NO,
			&tag.QTY,
			&tag.BATCH,
			&tag.BASKET,
			&tag.PO_NO,
			&tag.LOT_NO,
			&tag.LOCATION,
			&tag.ACTIVE,
			&tag.CREATED_AT,
			&tag.CREATED_BY,
			&tag.UPDATED_AT,
			&tag.UPDATED_BY,
		)

		if err != nil {
			return c.JSON(fiber.Map{"err": true, "msg": err.Error()})
		} else {
			tags = append(tags, tag)
		}
	}
	defer rows.Close()

	if len(tags) > 0 {
		return c.JSON(fiber.Map{"err": false, "msg": "", "status": "Ok", "results": tags})
	} else {
		return c.JSON(fiber.Map{"err": true, "msg": "Not Found"})
	}
}

func CancelTags(c *fiber.Ctx) error {

	var req model.TagsCancelBody

	errBody := c.BodyParser(&req)

	if errBody != nil {
		return c.JSON(fiber.Map{"err": true, "msg": errBody.Error()})
	}

	db, err := OpenConnectDatabase()
	if err != nil {
		return c.JSON(fiber.Map{"err": true, "msg": err.Error()})
	}

	if len(req.Ids) > 0 {
		// Loop update status of tags.
		for _, item := range req.Ids {
			_, err := db.Exec(`UPDATE [dbo].[TBL_TAGS] SET ACTIVE = 'N',[UPDATED_BY] = @actionBy,UPDATED_AT = GETDATE() WHERE [Id] = @id`, sql.Named("id", item.Id), sql.Named("actionBy", req.ActionBy))
			if err != nil {
				return c.JSON(fiber.Map{"err": true, "msg": err.Error()})
			}
		}
	}

	defer db.Close()

	return c.JSON(fiber.Map{"err": false, "msg": "Success", "status": "Ok"})
}
