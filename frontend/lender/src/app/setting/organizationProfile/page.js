'use client';
import { OrgProfile } from '@repo/shared-methods';
import { Button, message } from 'antd';
import { useRef, useState } from 'react';
import { useUserStore } from '@/store/user';
const OrganizationProfile = () => {
	const { organization, getOrganization } = useUserStore();
	const orgRef = useRef();
	const [confirmLoading, setConfirmLoading] = useState(false);

	const handleSave = async () => {
		try {
			setConfirmLoading(true);
			await orgRef?.current?.handleSaveOrgProfile?.();
			await getOrganization();
			message.destroy();
			message.open({
				content: 'Organization information saved successfully.',
				type: 'success',
			});
			setConfirmLoading(false);
		} catch {
			setConfirmLoading(false);
		}
	};

	return (
		<div className=" py-6 px-8 flex-1 overflow-hidden w-full overflow-y-auto">
			<OrgProfile
				ref={orgRef}
				rowData={organization}
				api="/organization/profile"
				metadataName={['lender_metadata']}
				initMetadataName={['metadata', 'lender']}
			>
				<div className="flex justify-end">
					<Button type="primary" loading={confirmLoading} onClick={handleSave}>
						{confirmLoading ? 'Saving...' : 'Save'}
					</Button>
				</div>
			</OrgProfile>
		</div>
	);
};
export default OrganizationProfile;
