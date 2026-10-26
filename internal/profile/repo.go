package profile

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"errors"
	"time"

	"github.com/angel/angel/internal/db"
	"github.com/jackc/pgx/v5"
)

// ... (ساختار User بدون تغییر می‌ماند، همان ۳۶ فیلد)

func generateToken() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b) // یک رشته ۱۶ کاراکتری
}

func GetByID(ctx context.Context, id int64) (*User, error) {
	p := db.Pool()
	if p == nil { return nil, errors.New("db not init") }
	
	// ⬇️ ستون public_profile_token به کوئری اضافه شد ⬇️
	row := p.QueryRow(ctx, `SELECT id, COALESCE(username,''), COALESCE(first_name,''), COALESCE(game_username,''), COALESCE(display_name,''), COALESCE(avatar,''), COALESCE(photo_url,''), COALESCE(profile_color,'#8b5cf6'), COALESCE(profile_banner,'default'), COALESCE(name_effect,''), COALESCE(name_emoji,''), COALESCE(bio,''), COALESCE(banner_id,''), COALESCE(frame_id,''), COALESCE(hat_id,''), COALESCE(chat_bubble_id,''), COALESCE(badge_id,''), COALESCE(coins,0), COALESCE(gems,0), COALESCE(wins,0), COALESCE(losses,0), COALESCE(level,1), COALESCE(xp,0), COALESCE(likes_count,0), COALESCE(streak_days,0), COALESCE(active_days,0), COALESCE(referral_code,''), COALESCE(referrals_count,0), COALESCE(created_at,0), COALESCE(last_seen,0), COALESCE(public_profile_token,'') FROM users WHERE id = $1`, id)
	
	var u User
	err := row.Scan(&u.ID, &u.Username, &u.FirstName, &u.GameUsername, &u.DisplayName, &u.Avatar, &u.PhotoURL, &u.ProfileColor, &u.ProfileBanner, &u.NameEffect, &u.NameEmoji, &u.Bio, &u.BannerID, &u.FrameID, &u.HatID, &u.ChatBubbleID, &u.BadgeID, &u.Coins, &u.Gems, &u.Wins, &u.Losses, &u.Level, &u.XP, &u.LikesCount, &u.StreakDays, &u.ActiveDays, &u.ReferralCode, &u.ReferralsCount, &u.CreatedAt, &u.LastSeen, &u.PublicProfileToken)
	
	if err != nil { 
		if errors.Is(err, pgx.ErrNoRows) {
			// ساخت کاربر مهمان با توکن اختصاصی
			guestName := fmt.Sprintf("Guest_%d", id%10000)
			token := generateToken()
			_, insErr := p.Exec(ctx, `INSERT INTO users (id, game_username, display_name, created_at, last_seen, public_profile_token) VALUES ($1, $2, $3, $4, $5, $6)`, id, guestName, guestName, time.Now().Unix(), time.Now().Unix(), token)
			if insErr != nil {
				return nil, fmt.Errorf("failed to create guest: %w", insErr)
			}
			return GetByID(ctx, id)
		}
		return nil, err 
	}
	if u.CreatedAt > 0 { u.AgeDays = int((time.Now().Unix() - u.CreatedAt) / 86400) }
	return &u, nil
}

// ⬇️ تابع جدید: جستجو با لینک اختصاصی (Token) ⬇️
func GetByPublicToken(ctx context.Context, token string) (*User, error) {
	p := db.Pool()
	if p == nil { return nil, errors.New("db not init") }
	
	row := p.QueryRow(ctx, `SELECT id, COALESCE(username,''), COALESCE(first_name,''), COALESCE(game_username,''), COALESCE(display_name,''), COALESCE(avatar,''), COALESCE(photo_url,''), COALESCE(profile_color,'#8b5cf6'), COALESCE(profile_banner,'default'), COALESCE(name_effect,''), COALESCE(name_emoji,''), COALESCE(bio,''), COALESCE(banner_id,''), COALESCE(frame_id,''), COALESCE(hat_id,''), COALESCE(chat_bubble_id,''), COALESCE(badge_id,''), COALESCE(coins,0), COALESCE(gems,0), COALESCE(wins,0), COALESCE(losses,0), COALESCE(level,1), COALESCE(xp,0), COALESCE(likes_count,0), COALESCE(streak_days,0), COALESCE(active_days,0), COALESCE(referral_code,''), COALESCE(referrals_count,0), COALESCE(created_at,0), COALESCE(last_seen,0), COALESCE(public_profile_token,'') FROM users WHERE public_profile_token = $1 LIMIT 1`, token)
	
	var u User
	err := row.Scan(&u.ID, &u.Username, &u.FirstName, &u.GameUsername, &u.DisplayName, &u.Avatar, &u.PhotoURL, &u.ProfileColor, &u.ProfileBanner, &u.NameEffect, &u.NameEmoji, &u.Bio, &u.BannerID, &u.FrameID, &u.HatID, &u.ChatBubbleID, &u.BadgeID, &u.Coins, &u.Gems, &u.Wins, &u.Losses, &u.Level, &u.XP, &u.LikesCount, &u.StreakDays, &u.ActiveDays, &u.ReferralCode, &u.ReferralsCount, &u.CreatedAt, &u.LastSeen, &u.PublicProfileToken)
	
	if err != nil { return nil, err }
	if u.CreatedAt > 0 { u.AgeDays = int((time.Now().Unix() - u.CreatedAt) / 86400) }
	return &u, nil
}

// توابع GetByGameUsername و UpdateMe (بدون تغییر)
