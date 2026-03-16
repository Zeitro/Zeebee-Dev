'use client';
import Details from '@repo/shared-methods/componentsUi/zeeBeeDetails';
import { PERMISSIONS } from '@config/permissions';
const ZeeBeeDetails = () => {
	return <Details PERMISSIONS={PERMISSIONS} />;
};
export default ZeeBeeDetails;
