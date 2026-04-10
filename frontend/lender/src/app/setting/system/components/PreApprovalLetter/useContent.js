import { message, Switch } from 'antd';
import { useUserStore } from '@repo/shared-methods/store';
import { useMemo, useRef, useState } from 'react';
import { createRequest } from '@repo/shared-methods/utils';

const useContent = () => {
	const { request } = createRequest();
	const switchRef = useRef();
	const { organization, getOrganization } = useUserStore();
	const [loading, setLoading] = useState(false);
	const preApprovalLetterSetting = useMemo(() => {
		const list = organization?.organization_settings;
		return list?.find((item) => item.purpose === 'preApprovedLetter');
	}, [organization]);
	const handleClickPreApprovalLetter = async (enabled) => {
		try {
			setLoading(true);
			await request(`/organization_setting/${preApprovalLetterSetting?.id}/status`, {
				method: 'PATCH',
				body: {
					enabled,
				},
			});
			await getOrganization();
			setLoading(false);
			message.destroy();
			message.open({
				type: 'success',
				content: enabled ? 'Enabled successful！' : 'Disabled successful！',
			});
		} catch {
			setLoading(false);
		}
	};
	const switchBtn = (
		<Switch
			ref={switchRef}
			checked={preApprovalLetterSetting?.enabled}
			loading={loading}
			onChange={handleClickPreApprovalLetter}
		/>
	);
	return { switchBtn, switchRef, preApprovalLetterSetting };
};
export default useContent;
