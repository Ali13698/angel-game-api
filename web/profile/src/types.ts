export interface User {
    id: number;
    game_username: string;
    display_name: string;
    bio: string;
    avatar: string;
    banner_id?: string;
    profile_color?: string;
    likes_count: number;
    coins: number;
    gems: number;
    wins: number;
    age_days: number;
}
