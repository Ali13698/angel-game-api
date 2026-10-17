export interface User {
  id: number; display_name: string; game_username: string; bio: string;
  avatar: string; name_emoji: string; profile_color: string;
  coins: number; gems: number; level: number;
}
export interface ProfileResponse { ok: boolean; user?: User; error?: string; }
