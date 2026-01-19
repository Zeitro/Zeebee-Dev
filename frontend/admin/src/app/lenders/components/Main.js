import { NewModal, OrgProfile } from '@repo/shared-methods';
import { message, Spin } from 'antd';
import { useCallback, useRef, useState } from 'react';
import { createRequest } from '@repo/shared-methods/utils';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
import _ from 'lodash';

// 用于登录、注册等位置，首次注册用户填写组织信息，才可登录系统
const MainModal = ({ rowData, setOpen, successCallback }) => {
	const { request } = createRequest();
	const orgRef = useRef();
	const [confirmLoading, setConfirmLoading] = useState(false);
	const [data, setData] = useState({});
	const [loading, setLoading] = useState(false);
	const handleSave = async () => {
		try {
			setConfirmLoading(true);
			await orgRef?.current?.handleSaveOrgProfile?.();
			await successCallback?.();
			message.destroy();
			message.open({
				content: `Organization ${rowData?.id ? 'Edited' : 'Added'} successfully.`,
				type: 'success',
			});
			setConfirmLoading(false);
			setOpen(false);
		} catch {
			setConfirmLoading(false);
		}
	};

	const getDetails = useCallback(async () => {
		try {
			setLoading(true);
			const res = await request(`/organization/${rowData?.id}`);
			// 组装数据，方便edit时提交
			const metadata = _.cloneDeep(res?.metadata?.lender);
			res.metadata = metadata;
			setData(res);
			setLoading(false);
		} catch {
			setLoading(false);
		}
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [rowData?.id]);

	useDeferredSyncEffect({
		predicate: () => rowData?.id,
		deferred: () => {
			getDetails();
		},
		deps: [rowData, getDetails],
		delay: 0,
	});
	return (
		<NewModal
			width={800}
			title={(rowData?.id ? 'Edit' : 'Add') + ' Lender Organization'}
			okText="Save"
			confirmLoading={confirmLoading}
			onOk={handleSave}
			onCancel={() => setOpen(false)}
		>
			{/* <Alert
				message="Please fill in your organization's complete information before you can log in to this system."
				banner
				className="!mb-4"
			/> */}
			<Spin spinning={loading}>
				<OrgProfile
					ref={orgRef}
					api={
						rowData?.id ? `/organization/lender/${rowData?.id}` : '/organization/lender'
					}
					method={rowData?.id ? 'PUT' : 'POST'}
					rowData={data}
					metadataName={['lender_metadata']}
					initMetadataName={['metadata', 'lender']}
					className={'h-[65vh] overflow-hidden overflow-y-auto'}
					leftCol={{ lg: 8, sm: 24, xl: 8 }}
					rightCol={{ lg: 16, sm: 24, xl: 16 }}
					isShowAccountInput={!rowData?.id}
					passwordName={['owner', 'password']}
					confirmPasswordName={['owner', 'confirm_password']}
					firstName={['owner', 'first_name']}
					lastName={['owner', 'last_name']}
					emailName={['owner', 'email']}
					orgName={rowData?.id ? ['name'] : ['organization', 'name']}
					contactName={rowData?.id ? ['contact_name'] : ['organization', 'contact_name']}
					contactEmailName={
						rowData?.id ? ['contact_email'] : ['organization', 'contact_email']
					}
					contactPhoneName={
						rowData?.id ? ['contact_phone'] : ['organization', 'contact_phone']
					}
					taglineName={rowData?.id ? ['tagline'] : ['organization', 'tagline']}
					websiteName={rowData?.id ? ['website'] : ['organization', 'website']}
					loanAmountMinName={[
						...(rowData?.id ? [] : ['organization']),
						'metadata',
						'loan_amount_range',
						'min',
					]}
					loanAmountMaxName={[
						...(rowData?.id ? [] : ['organization']),
						'metadata',
						'loan_amount_range',
						'max',
					]}
					loanRateMinName={[
						...(rowData?.id ? [] : ['organization']),
						'metadata',
						'loan_rate_range',
						'min',
					]}
					loanRateMaxName={[
						...(rowData?.id ? [] : ['organization']),
						'metadata',
						'loan_rate_range',
						'max',
					]}
					loanPointMinName={[
						...(rowData?.id ? [] : ['organization']),
						'metadata',
						'loan_point_range',
						'min',
					]}
					loanPointMaxName={[
						...(rowData?.id ? [] : ['organization']),
						'metadata',
						'loan_point_range',
						'max',
					]}
					loanTermMinName={[
						...(rowData?.id ? [] : ['organization']),
						'metadata',
						'loan_term_range',
						'min',
					]}
					loanTermMaxName={[
						...(rowData?.id ? [] : ['organization']),
						'metadata',
						'loan_term_range',
						'max',
					]}
					logoName={rowData?.id ? ['logo'] : ['organization', 'logo']}
				/>
			</Spin>
		</NewModal>
	);
};
export default MainModal;
