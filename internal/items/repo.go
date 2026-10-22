package items

import (
	"context"
	"encoding/json"

	"github.com/angel/angel/internal/db"
)

// GetUserActiveItems واکشی آیتم‌های مجهز شده کاربر از ستون JSONB
func GetUserActiveItems(ctx context.Context, userID int64) ([]InfiniteItem, error) {
	p := db.Pool()
	if p == nil {
		return nil, nil
	}

	rows, err := p.Query(ctx, `
		SELECT i.item_id, i.item_type, s.meta
		FROM inventory i
		JOIN shop_items s ON s.item_id = i.item_id
		WHERE i.user_id = $1 AND i.equipped = 1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []InfiniteItem
	for rows.Next() {
		var id, itype string
		var metaRaw []byte
		
		if err := rows.Scan(&id, &itype, &metaRaw); err != nil {
			continue
		}

		var meta map[string]any
		if len(metaRaw) > 0 {
			_ = json.Unmarshal(metaRaw, &meta)
		}

		items = append(items, InfiniteItem{
			ID:   id,
			Type: ItemType(itype),
			Meta: meta,
		})
	}
	return items, nil
}
