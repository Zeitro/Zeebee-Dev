'use client';

import React from 'react';

const DetailSection = ({ icon, title, extra, bordered = true, children }) => {
	return (
		<div
			className={[
				'flex flex-col gap-6 py-4',
				bordered ? 'border-b border-base-border pb-6' : '',
			].join(' ')}
		>
			<div className="flex items-center justify-between gap-4">
				<div className="flex items-center gap-2 text-sm font-medium leading-[22px] text-base-gray [&_.anticon]:text-base">
					{icon}
					<span>{title}</span>
				</div>
				{extra}
			</div>
			{children}
		</div>
	);
};

export default DetailSection;
