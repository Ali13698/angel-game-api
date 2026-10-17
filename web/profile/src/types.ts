export interface User {
  id: number; display_name: string; game_username: string; bio: string;
  avatar: string; name_emoji: string; profile_color: string;
  coins: number; gems: number; level: number; wins: number; likes_count: number; age_days: number;
}

export interface GameStat {
  id: string; name: string; icon: string;
  level: number; rank: string;
  played: number; won: number;
  progress: number; // 0-100
}

export interface ProfileResponse { ok: boolean; user?: User; error?: string; }
