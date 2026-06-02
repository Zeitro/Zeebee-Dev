'use client';

import React from 'react';
import styles from './index.module.scss';

export const ROW_LEVEL = {
	PARENT: 'parent',
	FILE: 'file',
};

/** @returns {'parent' | 'file' | 'both'} */
export function getRowPlacement(column) {
	return column?.rowPlacement ?? ROW_LEVEL.PARENT;
}

export function appliesToLevel(column, level) {
	const placement = getRowPlacement(column);
	if (placement === 'both') return true;
	return placement === level;
}

export function isPlaceholderCell(column, level) {
	return !appliesToLevel(column, level);
}

export function getDataIndex(column, level) {
	if (level === ROW_LEVEL.FILE) {
		return column.fileDataIndex ?? column.dataIndex;
	}
	return column.parentDataIndex ?? column.dataIndex;
}

export function getRecordValue(record, column, level) {
	if (isPlaceholderCell(column, level)) return undefined;
	const dataIndex = getDataIndex(column, level);
	return dataIndex ? record?.[dataIndex] : undefined;
}

export function getColumnFlexStyle(column, isParent) {
	const fontSize = isParent ? '16px' : '14px';
	if (column?.minWidth) {
		return { minWidth: column.minWidth, flex: '1 1 0', fontSize };
	}
	if (column?.width) {
		return { width: column.width, flex: `0 0 ${column.width}px`, fontSize };
	}

	return { flex: '0 0 auto', fontSize };
}

export function getColumnCellClassName(column, level, columnIndex, totalColumns) {
	return [
		column?.ellipsis && 'min-w-0 truncate',
		columnIndex === totalColumns - 1 && 'flex items-center',
		column?.align === 'right' && 'justify-end',
		column?.align === 'center' && 'justify-center',
		level === ROW_LEVEL.PARENT && columnIndex === 0 && 'font-medium text-foreground',
		level === ROW_LEVEL.FILE && columnIndex === 0 && 'min-w-0',
	]
		.filter(Boolean)
		.join(' ');
}

/** 仅占位对齐，不展示 `--` 或其它文案 */
export function renderPlaceholder(level) {
	return (
		<span
			className={[styles.cellSpacer, level === ROW_LEVEL.FILE && styles.cellSpacerFile]
				.filter(Boolean)
				.join(' ')}
			aria-hidden
		/>
	);
}
