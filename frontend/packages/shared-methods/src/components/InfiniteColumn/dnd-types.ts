export interface DragInfo {
	source: { columnId: string; index: number };
	destination: { columnId: string; index: number } | null;
	draggableId: string;
	type?: string;
}
