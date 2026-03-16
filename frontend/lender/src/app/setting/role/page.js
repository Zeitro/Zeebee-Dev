'use client';
import CommonSettingRole from '@repo/shared-methods/componentsUi/SettingRole';
import { PERMISSIONS } from '@config/permissions';
const SettingRole = () => {
	return <CommonSettingRole api="/role/page" systemType="Lender" PERMISSIONS={PERMISSIONS} />;
};
export default SettingRole;
