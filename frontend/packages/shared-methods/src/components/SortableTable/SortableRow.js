'use client';

import React, { createContext, useContext, useMemo } from 'react';
import { HolderOutlined } from '@ant-design/icons';
import { useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { Button } from 'antd';
import styles from './index.module.scss';

const RowContext = createContext({});

/**
 * 行拖拽手柄，须放在 SortableRow 渲染的单元格内。
 * 仅手柄可拖，避免与行内按钮 / 行点击冲突。
 */
export function DragHandle() {
	const { setActivatorNodeRef, listeners } = useContext(RowContext);
	return (
		<Button
			type="text"
			size="small"
			icon={<HolderOutlined />}
			style={{ cursor: 'move', touchAction: 'none' }}
			ref={setActivatorNodeRef}
			onClick={(e) => e.stopPropagation()}
			{...listeners}
		/>
	);
}

/**
 * antd Table `components.body.row`：把 tr 注册为 sortable 项。
 * 依赖 Table 写入的 `data-row-key`（与 rowKey 一致）。
 * 位移走 CSS 变量到 td：浏览器不会对 table-row 做 transform。
 */
export function SortableRow(props) {
	/** antd 传入的 data-row-key 可能是 number，必须和 SortableContext items 一样转成 string */
	const rowId = props['data-row-key'] == null ? '' : String(props['data-row-key']);
	const {
		attributes,
		setNodeRef,
		setActivatorNodeRef,
		transform,
		transition,
		isDragging,
		listeners,
	} = useSortable({
		id: rowId,
		disabled: rowId === '',
	});

	const style = {
		...props.style,
		'--row-transform': CSS.Translate.toString(transform) || 'none',
		'--row-transition': transition,
	};

	const contextValue = useMemo(
		() => ({ setActivatorNodeRef, listeners }),
		[setActivatorNodeRef, listeners],
	);

	return (
		<RowContext.Provider value={contextValue}>
			<tr
				{...props}
				ref={setNodeRef}
				style={style}
				className={[props.className, styles.sortableRow, isDragging ? styles.dragging : '']
					.filter(Boolean)
					.join(' ')}
				{...attributes}
			/>
		</RowContext.Provider>
	);
}
