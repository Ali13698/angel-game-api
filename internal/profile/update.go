package profile

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/angel/angel/internal/db"
)

func UpdateMe(ctx context.Context, uid int64, r *updateReq) error {
	p := db.Pool()
	if p == nil { return errStr("db not init") }

	var sets []string
	var args []any
	add := func(col string, v *string) {
		if v == nil { return }
		args = append(args, *v)
		sets = append(sets, col+" = $"+strconv.Itoa(len(args)))
	}
	add("display_name", r.DisplayName)
	add("bio", r.Bio)
	add("avatar", r.Avatar)
	add("name_emoji", r.NameEmoji)
	add("profile_color", r.ProfileColor)
	add("banner_id", r.BannerID)
	add("frame_id", r.FrameID)
	add("hat_id", r.HatID)
	add("chat_bubble_id", r.ChatBubbleID)
	add("badge_id", r.BadgeID)

	if len(sets) == 0 { return nil }

	args = append(args, time.Now().Unix())
	sets = append(sets, "last_seen = $"+strconv.Itoa(len(args)))
	args = append(args, uid)

	q := "UPDATE users SET " + strings.Join(sets, ", ") + " WHERE id = $" + strconv.Itoa(len(args))
	_, err := p.Exec(ctx, q, args...)
	if err != nil { return fmt.Errorf("update: %w", err) }
	return nil
}
