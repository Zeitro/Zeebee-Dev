'use client';

import React from 'react';
import { FieldTips } from '@repo/shared-methods';
import { EMPTY_DISPLAY } from '@repo/shared-methods/utils';

const DEFAULT_MAX_VISIBLE = 4;
const DEFAULT_TAG_CLASS =
	'rounded-[2px] px-2 text-xs border border-base-border truncate max-w-full w-max';

/**
 * 表格单元格标签：横向顺排、自动换行，最多展示 maxVisible 个；超出 hover 展示完整列表。
 * 纯展示，不修改入参 items。
 *
 * @param {string[]} [items] 标签文案列表
 * @param {number} [maxVisible=4] 最多展示标签数（约两行内顺排）
 * @param {string} [tagClassName] 单个标签基础样式
 * @param {{ className?: string }} [tagProps] 追加到每个标签的 props（如列配色）
 */
const TruncatedTagGrid = ({
	items = [],
	maxVisible = DEFAULT_MAX_VISIBLE,
	tagClassName = DEFAULT_TAG_CLASS,
	tagProps = {},
}) => {
	if (!items?.length) {
		return EMPTY_DISPLAY;
	}

	const visibleItems = items.slice(0, maxVisible);
	const hasOverflow = items.length > maxVisible;
	const tagClass = [tagClassName, tagProps?.className].filter(Boolean).join(' ');

	const renderTag = (item, index) => (
		<div key={`${item}-${index}`} className={tagClass}>
			{item}
		</div>
	);

	const content = (
		<div className="flex flex-wrap gap-1 overflow-hidden max-h-[52px]">
			{visibleItems.map(renderTag)}
		</div>
	);

	if (!hasOverflow) {
		return content;
	}

	const tooltipTitle = (
		<div className="flex flex-wrap gap-1" onClick={(e) => e.stopPropagation()}>
			{items.map(renderTag)}
		</div>
	);

	return <FieldTips title={tooltipTitle}>{content}</FieldTips>;
};

export default TruncatedTagGrid;
