import React from 'react';
import { NewModal } from '@repo/shared-methods';
import { Alert, Radio } from 'antd';
const SettingModal = ({ setOpen, handleSearch, dueDay, setDueDay }) => {
	const handleSave = () => {
		handleSearch?.(dueDay);
		setOpen(false);
	};
	return (
		<NewModal
			width={572}
			title="Maturity alerts settings"
			onCancel={() => setOpen(false)}
			onOk={handleSave}
		>
			<div className="text-base-text font-medium text-xl pb-4">Due soon criteria</div>
			<Radio.Group
				block
				options={[
					{
						label: '30 days',
						value: 30,
					},
					{
						label: '60 days',
						value: 60,
					},
					{
						label: '90 days',
						value: 90,
					},
				]}
				defaultValue="Apple"
				optionType="button"
				onChange={(e) => setDueDay(e.target.value)}
				size="large"
				className="gap-6 [&_.ant-radio-button-wrapper]:!rounded-lg !pb-4"
				value={dueDay}
			/>
			<Alert
				message="Note:"
				description={`Selecting ${dueDay} days means loans due within ${dueDay} days will be marked as "Due Soon", and those beyond ${dueDay} days will be marked as "Overdue".`}
				type="warning"
			/>
		</NewModal>
	);
};

export default SettingModal;
