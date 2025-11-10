'use client';
import React from 'react';
import styles from './index.module.scss';
import { propertiesFieldArr, propertyList } from '@/enum/properties';
import { Button, Col, message, Row, Space } from 'antd';
import { useRouter } from 'next/navigation';
import { formatDate, assemblyAddress, createRequest } from '@repo/shared-methods/utils';
import FieldTips from '@/components/FieldTips';
import { BtnPopUpModal, DeleteBtn } from '@repo/shared-methods';
import AddModal from '../MoreCard/Add';
import { DeleteOutlined } from '@ant-design/icons';
import { useState } from 'react';
const PropertiesItem = ({ data, handleSearch }) => {
	const { request } = createRequest();
	const [loading, setLoading] = useState(false);
	const router = useRouter();
	const nameList = [
		{
			label: 'Active borrower',
			value: 'applicants',
			render: (v) => <FieldTips title={v?.join(',')}>{v?.join(',') || '--'}</FieldTips>,
		},
		{
			label: 'Active lender',
			value: 'lenders',
			render: (v) => <FieldTips title={v?.join(',')}>{v?.join(',') || '--'}</FieldTips>,
		},
	];

	const getItem = (item, index) => {
		return (
			<div key={index} className={[item?.className || 'w-1/2', 'pb-4'].join(' ')}>
				<p className="text-base-gray text-sm pb-1">{item.label}</p>
				<span className="text-black-dark text-base font-medium">
					{(item?.render ? item?.render?.(data[item.value]) : data[item.value]) || '--'}
				</span>
			</div>
		);
	};

	const handleDelete = async () => {
		try {
			setLoading(true);
			await request(`/property/${data?.id}`, { method: 'DELETE' });
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
		<div className={[styles.propertiesItem].join(' ')}>
			<div className={styles.header}>
				<p className={styles.address}>{assemblyAddress(data?.address)}</p>
				<Space className="flex items-center shrink-0 gap-2">
					<Button onClick={() => router.push(`/properties/${data?.id}`)}>
						View details
					</Button>
					<BtnPopUpModal
						getContent={(setOpen) => (
							<AddModal
								setOpen={setOpen}
								rowData={data}
								handleSearch={handleSearch}
							/>
						)}
						type="default"
						size="medium"
					>
						Edit
					</BtnPopUpModal>

					<DeleteBtn
						type="default"
						title="Are you sure delete this Property?"
						content={assemblyAddress(data?.address)}
						size="medium"
						onOk={handleDelete}
						danger={!data?.loan_count}
						className={data?.loan_count && '!text-base-disabled'}
						predicate={() => {
							if (data?.loan_count) {
								message.destroy();
								message.open({
									type: 'warning',
									content:
										'This property has an associated loan, so it cannot be deleted.',
								});
								return false;
							}
							return true;
						}}
						confirmLoading={loading}
					>
						<DeleteOutlined />
					</DeleteBtn>
				</Space>
			</div>
			<Row className="py-4 px-4">
				<Col
					lg={8}
					sm={24}
					xs={24}
					className="border-base-border border-r-0 lg:pr-6 lg:border-r "
				>
					<p className="text-sm font-bold text-base-gray pb-4">Property Details</p>
					<div className="flex flex-wrap">
						{propertyList.map((item, index) => {
							return getItem(item, index);
						})}
					</div>
				</Col>
				<Col
					lg={9}
					sm={24}
					xs={24}
					className="border-base-border border-r-0 border-t lg:border-t-0 lg:px-6 lg:border-r "
				>
					<p className="text-sm font-bold text-base-gray pb-4 pt-4 lg:pt-0">
						Loan Information
					</p>
					<div className="flex flex-wrap">
						{propertiesFieldArr.map((item, index) => {
							return getItem(item, index);
						})}
					</div>
				</Col>
				<Col lg={7} sm={24} xs={24} className="border-t lg:border-t-0 lg:px-6 ">
					<p className="text-sm font-bold text-base-gray pb-4 pt-4 lg:pt-0">
						Loan Status
					</p>
					<div className="flex flex-wrap gap-1">
						{!!data?.active_loan_count && (
							<div className={styles.active}>
								{data?.active_loan_count || 0} Active Loans
							</div>
						)}
						{!!data?.pending_loan_count && (
							<div className={styles.pending}>
								{data?.pending_loan_count || 0} Pending Loans
							</div>
						)}
						{!data?.pending_loan_count && !data?.active_loan_count && (
							<span className="text-black-dark text-base font-medium">--</span>
						)}
					</div>
					<div className="pt-4">
						<p className="text-base-gray text-sm pb-1">Nearest maturity date</p>
						<span className="text-black-dark text-base font-medium">
							{formatDate(data?.nearest_maturity_date) || '--'}
						</span>
					</div>
				</Col>
			</Row>
			<Row className="bg-[#F5F5F5] rounded-b-md">
				{nameList.map((item, index) => {
					return (
						<Col
							lg={12}
							sm={24}
							xs={24}
							key={index}
							className={['border-t border-base-border pt-4 px-4']}
						>
							{getItem(item, index)}
						</Col>
					);
				})}
			</Row>
		</div>
	);
};
export default PropertiesItem;
