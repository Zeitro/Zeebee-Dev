import React from 'react';
import { Button, Space } from 'antd';
import { ApproveRejectModal } from '../ApproveRejectModal';
const HeaderBtn = ({
	title,
	children,
	isHaveRejectBtn,
	getContent = () => {},
	RejectedProps = {},
	pending_status,
	setNextStep,
	setRejectStep,
	disabledOperate,
	id,
	...ApproveProps
}) => {
	return (
		<div className="flex justify-between items-center border-b-4 border-[#f6f6f6] py-4 px-6 position-sticky top-0">
			<div
				className="flex-1 overflow-hidden text-xl font-medium"
				style={{ color: 'rgba(0, 0, 0, 0.65)' }}
			>
				{title || children}
			</div>
			{getContent() ?? (
				<Space className="shrink-0" size={16}>
					{isHaveRejectBtn && (
						<ApproveRejectModal
							{...RejectedProps}
							action="Reject"
							pending_status={pending_status}
							id={id}
							setNextStep={setRejectStep}
							disabled={disabledOperate}
						/>
					)}
					<ApproveRejectModal
						{...ApproveProps}
						action="Approve"
						pending_status={pending_status}
						id={id}
						setNextStep={setNextStep}
						disabled={disabledOperate}
					/>
				</Space>
			)}
		</div>
	);
};
export default HeaderBtn;
