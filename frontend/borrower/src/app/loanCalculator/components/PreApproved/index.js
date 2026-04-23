'use client';
import PreApprovedSuccess from './Success';
import PhoneSvg from '@/assets/images/phone.svg';
import { Button, message } from 'antd';
import { useRouter } from 'next/navigation';
import { DoubleLeftOutlined } from '@ant-design/icons';
import { createRequest } from '@repo/shared-methods/utils';
import { useLoanStore } from '@/store/loan';
import OtherLender from './OtherLender';
// preApproved-成功失败展示
const PreApproved = ({ result, setIsResult }) => {
	const router = useRouter();
	const { request } = createRequest();
	const { preApprovedResult, clearExpiredPreApprovedInStorage } = useLoanStore();
	const handleClick = async (lender_company_id, isOtherLender) => {
		if (isOtherLender) {
			// 清除过期的preApproved
			clearExpiredPreApprovedInStorage();
			// 发送通知: 没有request_id会重复发通知，有则24小时内只发1次
			await request('/loan/send_pre_approval_loan_notification', {
				method: 'POST',
				body: {
					lender_company_id,
					request_id: preApprovedResult?.primary?.request_id || '',
				},
			});
		}
		message.open({
			type: 'success',
			content: 'Notification sent successfully.',
		});
		router.push(
			`/login/register?email=${result?.email}&request_id=${result?.primary?.request_id}&lender_id=${lender_company_id}`,
		);
	};
	const back = (
		<Button
			className="text-[#0958D9] text-base font-medium !absolute left-0"
			onClick={() => {
				setIsResult(false);
			}}
			icon={<DoubleLeftOutlined />}
		>
			Back
		</Button>
	);
	return (
		<div className="w-full px-[10%] lg:px-[12%] 2xl:px-[20%] my-10 overflow-hidden overflow-y-auto flex flex-col gap-4 [&::-webkit-scrollbar-thumb]:border-[#E4F3FF] [&::-webkit-scrollbar]:bg-[#E4F3FF]">
			{result?.primary?.is_approved ? (
				<PreApprovedSuccess result={result} handleClick={handleClick} back={back} />
			) : (
				<div className="mx-auto py-6 px-4 bg-white rounded-lg">
					<div className="text-3xl/10 font-medium text-[#0958D9] pb-2 text-center relative">
						{back}
						Application Under Review
					</div>
					<div className="text-base text-black-dark font-medium text-center">
						We've carefully reviewed your application
					</div>
					<div className="text-xl/[22px] font-medium text-[#1D2129] pt-10 px-2">
						Application status update
					</div>
					<div className="text-base text-black-dark py-4 px-2">
						Thank you for your interest in our loan services. After careful
						consideration of your application, we were unable to approve your request at
						this time based on our current lending criteria
					</div>
					<div className="text-base text-black-dark pb-10 px-2">
						We understand this may be disappointing, and we want to help you explore
						your options.
					</div>
					<div className="py-4 px-6 bg-[#E6F4FF] rounded-md mb-10">
						<div className="flex items-center gap-2 text-base-primary pb-2">
							<PhoneSvg style={{ width: '20px', height: '20px' }} />
							<p>A Loan Officer Will Contact You</p>
						</div>
						<div className="text-black-dark">
							One of our experienced loan officers will reach out to you within the
							next 2-3 business days to discuss your application in detail. They will
							explain the factors that influenced this decision and explore
							alternative solutions that may better fit your current financial
							situation.
						</div>
					</div>
					<div className="flex justify-center">
						<Button onClick={() => handleClick(result?.primary?.lender_company_id)}>
							Send loan request & Create your account
						</Button>
					</div>
				</div>
			)}
			<OtherLender result={result} handleClick={handleClick} />
		</div>
	);
};
export default PreApproved;
