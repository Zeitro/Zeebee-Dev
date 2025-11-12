import React from 'react';
import { Row, Col, Space } from 'antd';
import { EditOutlined, DeleteOutlined, UserOutlined } from '@ant-design/icons';
import { assemblyAddress } from '@repo/shared-methods/utils';
import { BtnPopUpModal, DeleteBtn } from '@repo/shared-methods';
import AddModal from '../Add/borrower';
import { useBorrower } from './useCommon';
const BorrowerCardItem = ({
	isShowOperate = true,
	rowData,
	item,
	disabled,
	handleClickItem,
	getList,
	...props
}) => {
	const handleSearch = () => getList?.();
	const { loading, handleDelete } = useBorrower({
		handleSearch,
	});
	return (
		<Col xl={6} sm={24} xs={24} {...props}>
			<div
				onClick={() => {
					handleClickItem(item);
				}}
				style={{
					borderColor:
						rowData?.id && rowData?.id === item?.id ? '#0958D9' : 'var(--base-border)',
				}}
				className={[
					'border border-base-border rounded-md pl-2 py-3 pr-4',
					!disabled &&
						'cursor-pointer hover:!border-[#0958D9] [&:hover_.common-fill]:![color:#0958D9]',
				].join(' ')}
			>
				<div className="flex items-center gap-2">
					<div
						className={['common-fill text-[rgba(0,0,0,0.45)]'].join(' ')}
						style={{
							color:
								rowData?.id && rowData?.id === item?.id
									? '#0958D9'
									: 'rgba(0,0,0,0.45)',
						}}
					>
						<UserOutlined />
					</div>
					<p className="common-ellipsis text-base font-medium" title={item?.full_name}>
						{item?.full_name}
					</p>
				</div>
				<Row className="pl-2 pt-1 text-[rgba(0,0,0,0.65)]">
					<Col md={24} className="common-ellipsis" title={assemblyAddress(item?.address)}>
						{assemblyAddress(item?.address)}
					</Col>
				</Row>
				{isShowOperate && (
					<Space className="w-full justify-end pt-2">
						<BtnPopUpModal
							getContent={(setOpen) => (
								<AddModal
									setOpen={setOpen}
									handleSearch={async () => {
										const list = await handleSearch(rowData?.id);
										const newItem = list.find((i) => i.id === item?.id);
										// 修改的是选中的项
										if (rowData?.id === item?.id) {
											handleClickItem(newItem, list);
										}
									}}
									rowData={item}
								/>
							)}
							type="default"
							size="small"
							icon={<EditOutlined className="text-base-primary" />}
							disabled={disabled}
						></BtnPopUpModal>
						<DeleteBtn
							type="default"
							title="Are you sure delete this borrower?"
							content={item.full_name}
							onOk={() => handleDelete(item.id)}
							loading={loading}
							size="small"
							disabled={disabled}
						>
							<DeleteOutlined className="text-base-error" />
						</DeleteBtn>
					</Space>
				)}
			</div>
		</Col>
	);
};
export default BorrowerCardItem;
