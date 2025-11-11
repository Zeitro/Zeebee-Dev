import React from 'react';
import { Row, Col, Space, message } from 'antd';
import { EditOutlined, DeleteOutlined } from '@ant-design/icons';
import HomeSvg from '@/assets/images/home2.svg';
import { assemblyAddress, createRequest, formatNumber } from '@repo/shared-methods/utils';
import { BtnPopUpModal, DeleteBtn } from '@repo/shared-methods';
import AddModal from '../Add/property';
// import useBorrowingEntity from '@/app/management/borrowingEntity/components/useBorrowingEntity';
const CardItem = ({
	isShowOperate = true,
	rowData,
	item,
	disabled,
	handleClickItem,
	getList,
	...props
}) => {
	const { request } = createRequest();
	const [loading, setLoading] = React.useState(false);
	const handleSearch = () => getList?.();
	const handleDelete = async () => {
		try {
			setLoading(true);
			await request(`/property/${item?.id}`, { method: 'DELETE' });
			handleSearch();
			message.destroy();
			message.open({
				type: 'success',
				content: 'Delete successful！',
			});
			setLoading(false);
		} catch {
			setLoading(false);
		}
	};
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
						<HomeSvg />
					</div>
					<p
						className="common-ellipsis text-base font-medium"
						title={assemblyAddress(item?.address)}
					>
						{assemblyAddress(item?.address)}
					</p>
				</div>
				<Row className="pl-2 pt-1 text-[rgba(0,0,0,0.65)]">
					<Col md={10} className="common-ellipsis" title={item?.property_type}>
						{item?.property_type || '--'}
					</Col>
					<Col md={7} className="common-ellipsis" title={item?.unit}>
						{item?.unit ? `${item?.unit} units` : '--'}
					</Col>
					<Col
						md={7}
						className="common-ellipsis"
						title={formatNumber(item?.living_area_sqft)}
					>
						{formatNumber(item?.living_area_sqft)}
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
export default CardItem;
