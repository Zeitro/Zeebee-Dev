'use client';
import { Button } from 'antd';
import React, { useState } from 'react';
import { PlusOutlined } from '@ant-design/icons';
import { AddExperience } from './AddExperience';
import AddThirdService from './AddThirdService';
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
			{['Add third party service', 'Edit third party service'].includes(open) && (
				<AddThirdService
					title={title || children}
					rowData={rowData}
					onCancel={() => setOpen('')}
				/>
			)}
		</>
	);
};
