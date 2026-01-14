import { NewModal, OrgProfile } from '@repo/shared-methods';
import { Alert, message } from 'antd';
import { useRef, useState } from 'react';

// 用于登录、注册等位置，首次注册用户填写组织信息，才可登录系统
const MainModal = ({ rowData, setOpen, successCallback }) => {
	const orgRef = useRef();
	const [confirmLoading, setConfirmLoading] = useState(false);
	const [okText, setOkText] = useState('Save & Login');
	const handleSave = async () => {
		try {
			setConfirmLoading(true);
			setOkText('Saving...');
			await orgRef?.current?.handleSaveOrgProfile?.();
			setOkText('Logging in...');
			successCallback?.(() => {
				message.destroy();
				message.open({
					content: 'Organization information saved successfully.',
					type: 'success',
				});
				setConfirmLoading(false);
				setOpen(false);
			});
		} catch {
			setOkText('Save & Login');
			setConfirmLoading(false);
		}
	};
	return (
		<NewModal
			width={800}
			title={'Organization profile'}
			okText={okText}
			confirmLoading={confirmLoading}
			onOk={handleSave}
			onCancel={() => setOpen(false)}
		>
			<Alert
				message="Please fill in your organization's complete information before you can log in to this system."
				banner
				className="!mb-4"
			/>
			<OrgProfile
				ref={orgRef}
				api="/organization/profile"
				rowData={rowData}
				metadataName={['lender_metadata']}
				initMetadataName={['metadata', 'lender']}
				className={'h-[65vh] overflow-hidden overflow-y-auto'}
				leftCol={{ lg: 8, sm: 24, xl: 8 }}
				rightCol={{ lg: 16, sm: 24, xl: 16 }}
			/>
		</NewModal>
	);
};
export default MainModal;
