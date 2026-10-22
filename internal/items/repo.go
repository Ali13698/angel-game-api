package items

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/angel/angel/internal/db"
)

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

// EquipItem به صورت اتمی آیتم جدید را مجهز می‌کند (ضد تداخل)
func EquipItem(ctx context.Context, userID int64, itemID string) error {
	p := db.Pool()
	if p == nil {
		return errors.New("db not initialized")
	}

	tx, err := p.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var itemType string
	err = tx.QueryRow(ctx, "SELECT item_type FROM shop_items WHERE item_id = $1", itemID).Scan(&itemType)
	if err != nil {
		return errors.New("item not found")
	}

	var owned int
	err = tx.QueryRow(ctx, "SELECT 1 FROM inventory WHERE user_id = $1 AND item_id = $2", userID, itemID).Scan(&owned)
	if err != nil {
		return errors.New("you do not own this item")
	}

	_, err = tx.Exec(ctx, "UPDATE inventory SET equipped = 0 WHERE user_id = $1 AND item_type = $2", userID, itemType)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE inventory SET equipped = 1 WHERE user_id = $1 AND item_id = $2", userID, itemID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
