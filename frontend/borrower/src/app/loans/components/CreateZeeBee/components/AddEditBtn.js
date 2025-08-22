'use client';
import { Button } from 'antd';
import React, { useState } from 'react';
import { PlusOutlined } from '@ant-design/icons';
import { AddExperience } from './AddExperience';
export const AddEditBtn = ({ title, children, rowData, ...props }) => {
	const [open, setOpen] = useState('');
	const handleClick = () => {
		setOpen(title || children);
	};

	return (
		<>
			<Button
				type="primary"
				icon={<PlusOutlined />}
				onClick={handleClick}
				className="w-max"
				{...props}
			>
				{children}
			</Button>
			{[
				'Add transaction history',
				'Add real estate',
				'Edit transaction history',
				'Edit real estate',
			].includes(open) && (
				<AddExperience
					title={title || children}
					rowData={rowData}
					onCancel={() => setOpen('')}
				/>
			)}
		</>
	);
};
