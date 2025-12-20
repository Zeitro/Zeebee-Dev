'use client';
import { MailTwoTone } from '@ant-design/icons';
import ContactEmail from './components/ContactEmail';
const SystemSetting = () => {
	return (
		<div className="py-6 px-8 flex-1 overflow-hidden w-full xl:w-1/2">
			<div className="flex items-center pb-6">
				<MailTwoTone style={{ fontSize: 24 }} />
				<span className="text-base font-medium pl-2">Contact Email</span>
			</div>
			<ContactEmail />
		</div>
	);
};
export default SystemSetting;
