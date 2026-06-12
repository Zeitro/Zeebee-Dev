'use client';
import { GoogleOutlined, InfoCircleOutlined, MailTwoTone } from '@ant-design/icons';
import ContactEmail from './components/ContactEmail';
import DriveIntegration from './components/DriveIntegration';
import { FieldTips } from '@repo/shared-methods';
const SystemSetting = () => {
	return (
		<div className="py-6 px-8 flex-1 overflow-hidden w-full xl:w-1/2">
			<div className="flex items-center pb-6">
				<MailTwoTone style={{ fontSize: 24 }} />
				<span className="text-base font-medium pl-2">Contact Email</span>
				<FieldTips title="This email is used for request quote">
					<InfoCircleOutlined className={['text-sm !text-base-warning ml-1'].join(' ')} />
				</FieldTips>
			</div>
			<ContactEmail />

			<div className="flex items-center pb-6 pt-10">
				<GoogleOutlined style={{ fontSize: 24, color: '#4285F4' }} />
				<span className="text-base font-medium pl-2">Google Drive</span>
				<FieldTips title="Connect Google Drive to sync transaction documents">
					<InfoCircleOutlined className={['text-sm !text-base-warning ml-1'].join(' ')} />
				</FieldTips>
			</div>
			<DriveIntegration />
		</div>
	);
};
export default SystemSetting;
