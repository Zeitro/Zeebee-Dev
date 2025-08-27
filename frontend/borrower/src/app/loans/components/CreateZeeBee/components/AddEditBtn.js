'use client';
import { Button } from 'antd';
import React, { useState } from 'react';
import { PlusOutlined } from '@ant-design/icons';
import { AddExperience } from './AddExperience';
import AddThirdService from './AddThirdService';
export const AddEditBtn = ({ title, commonType, children, rowData, onOk = () => {}, ...props }) => {
	const [open, setOpen] = useState('');
	const handleClick = () => {
		setOpen(commonType);
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
			{['transactionHistory', 'realEstate'].includes(open) && (
				<AddExperience
					title={title || children}
					rowData={rowData}
					onCancel={() => setOpen('')}
					onOk={(values, v) => onOk(values, v, commonType)}
					type={commonType}
				/>
			)}
			{['thirdPartyService'].includes(open) && (
				<AddThirdService
					title={title || children}
					rowData={rowData}
					onCancel={() => setOpen('')}
					onOk={(values, v) => onOk(values, v, commonType)}
				/>
			)}
		</>
	);
};
