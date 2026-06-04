'use client';

import React from 'react';
import { FieldTips } from '@repo/shared-methods';
import { EMPTY_DISPLAY } from '@repo/shared-methods/utils';

const DEFAULT_MAX_VISIBLE = 4;
const DEFAULT_COLS = 2;
const DEFAULT_TAG_CLASS =
	'rounded-[2px] px-2 text-xs border border-base-border truncate min-w-0 w-max';

/**
 * 表格单元格标签网格：每行固定列数，最多展示若干行；超出部分 hover 展示完整列表。
 * 纯展示，不修改入参 items。
 *
 * @param {string[]} [items] 标签文案列表
 * @param {number} [maxVisible=4] 最多展示标签数（默认 2 列 × 2 行）
 * @param {number} [cols=2] 每行标签列数
 * @param {string} [tagClassName] 单个标签基础样式
 * @param {{ className?: string }} [tagProps] 追加到每个标签的 props（如列配色）
 */
const TruncatedTagGrid = ({
	items = [],
	maxVisible = DEFAULT_MAX_VISIBLE,
	cols = DEFAULT_COLS,
	tagClassName = DEFAULT_TAG_CLASS,
	tagProps = {},
}) => {
	if (!items?.length) {
		return EMPTY_DISPLAY;
	}

	const visibleItems = items.slice(0, maxVisible);
	const hasOverflow = items.length > maxVisible;

	const grid = (
		<div
			className="grid gap-1"
			style={{ gridTemplateColumns: `repeat(${cols}, minmax(0, 1fr))` }}
		>
			{visibleItems.map((item, index) => (
				<div
					key={`${item}-${index}`}
					className={[tagClassName, tagProps?.className].filter(Boolean).join(' ')}
				>
					{item}
				</div>
			))}
		</div>
	);

	if (!hasOverflow) {
		return grid;
	}

	const tooltipTitle = (
		<div className="flex flex-wrap gap-1" onClick={(e) => e.stopPropagation()}>
			{items.map((item, index) => (
				<div
					key={`${item}-${index}`}
					className={[tagClassName, tagProps?.className].filter(Boolean).join(' ')}
				>
					{item}
				</div>
			))}
		</div>
	);

	return <FieldTips title={tooltipTitle}>{grid}</FieldTips>;
};

export default TruncatedTagGrid;
