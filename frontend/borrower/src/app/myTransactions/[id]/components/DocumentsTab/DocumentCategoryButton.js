'use client';

import React from 'react';

function CategoryStatusTag({ completed, pendingCount }) {
	if (completed) {
		return (
			<span className="inline-flex items-center rounded-[2px] border border-[var(--green-light-100)] bg-green-light px-2 py-px text-xs leading-5 text-green-medium">
				√
			</span>
		);
	}
	if (pendingCount > 0) {
		return (
			<span className="inline-flex items-center rounded-[2px] border border-[var(--orange-medium)] bg-[var(--orange-light)] px-2 py-px text-xs leading-5 text-[var(--orange-dark)]">
				{pendingCount}
			</span>
		);
	}
	return null;
}

export default function DocumentCategoryButton({ category, isActive, onSelect }) {
	const { key, name, completed, pendingCount } = category;

	return (
		<div
			className={[
				'flex w-full cursor-pointer items-center justify-between border-l-2 px-4 py-4 text-left',
				isActive ? 'border-l-blue-dark bg-gray-light' : 'border-l-transparent',
			].join(' ')}
			onClick={() => onSelect(key)}
		>
			<span
				className={[
					'text-sm font-medium leading-[22px]',
					isActive ? 'text-foreground' : 'text-base-gray',
				].join(' ')}
			>
				{name}
			</span>
			<CategoryStatusTag completed={completed} pendingCount={pendingCount} />
		</div>
	);
}
