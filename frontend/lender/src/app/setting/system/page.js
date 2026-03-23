'use client';
// import { InfoCircleOutlined, MailTwoTone } from '@ant-design/icons';
// import ContactEmail from './components/ContactEmail';
// import { FieldTips } from '@repo/shared-methods';
import PreApprovedLetterSetting from './components/PreApprovalLetter';
import { PERMISSIONS } from '@config/permissions';
import { PermissionSlot } from '@repo/shared-methods';
const SystemSetting = () => {
	return (
		<div className="py-6 px-8 flex-1 overflow-hidden w-full">
			{/* <div className="flex items-center pb-6">
				<MailTwoTone style={{ fontSize: 24 }} />
				<span className="text-base font-medium pl-2">Contact Email</span>
				<FieldTips title="This email is used for request quote">
					<InfoCircleOutlined className={['text-sm !text-base-warning ml-1'].join(' ')} />
				</FieldTips>
			</div>
			<ContactEmail /> */}
			<div className="text-xl font-medium pb-4">System Settings</div>
			<PermissionSlot permission={PERMISSIONS.SETTING_SYSTEM.PRE_APPROVED_LETTER}>
				<PreApprovedLetterSetting />
			</PermissionSlot>
		</div>
	);
};
export default SystemSetting;
