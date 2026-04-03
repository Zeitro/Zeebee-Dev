'use client';
import { InputNameCombination } from '@repo/shared-methods';
import { Button, Form, message, Row } from 'antd';
import { useSearchParams } from 'next/navigation';
import { createRequest } from '@repo/shared-methods/utils';
import useData from '../components/useData';
// 邀请加入组织，邮件链接确认加入访问的页面
const JoinOrganization = () => {
	const { getUserPermissions } = useData();
	const { request } = createRequest();
	const searchParams = useSearchParams();
	const email = searchParams.get('email');
	const organizationName = searchParams.get('organization_name');
	const first_name = searchParams.get('first_name');
	const last_name = searchParams.get('last_name');
	const [form] = Form.useForm();

	const handleJoinAndLogin = async () => {
		const values = await form.validateFields();
		const organization_id = searchParams.get('organization_id');
		const code = searchParams.get('code');
		const access = await request('/organization_user/accept_invitation', {
			method: 'POST',
			body: { organization_id: Number(organization_id), email, code, ...values },
		});
		// await updateAccess(access);
		await getUserPermissions(access, () => {
			message.open({
				type: 'success',
				content: `Successfully joined the ${organizationName} organization. Logging in...`,
			});
		});
		// setCookie('token', access?.access_token, access?.access_token_expire * 1000);
		// await getUserInfo();
		// router.push('/');
	};
	return (
		<div className="flex-1 w-full bg-[#F5F5F5] overflow-hidden flex justify-center items-center ">
			<Form
				name="registerForm"
				size="large"
				layout="vertical"
				autoComplete="off"
				form={form}
				initialValues={{ email, last_name, first_name }}
				className="w-[400px] rounded-lg bg-white !px-12 !py-12 shadow-[0_1px_2px_0_rgba(0,0,0,0.03), 0_1px_6px_-1px_rgba(0,0,0,0.02), 0_2px_4px_0_rgba(0,0,0,0.02)]"
			>
				<div className="text-base text-[rgba(0,0,0,0.65)] font-medium text-center">
					You are accepting the invitation to join {organizationName} in
				</div>
				<div className="text-xl font-medium pt-4 pb-8 text-center">{email}</div>
				<div className="text-base text-[rgba(0,0,0,0.65)] font-medium text-center pb-4">
					Input your name in this organization
				</div>
				<Row gutter={[16]}>
					<InputNameCombination
						layoutCol={{ span: 24 }}
						firstName={['first_name']}
						lastName={['last_name']}
						hideMiddleName={true}
					/>
					{/* <Col span={24} className='flex'>
						
					</Col> */}
					{/* <Col span={24}>
						Choosing to join means you agree to the
						<span className="text-base-primary cursor-pointer" onClick={() => {}}>
							&nbsp;Terms of Service&nbsp;
						</span>
						and
						<span className="text-base-primary cursor-pointer" onClick={() => {}}>
							&nbsp;Privacy Policy&nbsp;
						</span>
					</Col> */}
				</Row>
				<div className="flex justify-center">
					<Button type="primary" onClick={handleJoinAndLogin}>
						Join
					</Button>
				</div>
			</Form>
		</div>
	);
};
export default JoinOrganization;
