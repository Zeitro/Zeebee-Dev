import React from 'react';
import { Button, Space } from 'antd';
import { ApproveModal, RejectedModal } from '../ApproveRejectModal';
const HeaderBtn = ({ title, children }) => {
	return (
		<div className="flex justify-between items-center border-b border-base-border py-4 px-6">
			<div className="flex-1 overflow-hidden text-xl">{title || children}</div>
			<Space className="shrink-0">
				<RejectedModal />
				<ApproveModal />
			</Space>
		</div>
	);
};
export default HeaderBtn;
