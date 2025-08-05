package handlers

import (
	"fmt"
	"strings"

	"github.com/aukawut/BackendCompoundHt/model"
	"github.com/gofiber/fiber/v2"
)

func FilterPartByLocation(location string, parts []model.Parts) model.Parts {

	for _, part := range parts {
		factory := fmt.Sprintf(`(%s)`, location)
		if strings.Contains(part.CODE, factory) {
			return part
		}
	}
	return parts[0]
}

func GetPartCodeByCompoundTags(c *fiber.Ctx) error {
	db, err := OpenConnectDatabase()
	if err != nil {
		return c.JSON(fiber.Map{"err": true, "msg": err.Error()})
	}
	defer db.Close() // ✅ Close DB connection

	tagNo := c.Query("tagNo")
	location := c.Query("location")

	var parts []model.Parts

	stmt := fmt.Sprintf(`SELECT [CODE],[NAME],ISNULL([RM_GROUP],'Null') as [RM_GROUP] 
		FROM [dbo].[fn_GetPartNoByCompoundTag]('%s');`, tagNo)

	rows, errQuery := db.Query(stmt)
	if errQuery != nil {
		return c.JSON(fiber.Map{"err": true, "msg": errQuery.Error()})
	}
	defer rows.Close()

	for rows.Next() {
		var part model.Parts
		if err := rows.Scan(&part.CODE, &part.NAME, &part.RM_GROUP); err != nil {
			return c.JSON(fiber.Map{"err": true, "msg": err.Error()})
		}
		parts = append(parts, part)
	}

	if len(parts) == 0 {
		return c.JSON(fiber.Map{"err": true, "msg": "Part is not found"})
	}

	// If multiple parts found, filter by location
	if len(parts) > 1 {
		part := FilterPartByLocation(location, parts)
		partFilter := []model.Parts{part} // ✅ Correct way to wrap in slice
		return c.JSON(fiber.Map{"err": false, "msg": "", "status": "Ok", "results": partFilter})
	}

	// Only one part found
	return c.JSON(fiber.Map{"err": false, "msg": "", "status": "Ok", "results": parts})
}
