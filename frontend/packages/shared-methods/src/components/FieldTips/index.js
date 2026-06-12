'use client';
import React from 'react';
import { QuestionCircleOutlined } from '@ant-design/icons';
import { Tooltip } from 'antd';

// 字段提示，内容超长浮框显示完整内容，或者用于问号内显示字段说明
// isShowDescription: 是否显示问号
export const FieldTips = ({
	children,
	title,
	className,
	isShowDescription,
	rootClassName,
	...props
}) => {
	return (
		<Tooltip
			title={title}
			classNames={{
				root: `!max-w-[300px] ${rootClassName}`,
			}}
			{...props}
		>
			{children}
			{isShowDescription && (
				<QuestionCircleOutlined
					className={['text-sm !text-base-warning', className].join(' ')}
				/>
			)}
		</Tooltip>
	);
};
