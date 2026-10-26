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

type User struct {
	ID                 int64  `json:"id"`
	Username           string `json:"username,omitempty"`
	FirstName          string `json:"first_name,omitempty"`
	GameUsername       string `json:"game_username,omitempty"`
	DisplayName        string `json:"display_name,omitempty"`
	Avatar             string `json:"avatar,omitempty"`
	PhotoURL           string `json:"photo_url,omitempty"`
	ProfileColor       string `json:"profile_color,omitempty"`
	ProfileBanner      string `json:"profile_banner,omitempty"`
	NameEffect         string `json:"name_effect,omitempty"`
	NameEmoji          string `json:"name_emoji,omitempty"`
	Bio                string `json:"bio,omitempty"`
	BannerID           string `json:"banner_id,omitempty"`
	FrameID            string `json:"frame_id,omitempty"`
	HatID              string `json:"hat_id,omitempty"`
	ChatBubbleID       string `json:"chat_bubble_id,omitempty"`
	BadgeID            string `json:"badge_id,omitempty"`
	Coins              int    `json:"coins"`
	Gems               int    `json:"gems"`
	Wins               int    `json:"wins"`
	Losses             int    `json:"losses"`
	Level              int    `json:"level"`
	XP                 int    `json:"xp"`
	LikesCount         int    `json:"likes_count"`
	StreakDays         int    `json:"streak_days"`
	ActiveDays         int    `json:"active_days"`
	ReferralCode       string `json:"referral_code,omitempty"`
	ReferralsCount     int    `json:"referrals_count"`
	CreatedAt          int64  `json:"created_at"`
	LastSeen           int64  `json:"last_seen"`
	AgeDays            int    `json:"age_days"`
	PublicProfileToken string `json:"public_profile_token,omitempty"`
}

func generateToken() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func GetByID(ctx context.Context, id int64) (*User, error) {
	p := db.Pool()
	if p == nil { return nil, errors.New("db not init") }
	row := p.QueryRow(ctx, `SELECT id, COALESCE(username,''), COALESCE(first_name,''), COALESCE(game_username,''), COALESCE(display_name,''), COALESCE(avatar,''), COALESCE(photo_url,''), COALESCE(profile_color,'#8b5cf6'), COALESCE(profile_banner,'default'), COALESCE(name_effect,''), COALESCE(name_emoji,''), COALESCE(bio,''), COALESCE(banner_id,''), COALESCE(frame_id,''), COALESCE(hat_id,''), COALESCE(chat_bubble_id,''), COALESCE(badge_id,''), COALESCE(coins,0), COALESCE(gems,0), COALESCE(wins,0), COALESCE(losses,0), COALESCE(level,1), COALESCE(xp,0), COALESCE(likes_count,0), COALESCE(streak_days,0), COALESCE(active_days,0), COALESCE(referral_code,''), COALESCE(referrals_count,0), COALESCE(created_at,0), COALESCE(last_seen,0), COALESCE(public_profile_token,'') FROM users WHERE id = $1`, id)
	var u User
	err := row.Scan(&u.ID, &u.Username, &u.FirstName, &u.GameUsername, &u.DisplayName, &u.Avatar, &u.PhotoURL, &u.ProfileColor, &u.ProfileBanner, &u.NameEffect, &u.NameEmoji, &u.Bio, &u.BannerID, &u.FrameID, &u.HatID, &u.ChatBubbleID, &u.BadgeID, &u.Coins, &u.Gems, &u.Wins, &u.Losses, &u.Level, &u.XP, &u.LikesCount, &u.StreakDays, &u.ActiveDays, &u.ReferralCode, &u.ReferralsCount, &u.CreatedAt, &u.LastSeen, &u.PublicProfileToken)
	if err != nil { 
		if errors.Is(err, pgx.ErrNoRows) {
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

func GetByGameUsername(ctx context.Context, username string) (*User, error) {
	p := db.Pool()
	if p == nil { return nil, errors.New("db not init") }
	row := p.QueryRow(ctx, `SELECT id, COALESCE(username,''), COALESCE(first_name,''), COALESCE(game_username,''), COALESCE(display_name,''), COALESCE(avatar,''), COALESCE(photo_url,''), COALESCE(profile_color,'#8b5cf6'), COALESCE(profile_banner,'default'), COALESCE(name_effect,''), COALESCE(name_emoji,''), COALESCE(bio,''), COALESCE(banner_id,''), COALESCE(frame_id,''), COALESCE(hat_id,''), COALESCE(chat_bubble_id,''), COALESCE(badge_id,''), COALESCE(coins,0), COALESCE(gems,0), COALESCE(wins,0), COALESCE(losses,0), COALESCE(level,1), COALESCE(xp,0), COALESCE(likes_count,0), COALESCE(streak_days,0), COALESCE(active_days,0), COALESCE(referral_code,''), COALESCE(referrals_count,0), COALESCE(created_at,0), COALESCE(last_seen,0), COALESCE(public_profile_token,'') FROM users WHERE game_username = $1 LIMIT 1`, username)
	var u User
	err := row.Scan(&u.ID, &u.Username, &u.FirstName, &u.GameUsername, &u.DisplayName, &u.Avatar, &u.PhotoURL, &u.ProfileColor, &u.ProfileBanner, &u.NameEffect, &u.NameEmoji, &u.Bio, &u.BannerID, &u.FrameID, &u.HatID, &u.ChatBubbleID, &u.BadgeID, &u.Coins, &u.Gems, &u.Wins, &u.Losses, &u.Level, &u.XP, &u.LikesCount, &u.StreakDays, &u.ActiveDays, &u.ReferralCode, &u.ReferralsCount, &u.CreatedAt, &u.LastSeen, &u.PublicProfileToken)
	if err != nil { return nil, err }
	if u.CreatedAt > 0 { u.AgeDays = int((time.Now().Unix() - u.CreatedAt) / 86400) }
	return &u, nil
}
