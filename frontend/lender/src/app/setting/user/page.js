'use client';
import CommonSettingUser from '@repo/shared-methods/componentsUi/SettingUser';
import { PERMISSIONS } from '@config/permissions';
const SettingUser = () => {
	return <CommonSettingUser PERMISSIONS={PERMISSIONS} />;
};
export default SettingUser;
