package profile

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/angel/angel/internal/db"
)

func UpdateMe(ctx context.Context, uid int64, r *updateReq) error {
	log.Printf("[UPDATE] User %d requesting profile update", uid)
	p := db.Pool()
	if p == nil { 
		log.Printf("[UPDATE] DB pool is nil")
		return errStr("db not init") 
	}

	if r.GameUsername != nil {
		log.Printf("[UPDATE] Validating username: %s", *r.GameUsername)
		if len(*r.GameUsername) > 10 {
			log.Printf("[UPDATE] Validation failed: username too long")
			return errStr("نام کاربری حداکثر ۱۰ کاراکتر مجاز است")
		}
		match, _ := regexp.MatchString("^[a-zA-Z0-9]+$", *r.GameUsername)
		if !match {
			log.Printf("[UPDATE] Validation failed: username contains invalid chars")
			return errStr("نام کاربری فقط شامل حروف و اعداد باشد")
		}
	}

	if r.Bio != nil {
		log.Printf("[UPDATE] Validating bio length...")
		if len([]rune(*r.Bio)) > 100 {
			log.Printf("[UPDATE] Validation failed: bio too long")
			return errStr("بیوگرافی حداکثر ۱۰۰ کاراکتر مجاز است")
		}
	}

	var sets []string
	var args []any
	add := func(col string, v *string) {
		if v == nil { return }
		args = append(args, *v)
		sets = append(sets, col+" = $"+strconv.Itoa(len(args)))
	}
	add("display_name", r.DisplayName)
	add("game_username", r.GameUsername)
	add("bio", r.Bio)
	add("avatar", r.Avatar)
	add("name_emoji", r.NameEmoji)
	add("profile_color", r.ProfileColor)
	add("banner_id", r.BannerID)
	add("frame_id", r.FrameID)
	add("hat_id", r.HatID)
	add("chat_bubble_id", r.ChatBubbleID)
	add("badge_id", r.BadgeID)

	if len(sets) == 0 { 
		log.Printf("[UPDATE] No changes detected for user %d", uid)
		return nil 
	}

	args = append(args, time.Now().Unix())
	sets = append(sets, "last_seen = $"+strconv.Itoa(len(args)))
	args = append(args, uid)

	q := "UPDATE users SET " + strings.Join(sets, ", ") + " WHERE id = $" + strconv.Itoa(len(args))
	log.Printf("[UPDATE] Executing DB query for user %d", uid)
	_, err := p.Exec(ctx, q, args...)
	if err != nil { 
		log.Printf("[UPDATE] DB Exec failed: %v", err)
		return fmt.Errorf("update: %w", err) 
	}
	log.Printf("[UPDATE] DB Update successful for user %d", uid)
	return nil
}