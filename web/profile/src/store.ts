export interface StoreItem { id: string; type: 'banner' | 'avatar' | 'frame' | 'hat'; url: string; price: number; owned: boolean; }
export const DEFAULT_BANNERS: StoreItem[] = [
  { id: 'b1', type: 'banner', url: 'linear-gradient(135deg, #7c3aed 0%, #ec4899 100%)', price: 0, owned: true },
  { id: 'b2', type: 'banner', url: 'linear-gradient(135deg, #06b6d4 0%, #3b82f6 100%)', price: 0, owned: true },
  { id: 'b3', type: 'banner', url: 'linear-gradient(135deg, #f59e0b 0%, #ef4444 100%)', price: 0, owned: true },
  { id: 'b4', type: 'banner', url: 'linear-gradient(135deg, #10b981 0%, #3b82f6 100%)', price: 100, owned: false },
  { id: 'b5', type: 'banner', url: "url('https://images.unsplash.com/photo-1579546929518-9e396f3cc809?w=600&q=80')", price: 200, owned: false },
  { id: 'b6', type: 'banner', url: "url('https://images.unsplash.com/photo-1557682250-33bd709cbe85?w=600&q=80')", price: 200, owned: false }
];
export const DEFAULT_AVATARS: StoreItem[] = [
  { id: 'a1', type: 'avatar', url: 'https://api.dicebear.com/7.x/adventurer/svg?seed=1', price: 0, owned: true },
  { id: 'a2', type: 'avatar', url: 'https://api.dicebear.com/7.x/avataaars/svg?seed=2', price: 0, owned: true },
  { id: 'a3', type: 'avatar', url: 'https://api.dicebear.com/7.x/bottts/svg?seed=3', price: 0, owned: true },
  { id: 'a4', type: 'avatar', url: 'https://api.dicebear.com/7.x/fun-emoji/svg?seed=4', price: 0, owned: true },
  { id: 'a5', type: 'avatar', url: 'https://api.dicebear.com/7.x/lorelei/svg?seed=5', price: 0, owned: true }
];
