'use client';
import CommonSettingUser from '@repo/shared-methods/componentsUi/SettingUser';
import UserMainModal from '../components/UserMain';
import { PERMISSIONS } from '@config/permissions';
const SettingUser = () => {
	return (
		<CommonSettingUser
			PERMISSIONS={PERMISSIONS}
			getNewMainModal={(props) => <UserMainModal {...props} />}
		/>
	);
};
export default SettingUser;
