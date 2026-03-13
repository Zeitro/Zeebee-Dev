'use client';
import SettingRole from '@repo/shared-methods/componentsUi/SettingRole';
import { PERMISSIONS } from '@config/permissions';
const Loans = () => {
	return <SettingRole PERMISSIONS={PERMISSIONS} api="/role/page" systemType="Admin" />;
};
export default Loans;
