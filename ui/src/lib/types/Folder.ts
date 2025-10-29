import type { Note } from './Note';
export type Folder =
	{
		id: number;
		title: string;
		notes: Note[];
		parent_folder_id: {
			Int64: number;
			Valid: boolean;
		}
		subfolders?: Folder[];
	}