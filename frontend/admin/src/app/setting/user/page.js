'use client';
import CommonSettingUser from '@repo/shared-methods/componentsUi/SettingUser';
import UserMainModal from '../components/UserMain';
const SettingUser = () => {
	return <CommonSettingUser getNewMainModal={(props) => <UserMainModal {...props} />} />;
};
export default SettingUser;
