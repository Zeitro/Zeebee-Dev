import React from 'react';
import { Row, Col, Space } from 'antd';
import { EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { formatTINOrSSN } from '@repo/shared-methods/utils';
import { BtnPopUpModal, DeleteBtn } from '@repo/shared-methods';
import BorrowingEntitySvg from '@repo/shared-methods/assets/images/borrowingEntity.svg';
import AddModal from '../Add/borrowingEntity';
import useBorrowingEntity from '@/app/management/borrowingEntity/components/useBorrowingEntity';
const BorrowingEntityCardItem = ({
	isShowOperate = true,
	rowData,
	item,
	disabled,
	handleClickItem,
	getList,
	...props
}) => {
	const handleSearch = (id) => getList?.(id);
	const { loading, handleDelete } = useBorrowingEntity({
		handleSearch,
	});
	return (
		<Col
			xl={6}
			// lg={list.length > 2 ? 24 : 6}
			sm={24}
			xs={24}
			{...props}
		>
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
						<BorrowingEntitySvg />
					</div>
					<p className="common-ellipsis text-base font-medium" title={item?.name}>
						{item?.name}
					</p>
				</div>
				<Row className="pl-2 pt-1 text-[rgba(0,0,0,0.65)]">
					<Col md={12} className="common-ellipsis" title={item?.type}>
						{item?.type || '--'}
					</Col>
					<Col md={12} className="common-ellipsis" title={formatTINOrSSN(item?.tin)}>
						{formatTINOrSSN(item?.tin) || '--'}
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
							title="Are you sure delete this borrowing entity?"
							content={item.full_name}
							onOk={() => handleDelete(item.id)}
							loading={loading}
							size="small"
							className={'!border-base-border hover:!border-base-error !px-1'}
							disabled={disabled}
						>
							<DeleteOutlined className="text-base-error" />
						</DeleteBtn>
						{/* <Button
							size="small"
							icon={<EditOutlined className="text-base-primary" />}
						></Button>
						<Button
							size="small"
							icon={<DeleteOutlined className="text-base-error" />}
						></Button> */}
					</Space>
				)}
			</div>
		</Col>
	);
};
export default BorrowingEntityCardItem;
