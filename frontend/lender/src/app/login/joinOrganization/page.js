'use client';
import { InputNameCombination } from '@repo/shared-methods';
import { WarningOutlined } from '@ant-design/icons';
import { Button, Form, message, Row, Spin } from 'antd';
import { useSearchParams } from 'next/navigation';
import { useEffect, useMemo, useState } from 'react';
import { createRequest } from '@repo/shared-methods/utils';
import useContent from '../components/useContent';

const formCardClassName =
	'w-[400px] rounded-lg bg-white !px-12 !py-12 shadow-[0_1px_2px_0_rgba(0,0,0,0.03), 0_1px_6px_-1px_rgba(0,0,0,0.02), 0_2px_4px_0_rgba(0,0,0,0.02)]';

// 邀请加入组织，邮件链接确认加入访问的页面
const JoinOrganization = () => {
	const { loginByToken } = useContent();
	const { request } = useMemo(() => createRequest(), []);
	const searchParams = useSearchParams();
	const email = searchParams.get('email');
	const organization_id = searchParams.get('organization_id');
	const code = searchParams.get('code');
	const organizationName = searchParams.get('organization_name');
	const first_name = searchParams.get('first_name');
	const last_name = searchParams.get('last_name');
	const [form] = Form.useForm();
	const [isCheckingInvitation, setIsCheckingInvitation] = useState(true);
	const [errTips, setErrTips] = useState('');

	/** 进入页面时校验邀请码是否过期（POST /organization_user/check_invitation） */
	useEffect(() => {
		const checkInvitation = async () => {
			if (!email || !code || !organization_id) {
				setIsCheckingInvitation(false);
				setErrTips('Invalid user information. Please contact the administrator.');
				return;
			}
			try {
				const data = await request('/organization_user/check_invitation', {
					method: 'POST',
					body: {
						code,
						email,
						organization_id: Number(organization_id),
					},
					isHideTips: true,
				});
				setErrTips(
					data?.is_expired
						? 'The invitation code has been expired. Please contact the organization administrator.'
						: '',
				);
			} catch (error) {
				setErrTips(error.message);
			} finally {
				setIsCheckingInvitation(false);
			}
		};
		checkInvitation();
	}, [code, email, organization_id, request]);

	const handleJoinAndLogin = async () => {
		if (errTips) return;
		const values = await form.validateFields();
		const access = await request('/organization_user/accept_invitation', {
			method: 'POST',
			body: { organization_id: Number(organization_id), email, code, ...values },
		});
		loginByToken(access, () => {
			message.open({
				type: 'success',
				content: `Successfully joined the ${organizationName} organization. Logging in...`,
			});
		});
		// setCookie('token', access?.access_token, access?.access_token_expire * 1000);
		// await updateAccess(access);
		// await getUserInfo();
		// router.push('/');
	};

	if (isCheckingInvitation) {
		return (
			<div className="flex-1 w-full bg-[#F5F5F5] overflow-hidden flex justify-center items-center">
				<Spin />
			</div>
		);
	}

	if (errTips) {
		return (
			<div className="flex-1 w-full bg-[#F5F5F5] overflow-hidden flex justify-center items-center">
				<div className={`${formCardClassName} text-center`}>
					<WarningOutlined style={{ fontSize: 48, color: 'var(--base-warning)' }} />
					<div className="text-base pt-4 text-base-gray ">
						Note: <span className="font-medium text-sm">{errTips}</span>
					</div>
				</div>
			</div>
		);
	}

	return (
		<div className="flex-1 w-full bg-[#F5F5F5] overflow-hidden flex justify-center items-center ">
			<Form
				name="registerForm"
				size="large"
				layout="vertical"
				autoComplete="off"
				form={form}
				initialValues={{ email, last_name, first_name }}
				className={formCardClassName}
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
