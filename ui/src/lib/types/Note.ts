export type Note = {
  id: number;
  user_id: number;
  title: string;
  content: string;
  filepath: string;
  is_favorite: boolean;
  created_at: string;
  updated_at: string;
}