import { NewModal } from '@repo/shared-methods';
import { Alert, message, Select } from 'antd';
import { useState } from 'react';
import { createRequest } from '@repo/shared-methods/utils';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
const ReassignLender = ({ setOpen, data, successCallback }) => {
	const { request } = createRequest();
	const [confirmLoading, setConfirmLoading] = useState(false);
	const [loading, setLoading] = useState(false);
	const [value, setValue] = useState(null);
	const [options, setOptions] = useState([]);

	const getOptions = async () => {
		try {
			setLoading(true);
			const params = new URLSearchParams({ type: 'Lender' });
			const { list } = await request(`/organization/select_options?${params}`);
			setOptions(list);
			setLoading(false);
		} catch {
			setLoading(false);
		}
	};
	const handleSave = async () => {
		try {
			if (!value) {
				message.destroy();
				message.open({
					type: 'error',
					content: 'Please select the lender to be assigned.',
				});
				return;
			}
			setConfirmLoading(true);
			await request(`/loan/${data?.id}/reassign`, {
				method: 'POST',
				body: {
					lender_organization_id: value,
				},
			});

			message.open({
				type: 'success',
				content: 'Assigned successfully!',
			});
			successCallback?.();
			setConfirmLoading(false);
			setOpen(false);
		} catch {
			setConfirmLoading(false);
		}
	};

	useDeferredSyncEffect({
		deferred: () => {
			getOptions();
		},
		deps: [],
		delay: 0,
	});
	return (
		<NewModal
			width={500}
			title="Reassign Lender"
			onCancel={() => setOpen(false)}
			onOk={handleSave}
			confirmLoading={confirmLoading}
		>
			<Alert
				message={
					<div className="text-base">
						Current lender:
						<span className="font-medium pl-2">{data?.lender_company}</span>
					</div>
				}
				type="info"
			/>
			<div className="pt-4 text-base">New lender:</div>
			<Select
				placeholder="Select a lender"
				style={{ width: '100%' }}
				onChange={(value) => setValue(value)}
				value={value}
				options={options}
				size="large"
				loading={loading}
				fieldNames={{ label: 'name', value: 'id' }}
			/>
		</NewModal>
	);
};
export default ReassignLender;
