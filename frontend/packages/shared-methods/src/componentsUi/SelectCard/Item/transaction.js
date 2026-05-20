'use client';

import React from 'react';
import { Row, Col, Space, Tag } from 'antd';
import { EditOutlined, DeleteOutlined, SwapOutlined } from '@ant-design/icons';
import { formatDate, EMPTY_DISPLAY } from '@repo/shared-methods/utils';
import { BtnPopUpModal, DeleteBtn } from '@repo/shared-methods';
// import AddModal from '../Add/transaction';
// import { useTransaction } from './useCommon';

function transactionStatusTag(status) {
	if (!status) return null;
	const color = status === 'Pending' ? 'gold' : status === 'Cancelled' ? 'red' : 'success';
	return (
		<Tag bordered={false} color={color} className="!m-0 shrink-0 !text-xs leading-5">
			{status}
		</Tag>
	);
}

function formatTransactionId(item) {
	if (item?.transaction_no) return item.transaction_no;
	if (item?.id != null && item?.id !== '') return `TS-${item.id}`;
	return EMPTY_DISPLAY;
}

const TransactionCardItem = ({
	isShowOperate = true,
	rowData,
	item,
	disabled,
	handleClickItem,
	getList,
	...props
}) => {
	const handleSearch = () => getList?.();
	// const { loading, handleDelete } = useTransaction({
	// 	handleSearch,
	// 	handleClickItem,
	// 	rowData,
	// });

	const isActive = rowData?.id && rowData?.id === item?.id;
	const activeColor = isActive ? '#0958D9' : 'rgba(0,0,0,0.45)';
	const coeRaw = item?.coe_date ?? item?.coeDate;
	const coeText = coeRaw ? formatDate(coeRaw) : EMPTY_DISPLAY;

	return (
		<Col xl={6} sm={24} xs={24} {...props}>
			<div
				onClick={() => handleClickItem(item)}
				style={{
					borderColor: isActive ? '#0958D9' : 'var(--base-border)',
				}}
				className={[
					'border border-base-border rounded-md pl-2 py-3 pr-4',
					!disabled &&
						'cursor-pointer hover:!border-[#0958D9] [&:hover_.common-fill]:![color:#0958D9]',
				].join(' ')}
			>
				<div className="flex items-center gap-2 min-w-0">
					<div
						className={['common-fill shrink-0'].join(' ')}
						style={{ color: activeColor }}
					>
						<SwapOutlined className="text-lg" />
					</div>
					<p
						className="common-ellipsis text-base font-medium flex-1 min-w-0"
						title={item?.type}
					>
						{item?.type || EMPTY_DISPLAY}
					</p>
					{transactionStatusTag(item?.status)}
				</div>
				<Row className="pl-2 pt-1 text-[rgba(0,0,0,0.65)]">
					<Col span={12} className="common-ellipsis" title={formatTransactionId(item)}>
						{formatTransactionId(item)}
					</Col>
					<Col span={12} className="common-ellipsis text-right" title={coeText}>
						{coeText}
					</Col>
				</Row>
				{/* {isShowOperate && (
					<Space className="w-full justify-end pt-2">
						<BtnPopUpModal
							getContent={(setOpen) => (
								<AddModal
									setOpen={setOpen}
									rowData={item}
									handleSearch={async () => {
										const list = await handleSearch(rowData?.id);
										const newItem = list?.find?.((i) => i.id === item?.id);
										if (rowData?.id === item?.id && newItem) {
											handleClickItem(newItem, list);
										}
									}}
								/>
							)}
							type="default"
							size="small"
							icon={<EditOutlined className="text-base-primary" />}
							disabled={disabled}
						/>
						<DeleteBtn
							type="default"
							title="Are you sure delete this transaction?"
							content={formatTransactionId(item)}
							onOk={() => handleDelete(item.id)}
							loading={loading}
							size="small"
							className="!border-base-border hover:!border-base-error !px-1"
							disabled={disabled}
						>
							<DeleteOutlined className="text-base-error" />
						</DeleteBtn>
					</Space>
				)} */}
			</div>
		</Col>
	);
};

export default TransactionCardItem;
